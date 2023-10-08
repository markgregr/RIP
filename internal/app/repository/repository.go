package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/markgregr/RIP/internal/app/ds"
)

type Repository struct {
	db *gorm.DB
}


