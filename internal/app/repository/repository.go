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

//методы для таблицы baggage
func (r *Repository) GetBaggages() ([]ds.Baggage, error) {
	var baggages []ds.Baggage
	if err := r.db.Find(&baggages, "baggage_status = ?", ds.BAGGAGE_STATUS_ACTIVE).Error; err != nil {
		return nil, err
	}
	return baggages, nil
}
func (r *Repository) GetBaggageByID(baggageID int) (*ds.Baggage, error) {
	baggage := &ds.Baggage{}

	err := r.db.First(baggage, "baggage_id = ? AND baggage_status = ?", baggageID, ds.BAGGAGE_STATUS_ACTIVE).Error
	if err != nil {
		return nil, err
	}

	return baggage, nil
}
func (r *Repository) CreateBaggage(baggage *ds.Baggage) error {
	return r.db.Create(baggage).Error
}
func (r *Repository) DeleteBaggage(baggageID int) error {
	return r.db.Exec("UPDATE baggages SET baggage_status = ? WHERE baggage_id = ?", ds.BAGGAGE_STATUS_DELETED, baggageID).Error
}
func (r *Repository) UpdateBaggage(baggageID int, updatedBaggage *ds.Baggage) error {
	err := r.db.Model(&ds.Baggage{}).Where("baggage_id = ? AND baggage_status = ?", baggageID, ds.BAGGAGE_STATUS_ACTIVE).Updates(updatedBaggage).Error
	if err != nil {
		return err
	}
	return nil
}


//методы для таблицы delivery
func (r *Repository) GetDeliveries() ([]ds.Delivery, error) {
	var deliveries []ds.Delivery
	if err := r.db.Find(&deliveries).Error; err != nil {
		return nil, err
	}
	return deliveries, nil
}
func (r *Repository) GetDeliveryByID(deliveryID int) (*ds.Delivery, error) {
	delivery := &ds.Delivery{}

	err := r.db.First(delivery, "delivery_id = ?", deliveryID).Error
	if err != nil {
		return nil, err
	}

	return delivery, nil
}
func (r *Repository) CreateDelivery(delivery *ds.Delivery) error {
	return r.db.Create(delivery).Error
}
func (r *Repository) DeleteDelivery(deliveryID int) error {
	return r.db.Exec("DELETE FROM deliveries WHERE delivery_id = ?", deliveryID).Error
}
func (r *Repository) UpdateDelivery(deliveryID int, updatedDelivery *ds.Delivery) error {
	err := r.db.Model(&ds.Delivery{}).Where("delivery_id = ?", deliveryID).Updates(updatedDelivery).Error
	if err != nil {
		return err
	}
	return nil
}


//методы для таблицы user

func (r *Repository) GetUsers() ([]ds.User, error) {
	var users []ds.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
func (r *Repository) GetUserByID(userID int) (*ds.User, error) {
	user := &ds.User{}

	err := r.db.First(user, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (r *Repository) CreateUser(user *ds.User) error {
	return r.db.Create(user).Error
}
func (r *Repository) DeleteUser(userID int) error {
	return r.db.Exec("DELETE FROM users WHERE user_id = ?", userID).Error
}
func (r *Repository) UpdateUser(userID int, updatedUser *ds.User) error {
	err := r.db.Model(&ds.User{}).Where("user_id = ?", userID).Updates(updatedUser).Error
	if err != nil {
		return err
	}
	return nil
}


