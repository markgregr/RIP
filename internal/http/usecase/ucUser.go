package usecase

import (
	"errors"

	"github.com/markgregr/RIP/internal/model"
	"github.com/markgregr/RIP/internal/pkg/middleware"
)

type UserUseCase interface {
	RegisterUser(requestUser model.User) (model.User, error)
	LoginUser(requestUser model.User) (model.User, error)
	GetUserById(userID uint) (model.User, error)
}

func (uc *UseCase) RegisterUser(requestUser model.UserRegisterRequest) (model.UserLoginResponse, error) {
	if requestUser.FullName == "" {
		return model.UserLoginResponse{}, errors.New("ФИО должно быть заполнен")
	}
	if requestUser.Email == "" {
		return model.UserLoginResponse{}, errors.New("почта должна быть заполнена")
	}
	if requestUser.Password == "" {
		return model.UserLoginResponse{}, errors.New("пароль должен быть заполнен")
	}
	if len(requestUser.Password) < 8 || len(requestUser.Password) > 20 {
		return model.UserLoginResponse{}, errors.New("пароль должен содержать от 8 до 20 символов")
	}

	candidate, err := uc.Repository.GetByEmail(requestUser.Email)
	if err != nil {
		return model.UserLoginResponse{}, err
	}

	if candidate.Email == requestUser.Email {
		return model.UserLoginResponse{}, errors.New("такой пользователь уже существует")
	}

	requestUser.Password, err = middleware.HashPassword(requestUser.Password)
	if err != nil {
		return model.UserLoginResponse{}, err
	}

	user := model.User{
		FullName:    requestUser.FullName,
		Email:       requestUser.Email,
		Password:    requestUser.Password,
		Role: 		 "пользователь",
	}

	err = uc.Repository.CreateUser(user)
	if err != nil {
		return model.UserLoginResponse{}, err
	}
	user, err = uc.Repository.GetByEmail(requestUser.Email)
	if err != nil {
		return model.UserLoginResponse{}, err
	}

	tokenPair, err := middleware.GenerateJWTTokenPair(uint(user.UserID))
	if err != nil {
		return model.UserLoginResponse{}, err
	}

	err = uc.Repository.SaveJWTTokenPair(uint(user.UserID), tokenPair.AccessToken, tokenPair.RefreshToken)
	if err != nil {
		return model.UserLoginResponse{}, err
	}
	response := model.UserLoginResponse{
		AccessToken: tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		FullName: user.FullName,
	}
	return response, nil
}

func (uc *UseCase) LoginUser(requestUser model.UserLoginRequest) (model.UserLoginResponse, error) {
	if requestUser.Email == "" {
		return model.UserLoginResponse{}, errors.New("заполните почту")
	}

	if requestUser.Password == "" {
		return model.UserLoginResponse{}, errors.New("заполните пароль")
	}

	candidate, err := uc.Repository.GetByEmail(requestUser.Email)
	if err != nil {
		return model.UserLoginResponse{}, err
	}

	if ok, err := middleware.CheckPasswordHash(requestUser.Password, candidate.Password); !ok {
		return model.UserLoginResponse{}, err
	}

	tokenPair, err := middleware.GenerateJWTTokenPair(uint(candidate.UserID))
	if err != nil {
		return model.UserLoginResponse{}, err
	}

	err = uc.Repository.SaveJWTTokenPair(uint(candidate.UserID), tokenPair.AccessToken, tokenPair.RefreshToken)
	if err != nil {
		return model.UserLoginResponse{}, err
	}
	response := model.UserLoginResponse{
		AccessToken: tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		FullName: candidate.FullName,
	}
	return response, nil
}

func (uc *UseCase) GetUserByID(userID uint) (model.User, error) {
	if userID < 1 {
		return model.User{}, errors.New("ID не может быть отрицательным")
	}

	user, err := uc.Repository.GetUserByID(userID)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (uc *UseCase) GetUsers() ([]model.User, error) {
	users, err := uc.Repository.GetUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (uc *UseCase) LogoutUser(userID uint) error{
	err := uc.Repository.DeleteJWTTokenPair(userID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) RefreshToken(refreshToken string) (model.TokenPair, error){
	tokenPair, err := middleware.RefreshToken(refreshToken, uc.Repository, []byte("RefreshSecretKey"))
	if err != nil {
		return model.TokenPair{}, err
	}

	return tokenPair, nil
}