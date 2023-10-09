package repository

import (
	"errors"

	"github.com/markgregr/RIP/internal/app/ds"
)

func (r *Repository) AddBaggageToDelivery(baggageID uint, deliveryID uint) error {
    // Проверяем, существуют ли указанный багаж и доставка в базе данных
    var baggage ds.Baggage
    if err := r.db.First(&baggage, baggageID).Error; err != nil {
        return errors.New("Baggage not found")
    }

    var delivery ds.Delivery
    if err := r.db.First(&delivery, deliveryID).Error; err != nil {
        return errors.New("Delivery not found")
    }

    // Проверяем, что багаж и доставка активны (или другие условия, которые вам нужны)
    if baggage.BaggageStatus != ds.BAGGAGE_STATUS_ACTIVE {
        return errors.New("Baggage is not active")
    }

	// Создаем связь между багажом и доставкой в промежуточной таблице
    relation := &ds.DeliveryBaggage{
        BaggageID:  baggageID,
        DeliveryID: deliveryID,
    }

    if err := r.db.Create(relation).Error; err != nil {
        return r.db.Exec("UPDATE deliveries SET delivery_status = ? WHERE delievry_id = ?", ds.DELIVERY_STATUS_WORK, deliveryID).Error
    }

    return nil
}

func (r *Repository) RemoveBaggageFromDelivery(baggageID uint, deliveryID uint) error {
    // Поиск связи между багажом и доставкой в базе данных
    var relation ds.DeliveryBaggage
    if err := r.db.Where("baggage_id = ? AND delivery_id = ?", baggageID, deliveryID).First(&relation).Error; err != nil {
        return errors.New("Relation not found")
    }

    // Удаление связи из базы данных
    if err :=r.db.Delete(&relation).Error; err != nil {
        return r.db.Exec("UPDATE deliveries SET delivery_status = ? WHERE delievry_id = ?", ds.DELIVERY_STATUS_REJECTED, deliveryID).Error;
    }

    return nil
}
