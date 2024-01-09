package main

import (
	"context"
	"log"

	"github.com/markgregr/RIP/internal/app"
)

// @title BagTracker RestAPI
// @version 1.0
// @description API server for BagTracker application

// @host http://localhost:8081
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	log.Println("Application start!")
	// Создаем контекст
	ctx := context.Background()

	// Создаем Aplication
	application, err := app.New(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Запустите сервер, вызвав метод StartServer у объекта Application
	application.Run()
	log.Println("Application terminated!")
}
	