package repository

import (
	"errors"
	"strings"

	"github.com/markgregr/RIP/internal/app/ds"
)

func (r *Repository) GetDeliveries(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus string) ([]map[string]interface{}, error) {
    searchFlightNumber = strings.ToUpper(searchFlightNumber + "%")
    deliveryStatus = strings.ToLower(deliveryStatus + "%")

    // Построение основного запроса для получения доставок.
    query := r.db.Table("deliveries").
        Select("DISTINCT deliveries.flight_number, deliveries.creation_date, deliveries.formation_date, deliveries.completion_date, deliveries.delivery_status").
        Joins("JOIN delivery_baggages ON deliveries.delivery_id = delivery_baggages.delivery_id").
        Joins("JOIN baggages ON baggages.baggage_id = delivery_baggages.baggage_id").
        Where("deliveries.delivery_status LIKE ? AND deliveries.flight_number LIKE ?", deliveryStatus, searchFlightNumber)
    // Добавление условия фильтрации по дате формирования, если она указана.
    if startFormationDate != "" && endFormationDate != "" {
        query = query.Where("deliveries.formation_date BETWEEN ? AND ?", startFormationDate, endFormationDate)
    }

    // Выполнение запроса и сканирование результатов в структуру deliveries.
    var deliveries []map[string]interface{}
    if err := query.Scan(&deliveries).Error; err != nil {
        return nil, err
    }

    // Для каждой доставки получаем информацию о багаже по номеру рейса и статусу доставки.
    for _, delivery := range deliveries {
        baggages, err := r.GetBaggagesByFlightNumber(delivery["flight_number"].(string))
        if err != nil {
            return nil, err
        }
        // Добавляем информацию о багаже в поле "baggages" внутри каждой доставки.
        delivery["baggages"] = baggages
    }

    return deliveries, nil
}

func (r *Repository) GetBaggagesByFlightNumber(flightNumber string) ([]map[string]interface{}, error) {
    // Инициализация списка багажа для данного рейса и статуса доставки.
    var baggages []map[string]interface{}
    // Выполнение запроса к базе данных для получения багажа с указанными параметрами.
    if err := r.db.
        Table("baggages").
        Select("baggage_code, weight, size, baggage_status, baggage_type, owner_name, pasport_details, airline, photo_url").
        Joins("JOIN delivery_baggages ON baggages.baggage_id = delivery_baggages.baggage_id").
        Joins("JOIN deliveries ON delivery_baggages.delivery_id = deliveries.delivery_id").
        Where("baggages.baggage_status = ? AND deliveries.flight_number = ?", ds.BAGGAGE_STATUS_ACTIVE, flightNumber).
        Scan(&baggages).Error; err != nil {
        return nil, err
    }

    return baggages, nil
}


func (r *Repository) GetDeliveryByID(deliveryID int) (map[string]interface{}, error) {
    var delivery map[string]interface{}
    // Получение информации о доставке по deliveryID.
    if err := r.db.
        Table("deliveries").
        Select("deliveries.flight_number, deliveries.creation_date, deliveries.formation_date, deliveries.completion_date, deliveries.delivery_status").
        Where("deliveries.delivery_status != ? AND deliveries.delivery_id = ?", ds.DELIVERY_STATUS_DELETED, deliveryID).
        Scan(&delivery).Error; err != nil {
        return nil, err
    }

    // Получение багажей по указанному deliveryID.
    baggages, err := r.GetBaggagesByFlightNumber(delivery["flight_number"].(string))
    if err != nil {
        return nil, err
    }
    // Добавление информации о багаже в поле "baggages" внутри доставки.
    delivery["baggages"] = baggages

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
func (r *Repository) UpdateDeliveryStatusForUser(deliveryID int, userID uint, updateRequest *ds.Delivery) error {
    // Проверяем, существует ли указанная доставка в базе данных
    var delivery ds.Delivery
    if err := r.db.First(&delivery, deliveryID).Error; err != nil {
        return err
    }

    // Проверяем, что пользователь имеет право на изменение статуса этой доставки
    if delivery.UserID != userID {
        return errors.New("unauthorized: User does not have permission to update this delivery status")
    }

    // Обновляем только поле DeliveryStatus из JSON-запроса
    delivery.DeliveryStatus = updateRequest.DeliveryStatus

    // Обновляем доставку в базе данных
    if err := r.db.Save(&delivery).Error; err != nil {
        return err
    }

    return nil
}

func (r *Repository) UpdateDeliveryStatusForModerator(deliveryID int, moderatorID uint, updateRequest *ds.Delivery) error {
    // Проверяем, существует ли указанная доставка в базе данных
    var delivery ds.Delivery
    if err := r.db.First(&delivery, deliveryID).Error; err != nil {
        return err
    }

    // Проверяем, что модератор имеет право на изменение статуса этой доставки
    if delivery.ModeratorID != moderatorID {
        return errors.New("unauthorized: Moderator does not have permission to update this delivery status")
    }

    // Обновляем только поле DeliveryStatus из JSON-запроса
    delivery.DeliveryStatus = updateRequest.DeliveryStatus

    // Обновляем доставку в базе данных
    if err := r.db.Save(&delivery).Error; err != nil {
        return err
    }

    return nil
}



