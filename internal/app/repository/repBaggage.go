package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/markgregr/RIP/internal/app/ds"
)

func (r *Repository) GetBaggages(searchCode string, userID uint) (map[string]interface{}, error) {
    searchCode = strings.ToUpper(searchCode + "%")
    var deliveryID uint
    if err := r.db.
        Table("deliveries").
        Select("deliveries.delivery_id").
        Where("user_id = ? AND delivery_status = ?", userID, ds.DELIVERY_STATUS_DRAFT).
        Take(&deliveryID).Error; err != nil {
        return nil, errors.New("ошибка нахождения delivery_id черновика")
    }

    var baggages []map[string]interface{}
    if err := r.db.
        Table("baggages").
        Select("baggages.baggage_id, baggages.baggage_code, baggages.weight, baggages.size, baggages.baggage_status, baggages.baggage_type, baggages.owner_name, baggages.pasport_details, baggages.airline, baggages.photo_url").
        Where("baggages.baggage_status = ? AND baggages.baggage_code LIKE ?", ds.BAGGAGE_STATUS_ACTIVE, searchCode).
        Scan(&baggages).Error; err != nil {
        return nil, errors.New("ошибка нахождения списка багажа")
    }

    // Создаем объект JSON для включения delivery_id и baggages
    result := make(map[string]interface{})
    result["delivery_id"] = deliveryID
    result["baggages"] = baggages

    return result, nil
}

func (r *Repository) GetBaggageByID(baggageID int, userID uint) (map[string]interface{}, error) {
    var baggage map[string]interface{}
    if err := r.db.
        Table("baggages").
        Select("baggages.baggage_id, baggages.baggage_code, baggages.weight, baggages.size, baggages.baggage_status, baggages.baggage_type, baggages.owner_name, baggages.pasport_details, baggages.airline, baggages.photo_url").
        Where("baggages.baggage_status = ? AND baggages.baggage_id = ?", ds.BAGGAGE_STATUS_ACTIVE, baggageID).
        Scan(&baggage).Error; err != nil {
        return nil, errors.New("ошибка нахождения багажа по ID")
    }
    return baggage, nil
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
        return nil, errors.New("ошибка нахождения списка багажа по номеру рейса")
    }

    return baggages, nil
}

func (r *Repository) CreateBaggage(baggage *ds.Baggage) error {
	// Создаем багаж
	if err := r.db.Create(baggage).Error; err != nil {
		return errors.New("ошибка создания багажа")
	}

	return nil
}

func (r *Repository) DeleteBaggage(baggageID int, userID uint) error {
    // Удаление изображения из MinIO
    err := r.minioClient.RemoveServiceImage(baggageID)
    if err != nil {
        // Обработка ошибки удаления изображения из MinIO, если необходимо
        return err
    }
	return r.db.Exec("UPDATE baggages SET baggage_status = ? WHERE baggage_id = ?", ds.BAGGAGE_STATUS_DELETED, baggageID).Error
}

func (r *Repository) UpdateBaggage(baggageID int, updatedBaggage *ds.Baggage) error {
	err := r.db.Model(&ds.Baggage{}).Where("baggage_id = ? AND baggage_status = ?", baggageID, ds.BAGGAGE_STATUS_ACTIVE).Updates(updatedBaggage).Error
	if err != nil {
		return errors.New("ошибка изменения багажа")
	}
	return nil
}

func (r *Repository) AddBaggageToDelivery(baggageID uint, userID uint, moderatorID uint) error {
    // Проверяем, существует ли указанный багаж в базе данных
    var baggage ds.Baggage
    if err := r.db.First(&baggage, baggageID).Error; err != nil {
        return errors.New("недопустимый ID для багажа")
    }

    // Получаем последнюю заявку со статусом "черновик" для указанного пользователя, если такая существует
    var latestDraftDelivery ds.Delivery
    if err := r.db.Where("delivery_status = ? AND user_id = ?", ds.DELIVERY_STATUS_DRAFT, userID).Last(&latestDraftDelivery).Error; err != nil {
        // Если нет заявки со статусом "черновик", создаем новую
        currentTime := time.Now().In(time.FixedZone("UTC+3", 3*60*60)) // Часовой пояс Москвы
        latestDraftDelivery = ds.Delivery{
            DeliveryStatus: ds.DELIVERY_STATUS_DRAFT,
            CreationDate:   currentTime,
            UserID:         userID, // Устанавливаем ID пользователя для заявки
            ModeratorID:    moderatorID,
        }
        if err := r.db.Create(&latestDraftDelivery).Error; err != nil {
            return errors.New("ошибка создания доставки со статусом черновик")
        }
    }

    // Создаем связь между багажом и заявкой в промежуточной таблице
    relation := &ds.DeliveryBaggage{
        BaggageID:  baggageID,
        DeliveryID: latestDraftDelivery.DeliveryID,
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
        return errors.New("ошибка создания связи между доставкой и багажом")
    }

    // Фиксируем транзакцию
    tx.Commit()

    return nil
}

func (r *Repository) RemoveBaggageFromDelivery(baggageID uint, userID uint) error {
    // Начинаем транзакцию
    tx := r.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // Поиск связи между багажом и доставкой в базе данных
    var relation ds.DeliveryBaggage

    // Проверяем, принадлежит ли багаж текущему пользователю и находится ли он в статусе "черновик"
    if err := tx.Joins("JOIN deliveries ON delivery_baggages.delivery_id = deliveries.delivery_id").
        Where("delivery_baggages.baggage_id = ? AND deliveries.user_id = ? AND deliveries.delivery_status = ?", baggageID, userID, ds.DELIVERY_STATUS_DRAFT).
        First(&relation).Error; err != nil {
        tx.Rollback()
        return errors.New("багаж не принадлежит пользователю или находится не в статусе черновик")
    }

    // Удаление связи из базы данных
    if err := tx.Delete(&relation).Error; err != nil {
        tx.Rollback()
        return errors.New("ошибка удаления связи между доставкой и багажом")
    }

    // Фиксируем транзакцию
    tx.Commit()

    return nil
}


func (r *Repository) AddBaggageImage(baggageID int, imageBytes []byte, contentType string) error {
    // Удаление существующего изображения (если есть)
    err := r.minioClient.RemoveServiceImage(baggageID)
    if err != nil {
        return err
    }

    // Загрузка нового изображения в MinIO
    imageURL, err := r.minioClient.UploadServiceImage(baggageID, imageBytes, contentType)
    if err != nil {
        return err
    }

    // Обновление информации об изображении в БД (например, ссылки на MinIO)
    err = r.db.Model(&ds.Baggage{}).Where("baggage_id = ?", baggageID).Update("photo_url", imageURL).Error
    if err != nil {
        // Обработка ошибки обновления URL изображения в БД, если необходимо
        return errors.New("ошибка обновления url изображения в БД")
    }

    return nil
}




