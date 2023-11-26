package repository

import (
	"errors"
	"time"

	"github.com/markgregr/RIP/internal/model"
)

type BaggageRepository interface {
    GetBaggages(searchCode string, userID uint) (model.BaggagesGetResponse, error)
}

func (r *Repository) GetBaggages(searchCode string, userID uint) (model.BaggagesGetResponse, error) {
    var deliveryID uint
    if err := r.db.
        Table("deliveries").
        Select("deliveries.delivery_id").
        Where("user_id = ? AND delivery_status = ?", userID, model.DELIVERY_STATUS_DRAFT).
        Take(&deliveryID).Error; 
        err != nil {
    }

    var baggages []model.Baggage
    if err := r.db.Table("baggages").
        Where("baggages.baggage_status = ? AND baggages.baggage_code LIKE ?", model.BAGGAGE_STATUS_ACTIVE, searchCode).
        Scan(&baggages).Error; err != nil {
        return model.BaggagesGetResponse{}, errors.New("ошибка нахождения списка багажа")
    }

    baggageResponse := model.BaggagesGetResponse{
        Baggages:   baggages,
        DeliveryID: deliveryID,
    }

    return baggageResponse, nil
}

func (r *Repository) GetBaggageByID(baggageID, userID uint) (model.Baggage, error) {
    var baggage model.Baggage

	if err := r.db.Table("baggages").
    Where("baggage_status = ? AND baggage_id = ?", model.BAGGAGE_STATUS_ACTIVE, baggageID).
    First(&baggage).Error; 
    err != nil {
		return model.Baggage{}, errors.New("ошибка при получении активного багажа из БД")
	}

    return baggage, nil
}

func (r *Repository) CreateBaggage(userID uint, baggage model.Baggage) error {
	if err := r.db.Create(baggage).Error; err != nil {
		return errors.New("ошибка создания багажа")
	}

	return nil
}

func (r *Repository) DeleteBaggage(baggageID, userID uint) error {
    var baggage model.Baggage

	if err := r.db.Table("baggages").Where("baggage_id = ? AND baggage_status = ?", baggageID, model.BAGGAGE_STATUS_ACTIVE).First(baggage).Error; 
    err != nil {
		return errors.New("багаж не найден или уже удален")
	}

	baggage.BaggageStatus = model.BAGGAGE_STATUS_DELETED

	if err := r.db.Table("baggages").Save(baggage).Error; 
    err != nil {
		return errors.New("ошибка при обновлении статуса багажа в БД")
	}
    return nil
}

func (r *Repository) UpdateBaggage(baggageID, userID uint, baggage model.Baggage) error {
    if err := r.db.Table("baggages").
    Model(&model.Baggage{}).
    Where("baggage_id = ? AND baggage_status = ?", baggageID, model.BAGGAGE_STATUS_ACTIVE).
    Updates(baggage).Error; 
    err != nil {
		return errors.New("ошибка при обновлении информации о питомце в БД")
	}

	return nil
}

func (r *Repository) AddBaggageToDelivery(baggageID, userID, moderatorID uint) error {
    var baggage model.Baggage

	if err := r.db.Table("baggages").
    Where("baggage_id = ? AND baggage_status = ?", baggageID, model.BAGGAGE_STATUS_ACTIVE).
    First(baggage).Error; 
    err != nil {
		return errors.New("багаж не найден или удален")
	}

    var delivery model.Delivery

    if err := r.db.Table("delivery").
    Where("delivery_status = ? AND user_id = ?", model.DELIVERY_STATUS_DRAFT, userID).
    Last(delivery).Error; 
    err != nil {
        delivery = model.Delivery{
            DeliveryStatus: model.DELIVERY_STATUS_DRAFT,
            CreationDate:   time.Now(),
            UserID:         userID, 
            ModeratorID:    moderatorID,
        }

        if err := r.db.Table("delivery").
        Create(delivery).Error;
        err != nil {
            return errors.New("ошибка создания доставки со статусом черновик")
        }
    }

    deliveryBaggage := model.DeliveryBaggage{
        BaggageID:  baggageID,
        DeliveryID: delivery.DeliveryID,
    }

    if err := r.db.Table("deliveries_baggages").
    Create(deliveryBaggage).Error;
    err != nil {
		return errors.New("ошибка при создании связи между доставкой и багажом")
	}

    return nil
}

func (r *Repository) RemoveBaggageFromDelivery(baggageID, userID uint) error {
    var deliveryBaggage model.DeliveryBaggage

    if err := r.db.Joins("JOIN deliveries ON delivery_baggages.delivery_id = deliveries.delivery_id").
        Where("delivery_baggages.baggage_id = ? AND deliveries.user_id = ? AND deliveries.delivery_status = ?", baggageID, userID, model.DELIVERY_STATUS_DRAFT).
        First(deliveryBaggage).Error; 
        err != nil {
        return errors.New("багаж не принадлежит пользователю или находится не в статусе черновик")
    }

    if err := r.db.Table("deliveries_baggages").
    Delete(deliveryBaggage).Error; 
    err != nil {
        return errors.New("ошибка удаления связи между доставкой и багажом")
    }
 
    return nil
}

func (r *Repository) AddBaggageImage(baggageID, userID uint, imageURL string) error {
    err := r.db.Table("baggages").Where("baggage_id = ?", baggageID).Update("photo", imageURL).Error
    if err != nil {
        return errors.New("ошибка обновления url изображения в БД")
    }

    return nil
}




