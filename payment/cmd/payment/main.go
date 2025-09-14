package main

import (
	"log"
	"payment/internal/config"
	"payment/internal/controller"
	"payment/internal/repository"
	"payment/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

func ensureTopic(brokerAddress string, topic string, partitions int, replicationFactor int) {
	conn, err := kafka.Dial("tcp", brokerAddress)
	if err != nil {
		log.Fatal("failed to connect to kafka: %s", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Fatal("failed to get controller: %s", err)
	}

	controllerAddr := controller.Host + ":" + strconv.Itoa(controller.Port)

	controllerConn, err := kafka.Dial("tcp", controllerAddr)
	if err != nil {
		log.Fatal("failed to connect to controller: %s", err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     partitions,
			ReplicationFactor: replicationFactor,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		log.Println("create topic error (may already exist): %s", err)
	}
}

func main() {
	cfg := config.Load()

	db, err := sqlx.Connect("pgx", cfg.PostgresDSN)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP(cfg.KafkaBrokers...),
		Topic:    cfg.KafkaPaymentsTopic,
		Balancer: &kafka.LeastBytes{},
	}
	defer kafkaWriter.Close()

	conn, err := grpc.Dial(cfg.AccountAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("failed to connect to integration service:", err)
	}
	defer conn.Close()
	accountClient := repository.NewAccountClient(conn)

	paymentRepo := repository.NewPaymentRepository(db)
	kafkaProducer := repository.NewKafkaProducer(kafkaWriter)
	paymentService := service.NewPaymentService(paymentRepo, accountClient, kafkaProducer)
	paymentHandler := controller.NewPaymentHandler(paymentService)
	authHandler := controller.NewAuthHandler(cfg.JWTSecret)

	r := gin.Default()

	paymentGroup := r.Group("/")
	paymentGroup.Use(authHandler.Middleware())
	{
		paymentGroup.POST("/payments", paymentHandler.CreatePayment)
	}

	log.Println("🚀 Payment service started on", cfg.HTTPPort)
	r.Run(":" + cfg.HTTPPort)
}
