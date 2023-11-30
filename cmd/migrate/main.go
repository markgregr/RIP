package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/markgregr/RIP/internal/dsn"
	"github.com/markgregr/RIP/internal/model"
)

func main() {
    _ = godotenv.Load()
    db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    // Явно мигрировать только нужные таблицы
    err = db.AutoMigrate(&model.User{})
    if err != nil {
        panic("cant migrate db")
    }
}
