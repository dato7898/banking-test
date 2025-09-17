package repository

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer interface {
	Publish(key, value interface{}) error
}

type kafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(writer *kafka.Writer) KafkaProducer {
	return &kafkaProducer{writer: writer}
}

func (p *kafkaProducer) Publish(key, value interface{}) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}

	msg := kafka.Message{Value: payload}
	if key != nil {
		keyBytes, _ := json.Marshal(key)
		msg.Key = keyBytes
	}

	if err := p.writer.WriteMessages(context.Background(), msg); err != nil {
		log.Println("Kafka publish failed:", err)
		return err
	}
	return nil
}
