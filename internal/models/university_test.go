package models

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUniversity_Validation(t *testing.T) {
	tests := []struct {
		name    string
		uni     University
		wantErr bool
	}{
		{
			name: "valid university",
			uni: University{
				ID:        primitive.NewObjectID(),
				Name:      "Test University",
				Address:   "123 Test St",
				Phone:    "(11) 1234-5678",
				Email:    "test@university.edu",
				Website:  "https://test.edu",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			uni: University{
				ID:        primitive.NewObjectID(),
				Name:      "Test University",
				Address:   "123 Test St",
				Phone:    "(11) 1234-5678",
				Email:    "invalid-email",
				Website:  "https://test.edu",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.uni.ID.IsZero() {
				t.Error("ID should not be zero")
			}
			if tt.uni.CreatedAt.IsZero() {
				t.Error("CreatedAt should not be zero")
			}
			if tt.uni.UpdatedAt.IsZero() {
				t.Error("UpdatedAt should not be zero")
			}
		})
	}
}

func TestUniversityResponse_Structure(t *testing.T) {
	uni := University{
		ID:        primitive.NewObjectID(),
		Name:      "Test University",
		Address:   "123 Test St",
		Phone:    "(11) 1234-5678",
		Email:    "test@university.edu",
		Website:  "https://test.edu",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	resp := UniversityResponse{
		Status:  200,
		Message: "Success",
		Data:    uni,
	}

	if resp.Status != 200 {
		t.Errorf("Expected status 200, got %d", resp.Status)
	}

	if resp.Message != "Success" {
		t.Errorf("Expected message 'Success', got %s", resp.Message)
	}

	if resp.Data.ID != uni.ID {
		t.Error("University data mismatch")
	}
}