package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/markgregr/RIP/internal/app/ds"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetActiveBaggageByID(baggage_id int) (*ds.Baggage, error) {
	baggage := &ds.Baggage{}

	err := r.db.First(baggage, "baggage_id = ? AND baggage_status = ?", baggage_id, ds.BAGGAGE_STATUS_ACTIVE).Error // find product with id = 1
	if err != nil {
		return nil, err
	}

	return baggage, nil
}

func (r *Repository) GetActiveBaggages() ([]ds.Baggage,error) {
	var baggages []ds.Baggage
	if err := r.db.Find(&baggages, "baggage_status = ?", ds.BAGGAGE_STATUS_ACTIVE).Error; err != nil {
        return nil, err
    }
	return baggages, nil
}

func (r *Repository) DeleteBaggage(baggage_id int) error {
	return r.db.Exec("UPDATE baggages SET baggage_status = ? WHERE baggage_id = ?", ds.BAGGAGE_STATUS_DELETED, baggage_id).Error
}
