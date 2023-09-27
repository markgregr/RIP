package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/markgregr/RIP/internal/app/ds"
	"github.com/markgregr/RIP/internal/app/dsn"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&ds.Baggage{},&ds.Delivery{},&ds.DeliveryBagggage{},&ds.User{})
	if err != nil {
		panic("cant migrate db")
	}
	
}