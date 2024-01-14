package repository

import (
	"errors"
	"time"

	"github.com/markgregr/RIP/internal/model"
)

type DeliveryUserRepository interface{
    GetDeliveriesUser(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus string, userID uint) ([]model.DeliveryRequest, error)
    GetDeliveryByIDUser(deliveryID, userID uint) (model.DeliveryGetResponse, error)
    DeleteDeliveryUser(deliveryID, userID uint) error
    UpdateFlightNumberUser(deliveryID uint, userID uint, flightNumber model.DeliveryUpdateFlightNumberRequest) error
    UpdateDeliveryStatusUser(deliveryID, userID uint) error
}

func (r *Repository) GetDeliveriesUser(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus string, userID uint) ([]model.DeliveryRequest, error) {
    query := r.db.Table("deliveries").
        Select("DISTINCT deliveries.delivery_id, deliveries.flight_number, deliveries.creation_date, deliveries.formation_date, deliveries.completion_date, deliveries.delivery_status, users.full_name").
        Joins("JOIN users ON users.user_id = deliveries.user_id").
        Where("deliveries.delivery_status LIKE ? AND deliveries.flight_number LIKE ? AND deliveries.user_id = ? AND deliveries.delivery_status != ?", deliveryStatus, searchFlightNumber, userID, model.DELIVERY_STATUS_DELETED)
    
    if startFormationDate != "" && endFormationDate != "" {
        query = query.Where("deliveries.formation_date BETWEEN ? AND ?", startFormationDate, endFormationDate)
    }

    var deliveries []model.DeliveryRequest
    if err := query.Find(&deliveries).Error; err != nil {
        return nil, errors.New("ошибка получения доставок")
    }

    return deliveries, nil
}

func (r *Repository) GetDeliveryByIDUser(deliveryID, userID uint) (model.DeliveryGetResponse, error) {
    var delivery model.DeliveryGetResponse

    if err := r.db.
        Table("deliveries").
        Select("deliveries.delivery_id, deliveries.flight_number, deliveries.creation_date, deliveries.formation_date, deliveries.completion_date, deliveries.delivery_status, users.full_name").
        Joins("JOIN users ON users.user_id = deliveries.user_id").
        Where("deliveries.delivery_status != ? AND deliveries.delivery_id = ? AND deliveries.user_id = ?", model.DELIVERY_STATUS_DELETED, deliveryID, userID).
        First(&delivery).Error; err != nil {
        return model.DeliveryGetResponse{}, errors.New("ошибка получения доставки по ИД")
    }
    

    var baggages []model.Baggage
    if err := r.db.
        Table("baggages").
        Joins("JOIN delivery_baggages ON baggages.baggage_id = delivery_baggages.baggage_id").
        Where("delivery_baggages.delivery_id = ?", delivery.DeliveryID).
        Scan(&baggages).Error; err != nil {
        return model.DeliveryGetResponse{}, errors.New("ошибка получения багажей для доставки")
    }

    delivery.Baggages = baggages

    return delivery, nil
}

func (r *Repository) DeleteDeliveryUser(deliveryID, userID uint) error {
    var delivery model.Delivery
    if err := r.db.Table("deliveries").
        Where("delivery_id = ? AND user_id = ? AND delviery_status = ?", deliveryID, userID, model.DELIVERY_STATUS_DRAFT).
        First(&delivery).
        Error; err != nil {
        return errors.New("доставка не найдена или не принадлежит указанному пользователю или не находится в статусе черновик")
    }

    tx := r.db.Begin()
    if err := tx.Where("delivery_id = ?", deliveryID).Delete(&model.DeliveryBaggage{}).Error; 
    err != nil {
        tx.Rollback()
        return errors.New("ошибка удаления связей из таблицы-множества")
    }

    err := r.db.Model(&model.Delivery{}).Where("delivery_id = ?", deliveryID).Update("delivery_status", model.DELIVERY_STATUS_DELETED).Error
    if err != nil {
        return errors.New("ошибка обновления статуса на удален")
    }
     tx.Commit()
     
    return nil
}

func (r *Repository) UpdateFlightNumberUser(deliveryID uint, userID uint, flightNumber model.DeliveryUpdateFlightNumberRequest) error {
    var delivery model.Delivery
    if err := r.db.Table("deliveries").
        Where("delivery_id = ? AND user_id = ? AND delivery_status = ?", deliveryID, userID, model.DELIVERY_STATUS_DRAFT).
        First(&delivery).
        Error; err != nil {
        return errors.New("доставка не найдена или не принадлежит указанному пользователю")
    }

    if err := r.db.Table("deliveries").
        Model(&delivery).
        Update("flight_number", flightNumber.FlightNumber).
        Error; err != nil {
        return errors.New("ошибка обновления номера рейса")
    }

    return nil
}

func (r *Repository) UpdateDeliveryStatusUser(deliveryID, userID uint) error {
    var delivery model.Delivery
    if err := r.db.Table("deliveries").
        Where("delivery_id = ? AND user_id = ? AND delivery_status = ?", deliveryID, userID, model.DELIVERY_STATUS_DRAFT).
        First(&delivery).
        Error; err != nil {
        return errors.New("доставка не найдена, или не принадлежит указанному пользователю, или не имеет статус черновик")
    }

    delivery.DeliveryStatus = model.DELIVERY_STATUS_WORK
    currentTime := time.Now()
	delivery.FormationDate = &currentTime

    if err := r.db.Save(&delivery).Error; err != nil {
        return errors.New("ошибка обновления статуса доставки в БД")
    }

    return nil
}




