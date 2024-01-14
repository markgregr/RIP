package repository

import (
	"errors"

	"github.com/markgregr/RIP/internal/model"
)

type UserRepository interface {
	CreateUser(user model.User) error
	GetByEmail(email string) (model.User, error)
	GetUserRoleByID(userID uint) (model.Role, error) 
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


