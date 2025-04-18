package repository

import (
	"context"
	"time"

	"github.com/university-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UniversityRepository struct {
	collection *mongo.Collection
}

func NewUniversityRepository(db *mongo.Database) *UniversityRepository {
	return &UniversityRepository{
		collection: db.Collection("universities"),
	}
}

func (r *UniversityRepository) Create(ctx context.Context, university *models.University) error {
	university.CreatedAt = time.Now()
	university.UpdatedAt = time.Now()
	
	result, err := r.collection.InsertOne(ctx, university)
	if err != nil {
		return err
	}
	
	university.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *UniversityRepository) GetByID(ctx context.Context, id string) (*models.University, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var university models.University
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&university)
	if err != nil {
		return nil, err
	}

	return &university, nil
}

func (r *UniversityRepository) GetAll(ctx context.Context) ([]models.University, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var universities []models.University
	if err = cursor.All(ctx, &universities); err != nil {
		return nil, err
	}

	return universities, nil
}

func (r *UniversityRepository) Update(ctx context.Context, university *models.University) error {
	university.UpdatedAt = time.Now()

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": university.ID},
		bson.M{"$set": university},
	)
	return err
}

func (r *UniversityRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}