package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/university-service/internal/models"
	"github.com/university-service/internal/repository"
	"github.com/university-service/internal/service"
)

type Handler struct {
	repo  *repository.UniversityRepository
	kafka *service.KafkaService
}

func NewHandler(repo *repository.UniversityRepository, kafka *service.KafkaService) *Handler {
	return &Handler{
		repo:  repo,
		kafka: kafka,
	}
}

func (h *Handler) CreateUniversity(c *gin.Context) {
	var university models.University
	if err := c.ShouldBindJSON(&university); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Create(c.Request.Context(), &university); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Publicar evento no Kafka
	if err := h.kafka.PublishUniversityEvent(c.Request.Context(), "university_created", &university); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish event"})
		return
	}

	c.JSON(http.StatusCreated, models.UniversityResponse{
		Status:  http.StatusCreated,
		Message: "University created successfully",
		Data:    university,
	})
}

func (h *Handler) GetUniversity(c *gin.Context) {
	id := c.Param("id")
	university, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "university not found"})
		return
	}

	c.JSON(http.StatusOK, models.UniversityResponse{
		Status:  http.StatusOK,
		Message: "University retrieved successfully",
		Data:    *university,
	})
}

func (h *Handler) ListUniversities(c *gin.Context) {
	universities, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Universities retrieved successfully",
		"data":    universities,
	})
}

func (h *Handler) UpdateUniversity(c *gin.Context) {
	id := c.Param("id")
	var university models.University
	if err := c.ShouldBindJSON(&university); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUniversity, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "university not found"})
		return
	}

	university.ID = existingUniversity.ID
	if err := h.repo.Update(c.Request.Context(), &university); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Publicar evento no Kafka
	if err := h.kafka.PublishUniversityEvent(c.Request.Context(), "university_updated", &university); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish event"})
		return
	}

	c.JSON(http.StatusOK, models.UniversityResponse{
		Status:  http.StatusOK,
		Message: "University updated successfully",
		Data:    university,
	})
}

func (h *Handler) DeleteUniversity(c *gin.Context) {
	id := c.Param("id")
	
	university, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "university not found"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Publicar evento no Kafka
	if err := h.kafka.PublishUniversityEvent(c.Request.Context(), "university_deleted", university); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "University deleted successfully",
	})
}