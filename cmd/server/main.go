package main

import (
	"log"
	config "mobile/internal/configs"
	"mobile/internal/database"
	"mobile/internal/handler"
	"mobile/internal/repository"
	"mobile/internal/service"
	"mobile/internal/utils"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Загружаем окружение из .env файла в переменную config
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// База данных
	db, err := database.NewDataBase(cfg)
	if err != nil {
		panic(err)
	}

	//	Логи
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Could not parse log level: %v", err)
		os.Exit(1)
	}

	logger := utils.InitLogger(level)

	// Репозиторий
	subscriptionRepo := repository.NewSubscriptionRepo(db)

	// Бизнес-логика
	subscriptionSrv := service.NewSubscriptionService(subscriptionRepo)

	subscriptionHandler := handler.NewMsHandler(subscriptionSrv, logger)

	router := gin.Default()

	router.POST("/subscription/create", subscriptionHandler.CreateSubscription)
	router.GET("/subscription", subscriptionHandler.GetSubscription)
	router.PATCH("/subscription", subscriptionHandler.UpdateSubscription)
	router.DELETE("/subscription", subscriptionHandler.DeleteSubscription)
	router.GET("/subscription/report", subscriptionHandler.Report)

	//Запускаем сервер
	logger.Infof("Starting server on port %s...", cfg.ApiPort)
	if err := router.Run(":" + cfg.ApiPort); err != nil {
		logger.Fatalf("Error starting server: %v", err)
	}
}
