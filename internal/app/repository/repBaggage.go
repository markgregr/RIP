package repository



//методы для таблицы baggage
func (r *Repository) GetBaggages() ([]ds.Baggage, error) {
	searchCode = strings.ToUpper(searchCode+"%")
	var baggages []ds.Baggage
	if err := r.db.Find(&baggages, "baggage_status = ? AND baggage_code LIKE ?", ds.BAGGAGE_STATUS_ACTIVE, searchCode).Error; err != nil {
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