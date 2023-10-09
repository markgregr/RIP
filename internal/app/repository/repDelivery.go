package repository

import (
	"strings"

	"github.com/markgregr/RIP/internal/app/ds"
)

//методы для таблицы delivery
func (r *Repository) GetDeliveries(searchFlightNumber string) ([]ds.Delivery,error) {
	
	searchFlightNumber = strings.ToUpper(searchFlightNumber+"%")
	var deliveries []ds.Delivery
	if err := r.db.Find(&deliveries, "delivery_status != ? AND flight_number LIKE ?", ds.DELIVERY_STATUS_DELETED, searchFlightNumber).Error; err != nil {
        return nil, err
    }
	return deliveries, nil
}

func (r *Repository) GetDeliveryByID(deliveryID int) (*ds.Delivery, error) {
	delivery := &ds.Delivery{}

	err := r.db.First(delivery, "delivery_status != ? AND delivery_id = ?", ds.BAGGAGE_STATUS_DELETED, deliveryID).Error
	if err != nil {
		return nil, err
	}

	return delivery, nil
}
func (r *Repository) DeleteDelivery(deliveryID int) error {
	return r.db.Exec("UPDATE deliveries SET delivery_status = ? WHERE delivery_id = ?", ds.DELIVERY_STATUS_DELETED, deliveryID).Error
}
func (r *Repository) UpdateDelivery(deliveryID int, updatedDelivery *ds.Delivery) error {
	err := r.db.Model(&ds.Delivery{}).Where("delivery_status != ? AND delivery_id = ?", ds.BAGGAGE_STATUS_DELETED, deliveryID).Updates(updatedDelivery).Error
	if err != nil {
		return err
	}
	return nil
}
