package usecase

import (
	"errors"
	"strings"

	"github.com/markgregr/RIP/internal/model"
)

func (uc *UseCase) GetDeliveriesModerator(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus string, moderatorID uint) ([]model.DeliveryRequest, error) {
	searchFlightNumber = strings.ToUpper(searchFlightNumber + "%")
	deliveryStatus = strings.ToLower(deliveryStatus + "%")

	if moderatorID <= 0 {
		return nil, errors.New("недопустимый ИД модератора")
	}

	deliveries, err := uc.Repository.GetDeliveriesModerator(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus, moderatorID)
	if err != nil {
		return nil, err
	}

	return deliveries, nil
}

func (uc *UseCase) GetDeliveryByIDModerator(deliveryID, moderatorID uint) (model.DeliveryGetResponse, error) {
	if deliveryID <= 0 {
		return model.DeliveryGetResponse{}, errors.New("недопустимый ИД доставки")
	}
	if moderatorID <= 0 {
		return model.DeliveryGetResponse{}, errors.New("недопустимый ИД модератора")
	}

	deliveries, err := uc.Repository.GetDeliveryByIDModerator(deliveryID, moderatorID)
	if err != nil {
		return model.DeliveryGetResponse{}, err
	}

	return deliveries, nil
}

func (uc *UseCase) UpdateFlightNumberModerator(deliveryID, moderatorID uint, flightNumber model.DeliveryUpdateFlightNumberRequest) error{
	if deliveryID <= 0 {
		return errors.New("недопустимый ИД доставки")
	}
	if moderatorID <= 0 {
		return errors.New("недопустимый ИД модератора")
	}
	if len(flightNumber.FlightNumber) !=6 {
		return errors.New("недопустимый номер рейса")
	}

	err := uc.Repository.UpdateFlightNumberModerator(deliveryID, moderatorID, flightNumber)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) UpdateDeliveryStatusModerator(deliveryID, moderatorID uint, deliveryStatus model.DeliveryUpdateStatusRequest) error{
	if deliveryID <= 0 {
		return errors.New("недопустимый ИД доставки")
	}
	if moderatorID <= 0 {
		return errors.New("недопустимый ИД модератора")
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


