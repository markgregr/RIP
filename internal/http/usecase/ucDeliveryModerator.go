package usecase

import (
	"errors"
	"strings"

	"github.com/markgregr/RIP/internal/model"
)

type DeliveryModeratorUseCase interface {
	GetDeliveriesModerator(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus string) ([]model.DeliveryRequest, error)
	GetDeliveryByIDModerator(deliveryID uint) (model.DeliveryGetResponse, error)
	UpdateDeliveryStatusModerator(deliveryID, moderatorID uint, deliveryStatus model.DeliveryUpdateStatusRequest) error
}

func (uc *UseCase) GetDeliveriesModerator(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus string) ([]model.DeliveryRequest, error) {
	searchFlightNumber = strings.ToUpper(searchFlightNumber + "%")
	deliveryStatus = strings.ToLower(deliveryStatus + "%")

	deliveries, err := uc.Repository.GetDeliveriesModerator(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus)
	if err != nil {
		return nil, err
	}

	return deliveries, nil
}

func (uc *UseCase) GetDeliveryByIDModerator(deliveryID uint) (model.DeliveryGetResponse, error) {
	if deliveryID <= 0 {
		return model.DeliveryGetResponse{}, errors.New("недопустимый ИД доставки")
	}

	deliveries, err := uc.Repository.GetDeliveryByIDModerator(deliveryID)
	if err != nil {
		return model.DeliveryGetResponse{}, err
	}

	return deliveries, nil
}

func (uc *UseCase) UpdateDeliveryStatusModerator(deliveryID, moderatorID uint, deliveryStatus model.DeliveryUpdateStatusRequest) error{
	if deliveryID <= 0 {
		return errors.New("недопустимый ИД доставки")
	}
	if deliveryStatus.DeliveryStatus != model.DELIVERY_STATUS_COMPLETED && deliveryStatus.DeliveryStatus != model.DELIVERY_STATUS_REJECTED {
		return errors.New("текущий статус доставки уже завершен или отклонен")
	}

	err := uc.Repository.UpdateDeliveryStatusModerator(deliveryID, moderatorID, deliveryStatus)
	if err != nil {
		return err
	}

	return nil
}


