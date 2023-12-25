package repository

import (
	"errors"
	"time"

	"github.com/markgregr/RIP/internal/model"
)


func (r *Repository) GetDeliveriesModerator(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus string, moderatorID uint) ([]model.DeliveryRequest, error) {
    query := r.db.Table("deliveries").
        Select("DISTINCT deliveries.delivery_id, deliveries.flight_number, deliveries.creation_date, deliveries.formation_date, deliveries.completion_date, deliveries.delivery_status, users.full_name").
        Joins("JOIN users ON users.user_id = deliveries.user_id").
        Where("deliveries.delivery_status LIKE ? AND deliveries.flight_number LIKE ? AND deliveries.delivery_status != ?", deliveryStatus, searchFlightNumber, model.DELIVERY_STATUS_DELETED)
    
    if startFormationDate != "" && endFormationDate != "" {
        query = query.Where("deliveries.formation_date BETWEEN ? AND ?", startFormationDate, endFormationDate)
    }

    var deliveries []model.DeliveryRequest
    if err := query.Find(&deliveries).Error; err != nil {
        return nil, errors.New("ошибка получения доставок")
    }


    return deliveries, nil
}

func (r *Repository) GetDeliveryByIDModerator(deliveryID, moderatorID uint) (model.DeliveryGetResponse, error) {
    var delivery model.DeliveryGetResponse

    if err := r.db.
        Table("deliveries").
        Select("deliveries.delivery_id, deliveries.flight_number, deliveries.creation_date, deliveries.formation_date, deliveries.completion_date, deliveries.delivery_status, users.full_name").
        Joins("JOIN users ON users.user_id = deliveries.user_id").
        Where("deliveries.delivery_status != ? AND deliveries.delivery_id = ?", model.DELIVERY_STATUS_DELETED, deliveryID).
        Scan(&delivery).Error; err != nil {
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

func (r *Repository) UpdateFlightNumberModerator(deliveryID uint, moderatorID uint, flightNumber model.DeliveryUpdateFlightNumberRequest) error {
    var delivery model.Delivery
    if err := r.db.Table("deliveries").
        Where("delivery_id = ? AND moderator_id = ?", deliveryID, moderatorID).
        First(&delivery).
        Error; err != nil {
        return errors.New("доставка не найдена или не принадлежит указанному модератору")
    }

    if err := r.db.Table("deliveries").
        Model(&delivery).
        Update("flight_number", flightNumber.FlightNumber).
        Error; err != nil {
        return errors.New("ошибка обновления номера рейса")
    }

    return nil
}

func (r *Repository) UpdateDeliveryStatusModerator(deliveryID, moderatorID uint, deliveryStatus model.DeliveryUpdateStatusRequest) error {
    var delivery model.Delivery
    if err := r.db.Table("deliveries").
        Where("delivery_id = ? AND moderator_id = ? AND delivery_status = ?", deliveryID, moderatorID, model.DELIVERY_STATUS_WORK).
        First(&delivery).
        Error; err != nil {
        return errors.New("доставка не найдена или не принадлежит указанному модератору")
    }

    delivery.DeliveryStatus = deliveryStatus.DeliveryStatus
	delivery.CompletionDate = time.Now()

    if err := r.db.Save(&delivery).Error; err != nil {
        return errors.New("ошибка обновления статуса доставки в БД")
    }

    return nil
}