package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/markgregr/RIP/internal/app/ds"
)

func (r *Repository) GetDeliveriesForUser(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus string, userID uint) ([]ds.DeliveryRequest, error) {
    searchFlightNumber = strings.ToUpper(searchFlightNumber + "%")
    deliveryStatus = strings.ToLower(deliveryStatus + "%")

    // Построение основного запроса для получения доставок.
    query := r.db.Table("deliveries").
        Select("DISTINCT deliveries.delivery_id, deliveries.flight_number, deliveries.creation_date, deliveries.formation_date, deliveries.completion_date, deliveries.delivery_status, users.full_name").
        Joins("JOIN delivery_baggages ON deliveries.delivery_id = delivery_baggages.delivery_id").
        Joins("JOIN baggages ON baggages.baggage_id = delivery_baggages.baggage_id").
        Joins("JOIN users ON users.user_id = deliveries.user_id").
        Where("deliveries.delivery_status LIKE ? AND deliveries.flight_number LIKE ? AND deliveries.user_id = ? AND deliveries.delivery_status != ?", deliveryStatus, searchFlightNumber, userID, ds.DELIVERY_STATUS_DELETED)
    
    // Добавление условия фильтрации по дате формирования, если она указана.
    if startFormationDate != "" && endFormationDate != "" {
        query = query.Where("deliveries.formation_date BETWEEN ? AND ?", startFormationDate, endFormationDate)
    }

    // Выполнение запроса и сканирование результатов в слайс deliveries.
    var deliveries []ds.DeliveryRequest
    if err := query.Find(&deliveries).Error; err != nil {
        return nil, errors.New("ошибка получения доставок")
    }
    return deliveries, nil
}


func (r *Repository) GetDeliveryByIDForUser(deliveryID int, userID uint) (map[string]interface{}, error) {
    var delivery map[string]interface{}
    // Получение информации о доставке по deliveryID.
    if err := r.db.
        Table("deliveries").
        Select("deliveries.delivery_id, deliveries.flight_number, deliveries.creation_date, deliveries.formation_date, deliveries.completion_date, deliveries.delivery_status").
        Where("deliveries.delivery_status != ? AND deliveries.delivery_id = ? AND deliveries.user_id = ?", ds.DELIVERY_STATUS_DELETED, deliveryID, userID).
        Scan(&delivery).Error; err != nil {
        return nil, errors.New("ошибка получения доставки по ИД")
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

func (r *Repository) DeleteDeliveryForUser(deliveryID int, userID uint) error {
    // Проверяем, существует ли указанная доставка в базе данных
    var delivery ds.Delivery
    if err := r.db.First(&delivery, deliveryID).Error; err != nil {
        return errors.New("данная доставка не существует")
    }

    // Проверяем, что пользователь является создателем этой доставки
    if delivery.UserID != userID {
        return errors.New("пользователь не является создателем этой доставки")
    }

    // Начинаем транзакцию для атомарности операций
    tx := r.db.Begin()

    // Удаляем связанные записи из таблицы-множества (delivery_baggages)
    if err := tx.Where("delivery_id = ?", deliveryID).Delete(&ds.DeliveryBaggage{}).Error; err != nil {
        tx.Rollback()
        return errors.New("ошибка удаления связей из таблицы-множества")
    }

    // Обновляем статус доставки на "удален" с использованием GORM
    err := r.db.Model(&ds.Delivery{}).Where("delivery_id = ?", deliveryID).Update("delivery_status", ds.DELIVERY_STATUS_DELETED).Error
    if err != nil {
        return errors.New("ошибка обновления статуса на удален")
    }
     // Фиксируем транзакцию
     tx.Commit()
     
    return nil
}

func (r *Repository) UpdateDeliveryForUser(deliveryID int, userID uint, updatedDelivery *ds.Delivery) error {
    // Проверяем, существует ли указанная доставка в базе данных
    var delivery ds.Delivery
    if err := r.db.First(&delivery, deliveryID).Error; err != nil {
        return errors.New("данная доставка не существует")
    }

    // Проверяем, что доставка принадлежит указанному пользователю
    if delivery.UserID != userID {
        return errors.New("пользователь не является создателем этой доставки")
    }

    // Проверяем, что обновляем только поле FlightNumber
    if updatedDelivery.FlightNumber != "" {
        // Обновляем только поле FlightNumber из JSON-запроса
        if err := r.db.Model(&ds.Delivery{}).Where("delivery_id = ?", deliveryID).Update("flight_number", updatedDelivery.FlightNumber).Error; err != nil {
            return errors.New("ошибка обновления номера рейса")
        }
    } else {
        return errors.New("можно обновлять только номер рейса")
    }

    return nil
}

func (r *Repository) UpdateDeliveryStatusForUser(deliveryID int, userID uint) error {
    // Проверяем, существует ли указанная доставка в базе данных
    var delivery ds.Delivery
    if err := r.db.First(&delivery, deliveryID).Error; err != nil {
        return errors.New("данная доставка не существует")
    }

    // Проверяем, что пользователь имеет право на изменение статуса этой доставки
    if delivery.UserID != userID {
        return errors.New("пользователь не является создателем этой доставки")
    }

    // Проверяем, что текущий статус доставки - "черновик"
    if delivery.DeliveryStatus == ds.DELIVERY_STATUS_DRAFT {
        // Обновляем статус доставки на "в работе"
        delivery.DeliveryStatus = ds.DELIVERY_STATUS_WORK

        // Обновляем дату формирования на текущее московское время
        moscowTime, err := time.LoadLocation("Europe/Moscow")
        if err != nil {
            return err
        }
        delivery.FormationDate = time.Now().In(moscowTime)
    }else {
        return errors.New("доставка должна иметь статус черновик")
    }

    // Обновляем доставку в базе данных
    if err := r.db.Save(&delivery).Error; err != nil {
        return errors.New("ошибка обновления статуса")
    }

    return nil
}





