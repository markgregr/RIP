package usecase

import (
	"errors"
	"strings"

	"github.com/markgregr/RIP/internal/model"
)

func (uc *UseCase) GetDeliveriesUser(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus string, userID uint) ([]model.DeliveryRequest, error) {
	searchFlightNumber = strings.ToUpper(searchFlightNumber + "%")
	deliveryStatus = strings.ToLower(deliveryStatus + "%")

	if userID <= 0 {
		return nil, errors.New("недопустимый ИД пользователя")
	}

	deliveries, err := uc.Repository.GetDeliveriesUser(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus, userID)
	if err != nil {
		return nil, err
	}

	return deliveries, nil
}

func (uc *UseCase) GetDeliveryByIDUser(deliveryID, userID uint) (model.DeliveryGetResponse, error) {
	if deliveryID <= 0 {
		return model.DeliveryGetResponse{}, errors.New("недопустимый ИД доставки")
	}
	if userID <= 0 {
		return model.DeliveryGetResponse{}, errors.New("недопустимый ИД пользователя")
	}

	deliveries, err := uc.Repository.GetDeliveryByIDUser(deliveryID, userID)
	if err != nil {
		return model.DeliveryGetResponse{}, err
	}

	return deliveries, nil
}

func (uc *UseCase) DeleteDeliveryUser(deliveryID, userID uint) error{
	if deliveryID <= 0 {
		return errors.New("недопустимый ИД доставки")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}

	err := uc.Repository.DeleteDeliveryUser(deliveryID, userID)
	if err != nil {
		return err
	}

	return nil
}


func (uc *UseCase) UpdateFlightNumberUser(deliveryID, userID uint, flightNumber model.DeliveryUpdateFlightNumberRequest) error{
	if deliveryID <= 0 {
		return errors.New("недопустимый ИД доставки")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}
	if len(flightNumber.FlightNumber) !=6 {
		return errors.New("недопустимый номер рейса")
	}

	err := uc.Repository.UpdateFlightNumberUser(deliveryID, userID, flightNumber)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) UpdateDeliveryStatusUser(deliveryID, userID uint) error{
	if deliveryID <= 0 {
		return errors.New("недопустимый ИД доставки")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}

	err := uc.Repository.UpdateDeliveryStatusUser(deliveryID, userID)
	if err != nil {
		return err
	}

	return nil
}


