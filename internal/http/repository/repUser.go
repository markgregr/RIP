package repository

import (
	"errors"

	"github.com/markgregr/RIP/internal/model"
)

type UserRepository interface {
	CreateUser(user model.User) error
	GetByEmail(email string) (model.User, error)
}

func (r *Repository) GetUsers() ([]model.User, error) {
	var users []model.User

	err := r.db.Table("users").
	Scan(&users).Error; 
	if err != nil {
        return nil, errors.New("ошибка нахождения списка пользователей")
    }
	
	return users, nil
}

func (r *Repository) GetUserByID(userID uint) (model.User, error) {
	var user model.User

	err := r.db.Table("users").Where(`"user_id" = ?`, userID).Find(&user).Error
	if err != nil {
		return model.User{}, errors.New("пользователь с данным ID не найден")
	}

	return user, nil
}

func (r *Repository) CreateUser(user model.User) error {
	err := r.db.Table("users").Create(&user).Error
	if err != nil {
		return errors.New("неудалось создать нового пользователя")
	}
	return nil
}

func (r *Repository) GetByEmail(email string) (model.User, error) {
	var user model.User

	err := r.db.Table("users").Where(`"email" = ?`, email).Find(&user).Error
	if err != nil {
		return model.User{}, errors.New("не существует пользователя с данной почтой")
	}

	return user, nil
}

func (r *Repository) GetUserRoleByID(userID uint) (model.Role, error) {
	var role model.Role

	err := r.db.Table("users").Where(`"user_id" = ?`, userID).Select("role").Scan(&role).Error
	if err != nil {
		return "", errors.New("пользователь с таким ID не найден")
	}

	return role, nil
}


