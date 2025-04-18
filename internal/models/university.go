package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type University struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string            `bson:"name" json:"name" binding:"required"`
	Address     string            `bson:"address" json:"address" binding:"required"`
	Phone       string            `bson:"phone" json:"phone" binding:"required"`
	Email       string            `bson:"email" json:"email" binding:"required,email"`
	Website     string            `bson:"website" json:"website"`
	CreatedAt   time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `bson:"updated_at" json:"updated_at"`
}

type UniversityResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    University `json:"data"`
}