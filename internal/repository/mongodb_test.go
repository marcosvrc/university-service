package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/university-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) (*mongo.Database, func()) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	db := client.Database("test_db")

	return db, func() {
		if err := db.Drop(ctx); err != nil {
			t.Errorf("Failed to drop test database: %v", err)
		}
		if err := client.Disconnect(ctx); err != nil {
			t.Errorf("Failed to disconnect from MongoDB: %v", err)
		}
	}
}

func TestUniversityRepository_CRUD(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUniversityRepository(db)
	ctx := context.Background()

	// Test Create
	t.Run("Create", func(t *testing.T) {
		uni := &models.University{
			Name:    "Test University",
			Address: "123 Test St",
			Phone:   "(11) 1234-5678",
			Email:   "test@university.edu",
			Website: "https://test.edu",
		}

		err := repo.Create(ctx, uni)
		assert.NoError(t, err)
		assert.False(t, uni.ID.IsZero())
		assert.False(t, uni.CreatedAt.IsZero())
		assert.False(t, uni.UpdatedAt.IsZero())
	})

	// Test GetByID
	t.Run("GetByID", func(t *testing.T) {
		uni := &models.University{
			ID:      primitive.NewObjectID(),
			Name:    "Test University",
			Address: "123 Test St",
			Phone:   "(11) 1234-5678",
			Email:   "test@university.edu",
			Website: "https://test.edu",
		}

		_, err := repo.collection.InsertOne(ctx, uni)
		assert.NoError(t, err)

		found, err := repo.GetByID(ctx, uni.ID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, uni.ID, found.ID)
		assert.Equal(t, uni.Name, found.Name)
	})

	// Test GetAll
	t.Run("GetAll", func(t *testing.T) {
		// Clear collection
		_, err := repo.collection.DeleteMany(ctx, bson.M{})
		assert.NoError(t, err)

		// Insert test data
		unis := []interface{}{
			&models.University{
				ID:      primitive.NewObjectID(),
				Name:    "University 1",
				Address: "Address 1",
				Phone:   "(11) 1111-1111",
				Email:   "uni1@test.edu",
			},
			&models.University{
				ID:      primitive.NewObjectID(),
				Name:    "University 2",
				Address: "Address 2",
				Phone:   "(11) 2222-2222",
				Email:   "uni2@test.edu",
			},
		}

		_, err = repo.collection.InsertMany(ctx, unis)
		assert.NoError(t, err)

		results, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		assert.Len(t, results, 2)
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		uni := &models.University{
			ID:      primitive.NewObjectID(),
			Name:    "Old Name",
			Address: "Old Address",
			Phone:   "(11) 1234-5678",
			Email:   "old@test.edu",
			Website: "https://old.edu",
		}

		_, err := repo.collection.InsertOne(ctx, uni)
		assert.NoError(t, err)

		uni.Name = "New Name"
		uni.UpdatedAt = time.Now()

		err = repo.Update(ctx, uni)
		assert.NoError(t, err)

		updated, err := repo.GetByID(ctx, uni.ID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, "New Name", updated.Name)
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		uni := &models.University{
			ID:      primitive.NewObjectID(),
			Name:    "To Delete",
			Address: "Delete Address",
			Phone:   "(11) 1234-5678",
			Email:   "delete@test.edu",
		}

		_, err := repo.collection.InsertOne(ctx, uni)
		assert.NoError(t, err)

		err = repo.Delete(ctx, uni.ID.Hex())
		assert.NoError(t, err)

		_, err = repo.GetByID(ctx, uni.ID.Hex())
		assert.Error(t, err)
	})
}