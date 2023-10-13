package repository

import (
	"strings"

	"github.com/markgregr/RIP/internal/app/ds"
)

//методы для таблицы delivery
func (r *Repository) GetDeliveries(searchFlightNumber, startFormationDate, endFormationDate string) ([]map[string]interface{}, error) {
    searchFlightNumber = strings.ToUpper(searchFlightNumber + "%")
    var deliveries []map[string]interface{}
    query := r.db.Table("deliveries").
        Select("deliveries.flight_number, deliveries.creation_date, deliveries.formation_date, deliveries.completion_date, deliveries.delivery_status").
        Joins("JOIN delivery_baggages ON deliveries.delivery_id = delivery_baggages.delivery_id").
        Joins("JOIN baggages ON baggages.baggage_id = delivery_baggages.baggage_id").
        Where("deliveries.delivery_status != ? AND deliveries.flight_number LIKE ?", ds.DELIVERY_STATUS_DELETED, searchFlightNumber)

    if startFormationDate != "" && endFormationDate != "" {
        query = query.Where("deliveries.formation_date BETWEEN ? AND ?", startFormationDate, endFormationDate)
    }

    if err := query.Scan(&deliveries).Error; err != nil {
        return nil, err
    }

    return deliveries, nil
}

func (r *Repository) GetDeliveryByID(deliveryID int) (map[string]interface{}, error) {
    var delivery map[string]interface{}
    if err := r.db.
        Table("deliveries").
        Select("deliveries.flight_number, deliveries.creation_date, deliveries.formation_date, deliveries.completion_date, deliveries.delivery_status").
        Where("deliveries.delivery_status != ? AND deliveries.delivery_id = ?", ds.DELIVERY_STATUS_DELETED, deliveryID).
        Scan(&delivery).Error; err != nil {
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
