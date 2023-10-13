package repository

import (
	"strings"

	"github.com/markgregr/RIP/internal/app/ds"
)

//методы для таблицы baggage
func (r *Repository) GetBaggages(searchCode string) ([]map[string]interface{}, error) {
    searchCode = strings.ToUpper(searchCode + "%")
    var baggages []map[string]interface{}
    if err := r.db.
        Table("baggages").
        Select("baggages.baggage_code, baggages.weight, baggages.size, baggages.baggage_status, baggages.baggage_type, baggages.owner_name, baggages.pasport_details, baggages.airline, deliveries.flight_number, deliveries.delivery_status").
        Joins("JOIN delivery_baggages ON baggages.baggage_id = delivery_baggages.baggage_id").
        Joins("JOIN deliveries ON deliveries.delivery_id = delivery_baggages.delivery_id").
        Where("baggages.baggage_status = ? AND baggages.baggage_code LIKE ?", ds.BAGGAGE_STATUS_ACTIVE, searchCode).
        Scan(&baggages).Error; err != nil {
        return nil, err
    }

    return baggages, nil
}

func (r *Repository) GetBaggageByID(baggageID int) (map[string]interface{}, error) {
    var baggage map[string]interface{}
    if err := r.db.
        Table("baggages").
        Select("baggages.baggage_code, baggages.weight, baggages.size, baggages.baggage_status, baggages.baggage_type, baggages.owner_name, baggages.pasport_details, baggages.airline").
        Where("baggages.baggage_status = ? AND baggages.baggage_id = ?", ds.BAGGAGE_STATUS_ACTIVE, baggageID).
        Scan(&baggage).Error; err != nil {
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


