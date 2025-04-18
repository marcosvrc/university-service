package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/university-service/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock do UniversityRepository
type MockUniversityRepository struct {
	mock.Mock
}

func (m *MockUniversityRepository) Create(ctx context.Context, university *models.University) error {
	args := m.Called(ctx, university)
	return args.Error(0)
}

func (m *MockUniversityRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.University, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.University), args.Error(1)
}

func (m *MockUniversityRepository) GetAll(ctx context.Context) ([]*models.University, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.University), args.Error(1)
}

func (m *MockUniversityRepository) Update(ctx context.Context, university *models.University) error {
	args := m.Called(ctx, university)
	return args.Error(0)
}

func (m *MockUniversityRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Mock do KafkaService
type MockKafkaService struct {
	mock.Mock
}

func (m *MockKafkaService) PublishUniversityEvent(ctx context.Context, eventType string, university *models.University) error {
	args := m.Called(ctx, eventType, university)
	return args.Error(0)
}

func (m *MockKafkaService) Close() error {
	args := m.Called()
	return args.Error(0)
}

func setupTestRouter(repo *MockUniversityRepository, kafka *MockKafkaService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := NewHandler(repo, kafka)
	handler.RegisterRoutes(r)
	return r
}

func TestHandler_CreateUniversity(t *testing.T) {
	repo := new(MockUniversityRepository)
	kafka := new(MockKafkaService)
	router := setupTestRouter(repo, kafka)

	t.Run("Successful Creation", func(t *testing.T) {
		uni := &models.University{
			Name:    "Test University",
			Address: "123 Test St",
			Phone:   "(11) 1234-5678",
			Email:   "test@university.edu",
			Website: "https://test.edu",
		}

		repo.On("Create", mock.Anything, mock.AnythingOfType("*models.University")).Return(nil)
		kafka.On("PublishUniversityEvent", mock.Anything, "university_created", mock.AnythingOfType("*models.University")).Return(nil)

		body, _ := json.Marshal(uni)
		req := httptest.NewRequest("POST", "/universities", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response models.UniversityResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Universidade criada com sucesso", response.Message)

		repo.AssertExpectations(t)
		kafka.AssertExpectations(t)
	})
}

func TestHandler_GetUniversity(t *testing.T) {
	repo := new(MockUniversityRepository)
	kafka := new(MockKafkaService)
	router := setupTestRouter(repo, kafka)

	t.Run("Successful Retrieval", func(t *testing.T) {
		id := primitive.NewObjectID()
		uni := &models.University{
			ID:        id,
			Name:      "Test University",
			Address:   "123 Test St",
			Phone:    "(11) 1234-5678",
			Email:     "test@university.edu",
			Website:   "https://test.edu",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		repo.On("GetByID", mock.Anything, id).Return(uni, nil)

		req := httptest.NewRequest("GET", "/universities/"+id.Hex(), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.University
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, uni.Name, response.Name)

		repo.AssertExpectations(t)
	})
}

func TestHandler_ListUniversities(t *testing.T) {
	repo := new(MockUniversityRepository)
	kafka := new(MockKafkaService)
	router := setupTestRouter(repo, kafka)

	t.Run("Successful List", func(t *testing.T) {
		universities := []*models.University{
			{
				ID:        primitive.NewObjectID(),
				Name:      "University 1",
				Address:   "123 Test St",
				Phone:    "(11) 1234-5678",
				Email:     "test1@university.edu",
				Website:   "https://test1.edu",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        primitive.NewObjectID(),
				Name:      "University 2",
				Address:   "456 Test St",
				Phone:    "(11) 8765-4321",
				Email:     "test2@university.edu",
				Website:   "https://test2.edu",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		repo.On("GetAll", mock.Anything).Return(universities, nil)

		req := httptest.NewRequest("GET", "/universities", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []*models.University
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)

		repo.AssertExpectations(t)
	})
}

func TestHandler_UpdateUniversity(t *testing.T) {
	repo := new(MockUniversityRepository)
	kafka := new(MockKafkaService)
	router := setupTestRouter(repo, kafka)

	t.Run("Successful Update", func(t *testing.T) {
		id := primitive.NewObjectID()
		uni := &models.University{
			ID:      id,
			Name:    "Updated University",
			Address: "789 Test St",
			Phone:   "(11) 9999-9999",
			Email:   "updated@university.edu",
			Website: "https://updated.edu",
		}

		repo.On("Update", mock.Anything, mock.AnythingOfType("*models.University")).Return(nil)
		kafka.On("PublishUniversityEvent", mock.Anything, "university_updated", mock.AnythingOfType("*models.University")).Return(nil)

		body, _ := json.Marshal(uni)
		req := httptest.NewRequest("PUT", "/universities/"+id.Hex(), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.UniversityResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Universidade atualizada com sucesso", response.Message)

		repo.AssertExpectations(t)
		kafka.AssertExpectations(t)
	})
}

func TestHandler_DeleteUniversity(t *testing.T) {
	repo := new(MockUniversityRepository)
	kafka := new(MockKafkaService)
	router := setupTestRouter(repo, kafka)

	t.Run("Successful Delete", func(t *testing.T) {
		id := primitive.NewObjectID()
		uni := &models.University{
			ID:      id,
			Name:    "University to Delete",
			Address: "999 Test St",
			Phone:   "(11) 0000-0000",
			Email:   "delete@university.edu",
			Website: "https://delete.edu",
		}

		repo.On("GetByID", mock.Anything, id).Return(uni, nil)
		repo.On("Delete", mock.Anything, id).Return(nil)
		kafka.On("PublishUniversityEvent", mock.Anything, "university_deleted", mock.AnythingOfType("*models.University")).Return(nil)

		req := httptest.NewRequest("DELETE", "/universities/"+id.Hex(), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.UniversityResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Universidade excluída com sucesso", response.Message)

		repo.AssertExpectations(t)
		kafka.AssertExpectations(t)
	})
}