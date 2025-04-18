package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/university-service/api"
	"github.com/university-service/config"
	"github.com/university-service/internal/repository"
	"github.com/university-service/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Carregar configurações
	cfg := config.LoadConfig()

	// Conectar ao MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Ping no MongoDB para verificar a conexão
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(cfg.MongoDB.Database)

	// Inicializar repositório
	repo := repository.NewUniversityRepository(db)

	// Inicializar serviço Kafka
	kafkaService := service.NewKafkaService(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	defer kafkaService.Close()

	// Inicializar handler
	handler := api.NewHandler(repo, kafkaService)

	// Configurar router
	router := gin.Default()

	// Rotas
	router.POST("/universities", handler.CreateUniversity)
	router.GET("/universities/:id", handler.GetUniversity)
	router.GET("/universities", handler.ListUniversities)
	router.PUT("/universities/:id", handler.UpdateUniversity)
	router.DELETE("/universities/:id", handler.DeleteUniversity)

	// Iniciar servidor
	if err := router.Run(cfg.Server.Port); err != nil {
		log.Fatal(err)
	}
}