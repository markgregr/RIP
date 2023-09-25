package main

import (
	"context"
	"log"

	"github.com/markgregr/RIP/internal/app"
)

func main() {
	log.Println("Application start!")
	// Создайте контекст
	ctx := context.Background()

	// Создайте объект Application, передав контекст
	application, err := app.New(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Запустите сервер, вызвав метод StartServer у объекта Application
	application.StartServer()
	log.Println("Application terminated!")
}
