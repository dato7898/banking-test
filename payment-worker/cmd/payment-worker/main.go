package main

import (
	"log"

	"payment-worker/internal/config"
	"payment-worker/internal/controller"
	"payment-worker/internal/repository"
	"payment-worker/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	db, err := sqlx.Connect("pgx", cfg.PostgresDSN)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	paymentRepository := repository.NewPaymentRepository(db)

	conn, err := grpc.Dial(cfg.IntegrationAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("failed to connect to integration service:", err)
	}
	defer conn.Close()
	paymentIntegrationClient := repository.NewPaymentIntegrationClient(conn)

	paymentProcessor := service.NewPaymentProcessor(paymentRepository, paymentIntegrationClient)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.KafkaBrokers,
		Topic:   cfg.KafkaPaymentsTopic,
		GroupID: "payment-workers",
	})
	defer reader.Close()

	kafkaConsumer := controller.NewKafkaConsumer(paymentProcessor, reader)

	log.Println("🚀 Payment Worker started, listening for messages...")
	kafkaConsumer.Start()
}
