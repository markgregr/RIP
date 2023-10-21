package repository

import (
	"strings"
	"time"

	"github.com/markgregr/RIP/internal/app/ds"
)

func (r *Repository) GetBaggages(searchCode string) ([]map[string]interface{}, error) {
    searchCode = strings.ToUpper(searchCode + "%")
    var baggages []map[string]interface{}
    if err := r.db.
        Table("baggages").
        Select("baggage_code, weight, size, baggage_status, baggage_type, owner_name, pasport_details, airline, photo_url").
        Where("baggage_status = ? AND baggage_code LIKE ?", ds.BAGGAGE_STATUS_ACTIVE, searchCode).
        Scan(&baggages).Error; err != nil {
        return nil, err
    }

    return baggages, nil
}
func (r *Repository) GetBaggageByID(baggageID int) (map[string]interface{}, error) {
    var baggage map[string]interface{}
    if err := r.db.
        Table("baggages").
        Select("baggage_code, weight, size, baggage_status, baggage_type, owner_name, pasport_details, airline, photo_url").
        Where("baggages.baggage_status = ? AND baggages.baggage_id = ?", ds.BAGGAGE_STATUS_ACTIVE, baggageID).
        Scan(&baggage).Error; err != nil {
        return nil, err
    }
    return baggage, nil
}
func (r *Repository) CreateBaggage(baggage *ds.Baggage) error {
	// Создаем багаж
	if err := r.db.Create(baggage).Error; err != nil {
		return err
	}

	return nil
}
func (r *Repository) DeleteBaggage(baggageID int) error {
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
		return err
	}
	return nil
}

func (r *Repository) AddBaggageToDelivery(baggageID uint, userID uint, moderatorID uint) error {
    // Проверяем, существует ли указанный багаж в базе данных
    var baggage ds.Baggage
    if err := r.db.First(&baggage, baggageID).Error; err != nil {
        return err
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
            return err
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
        return err
    }

    if err := tx.Where("baggage_id = ? AND delivery_id = ?", baggageID, deliveryID).First(&relation).Error; err != nil {
        tx.Rollback()
        return err
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
        return err
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
        return err
    }

    return nil
}


