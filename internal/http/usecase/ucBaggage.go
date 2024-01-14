package usecase

import (
	"errors"
	"strings"

	"github.com/markgregr/RIP/internal/model"
)

type BaggageUseCase interface {
	GetBaggages(searchCode string, userID uint) (model.BaggagesGetResponse, error)
	GetBaggageByID(baggageID, userID uint) (model.Baggage, error)
	CreateBaggage(userID uint, requestBaggage model.BaggageRequest) error
	DeleteBaggage(baggageID, userID uint) error
	UpdateBaggage(baggageID, userID uint, requestBaggage model.BaggageRequest) error
	AddBaggageToDelivery(baggageID, userID, moderatorID uint) error
	RemoveBaggageFromDelivery(baggageID, userID uint) error
	AddBaggageImage(baggageID, userID uint, imageBytes []byte, ContentType string) error
}

func (uc *UseCase) GetBaggages(searchCode string, userID uint) (model.BaggagesGetResponse, error) {
	if userID <= 0 {
		return model.BaggagesGetResponse{}, errors.New("недопустимый ИД пользователя")
	}

	searchCode = strings.ToUpper(searchCode + "%")

	baggages, err := uc.Repository.GetBaggages(searchCode, userID)
	if err != nil {
		return model.BaggagesGetResponse{}, err
	}

	return baggages, nil
}

func (uc *UseCase) GetBaggageByID(baggageID, userID uint) (model.Baggage, error) {
	if baggageID <= 0 {
		return model.Baggage{}, errors.New("недопустимый ИД багажа")
	}
	if userID <= 0 {
		return model.Baggage{}, errors.New("недопустимый ИД пользователя")
	}

	baggage, err := uc.Repository.GetBaggageByID(baggageID, userID)
	if err != nil {
		return model.Baggage{}, err
	}

	return baggage, nil
}

func (uc *UseCase) CreateBaggage(userID uint, requestBaggage model.BaggageRequest) error {
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}
	if requestBaggage.BaggageCode == "" {
		return errors.New("код багажа должен быть заполнен")
	}
	if requestBaggage.Weight == 0 {
		return errors.New("вес багажа должен быть заполнен")
	}
	if requestBaggage.Size == "" {
		return errors.New("размер багажа должен быть заполнен")
	}
	if requestBaggage.BaggageType == "" {
		return errors.New("тип багажа должен быть заполнен")
	}
	if requestBaggage.OwnerName == "" {
		return errors.New("владелец багажа должен быть заполнен")
	}
	if requestBaggage.PasportDetails == "" {
		return errors.New("паспортные данные владельца багажа должны быть заполнен")
	}
	if requestBaggage.Airline == "" {
		return errors.New("авиакомпания должна быть заполнена")
	}

	baggage := model.Baggage{
		BaggageCode: 	requestBaggage.BaggageCode,
		Weight:			requestBaggage.Weight,
		Size:			requestBaggage.Size,
		BaggageType:	requestBaggage.BaggageType,
		OwnerName:		requestBaggage.OwnerName,
		PasportDetails:	requestBaggage.PasportDetails,
		Airline:		requestBaggage.Airline,
		BaggageStatus:  model.BAGGAGE_STATUS_ACTIVE,
	}

	err := uc.Repository.CreateBaggage(userID, baggage)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) DeleteBaggage(baggageID, userID uint) error {
	if baggageID <= 0 {
		return errors.New("недопустимый ИД багажа")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}

	err := uc.Repository.DeleteBaggage(baggageID, userID)
	if err != nil {
		return err
	}

	err = uc.Repository.RemoveServiceImage(baggageID, userID)
	if err != nil {
		return err
	}
	
	return nil
}

func (uc *UseCase) UpdateBaggage(baggageID, userID uint, requestBaggage model.BaggageRequest) error {
	if baggageID <= 0 {
		return errors.New("недопустимый ИД багажа")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}

	baggage := model.Baggage{
		BaggageCode: 	requestBaggage.BaggageCode,
		Weight:			requestBaggage.Weight,
		Size:			requestBaggage.Size,
		BaggageType:	requestBaggage.BaggageType,
		OwnerName:		requestBaggage.OwnerName,
		PasportDetails:	requestBaggage.PasportDetails,
		Airline:		requestBaggage.Airline,
	}

	err := uc.Repository.UpdateBaggage(baggageID, userID, baggage)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) AddBaggageToDelivery(baggageID, userID uint) error {
	if baggageID <= 0 {
		return errors.New("недопустимый ИД багажа")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}

	err := uc.Repository.AddBaggageToDelivery(baggageID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) RemoveBaggageFromDelivery(baggageID, userID uint) error {
	if baggageID <= 0 {
		return errors.New("недопустимый ИД багажа")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}

	err := uc.Repository.RemoveBaggageFromDelivery(baggageID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) AddBaggageImage(baggageID, userID uint, imageBytes []byte, ContentType string) error {
	if baggageID <= 0 {
		return errors.New("недопустимый ИД багажа")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}
	if imageBytes == nil {
		return errors.New("недопустимый imageBytes изображения")
	}
	if ContentType == "" {
		return errors.New("недопустимый ContentType изображения")
	}

	imageURL, err := uc.Repository.UploadServiceImage(baggageID, userID, imageBytes, ContentType)
	if err != nil {
		return err
	}

	err = uc.Repository.AddBaggageImage(baggageID, userID, imageURL)
	if err != nil {
		return err
	}

	return nil
}





