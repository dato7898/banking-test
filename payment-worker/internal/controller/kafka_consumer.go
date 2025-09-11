package controller

import (
	"context"
	"log"
	"payment-worker/internal/service"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	processor service.PaymentProcess
	reader    *kafka.Reader
}

func NewKafkaConsumer(processor service.PaymentProcess, reader *kafka.Reader) *KafkaConsumer {
	return &KafkaConsumer{processor: processor, reader: reader}
}

func (c *KafkaConsumer) Start() {
	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Kafka reader error:", err)
			continue
		}

		if err := c.processor.ProcessMessage(msg.Value); err != nil {
			log.Println("Processing failed:", err)
		} else {
			log.Println("✅ Payment processed successfully")
		}
	}
}
