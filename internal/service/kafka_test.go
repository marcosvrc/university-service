package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/university-service/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupTestKafka(t *testing.T) (*kafka.Reader, func()) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "test_topic",
		GroupID: "test_group",
	})

	return reader, func() {
		if err := reader.Close(); err != nil {
			t.Errorf("Failed to close Kafka reader: %v", err)
		}
	}
}

func TestKafkaService_PublishUniversityEvent(t *testing.T) {
	reader, cleanup := setupTestKafka(t)
	defer cleanup()

	service := NewKafkaService([]string{"localhost:9092"}, "test_topic")
	defer service.Close()

	ctx := context.Background()

	t.Run("Publish Create Event", func(t *testing.T) {
		uni := &models.University{
			ID:        primitive.NewObjectID(),
			Name:      "Test University",
			Address:   "123 Test St",
			Phone:    "(11) 1234-5678",
			Email:     "test@university.edu",
			Website:   "https://test.edu",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := service.PublishUniversityEvent(ctx, "university_created", uni)
		assert.NoError(t, err)

		// Read the message from Kafka
		msg, err := reader.ReadMessage(ctx)
		assert.NoError(t, err)

		var event struct {
			Type       string             `json:"type"`
			University *models.University `json:"university"`
		}

		err = json.Unmarshal(msg.Value, &event)
		assert.NoError(t, err)
		assert.Equal(t, "university_created", event.Type)
		assert.Equal(t, uni.ID, event.University.ID)
		assert.Equal(t, uni.Name, event.University.Name)
	})

	t.Run("Publish Update Event", func(t *testing.T) {
		uni := &models.University{
			ID:        primitive.NewObjectID(),
			Name:      "Updated University",
			Address:   "456 Test St",
			Phone:    "(11) 8765-4321",
			Email:     "updated@university.edu",
			Website:   "https://updated.edu",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := service.PublishUniversityEvent(ctx, "university_updated", uni)
		assert.NoError(t, err)

		msg, err := reader.ReadMessage(ctx)
		assert.NoError(t, err)

		var event struct {
			Type       string             `json:"type"`
			University *models.University `json:"university"`
		}

		err = json.Unmarshal(msg.Value, &event)
		assert.NoError(t, err)
		assert.Equal(t, "university_updated", event.Type)
		assert.Equal(t, uni.ID, event.University.ID)
		assert.Equal(t, uni.Name, event.University.Name)
	})

	t.Run("Publish Delete Event", func(t *testing.T) {
		uni := &models.University{
			ID:        primitive.NewObjectID(),
			Name:      "Deleted University",
			Address:   "789 Test St",
			Phone:    "(11) 9999-9999",
			Email:     "deleted@university.edu",
			Website:   "https://deleted.edu",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := service.PublishUniversityEvent(ctx, "university_deleted", uni)
		assert.NoError(t, err)

		msg, err := reader.ReadMessage(ctx)
		assert.NoError(t, err)

		var event struct {
			Type       string             `json:"type"`
			University *models.University `json:"university"`
		}

		err = json.Unmarshal(msg.Value, &event)
		assert.NoError(t, err)
		assert.Equal(t, "university_deleted", event.Type)
		assert.Equal(t, uni.ID, event.University.ID)
		assert.Equal(t, uni.Name, event.University.Name)
	})
}