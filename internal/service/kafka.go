package service

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"github.com/university-service/internal/models"
)

type KafkaService struct {
	writer *kafka.Writer
}

func NewKafkaService(brokers []string, topic string) *KafkaService {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})

	return &KafkaService{
		writer: writer,
	}
}

func (s *KafkaService) PublishUniversityEvent(ctx context.Context, eventType string, university *models.University) error {
	event := struct {
		Type       string             `json:"type"`
		University *models.University `json:"university"`
	}{
		Type:       eventType,
		University: university,
	}

	value, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return s.writer.WriteMessages(ctx, kafka.Message{
		Value: value,
	})
}

func (s *KafkaService) Close() error {
	return s.writer.Close()
}