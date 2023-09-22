package main

import (
	"log"

	"github.com/markgregr/RIP/internal/api"
)

func main() {
	log.Println("Application start!")
	api.StartServer()
	log.Println("Application terminated!")
}
