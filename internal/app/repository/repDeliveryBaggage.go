package repository

import (
	"errors"
	"time"

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

    // Устанавливаем formation_date для доставки на текущую дату и время в часовом поясе "Europe/Moscow"
    moscowLocation, err := time.LoadLocation("Europe/Moscow")
    if err != nil {
        return err
    }
    currentTime := time.Now().In(moscowLocation)

    // Создаем связь между багажом и доставкой в промежуточной таблице
    relation := &ds.DeliveryBaggage{
        BaggageID:  baggageID,
        DeliveryID: deliveryID,
    }

    // Начинаем транзакцию
    tx := r.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // Создаем связь в таблице delivery_baggages
    if err := tx.Create(relation).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Обновляем formation_date для доставки
    if err := tx.Model(&delivery).Update("formation_date", currentTime).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Фиксируем транзакцию
    tx.Commit()

    return nil
}

func (r *Repository) RemoveBaggageFromDelivery(baggageID uint, deliveryID uint) error {
    // Начинаем транзакцию
    tx := r.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // Поиск связи между багажом и доставкой в базе данных
    var relation ds.DeliveryBaggage
    var delivery ds.Delivery
    if err := r.db.First(&delivery, deliveryID).Error; err != nil {
        return errors.New("Delivery not found")
    }

    if err := tx.Where("baggage_id = ? AND delivery_id = ?", baggageID, deliveryID).First(&relation).Error; err != nil {
        tx.Rollback()
        return errors.New("Relation not found")
    }
    // Устанавливаем formation_date для доставки на текущую дату и время в часовом поясе "Europe/Moscow"
    moscowLocation, err := time.LoadLocation("Europe/Moscow")
    if err != nil {
        return err
    }   
    // Сохраняем текущее значение времени как новую formation_date
    currentTime := time.Now().In(moscowLocation)

    // Удаление связи из базы данных
    if err := tx.Delete(&relation).Error; err != nil {
        tx.Rollback()
        return errors.New("Relation not deleted")
    }

    // Обновляем formation_date для доставки на текущую дату и время
    if err := tx.Model(&delivery).Update("formation_date", currentTime).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Фиксируем транзакцию
    tx.Commit()

    return nil
}

