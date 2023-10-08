package repository

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
