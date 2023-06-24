package service

import (
	"srbolab_cpc/db"
	"srbolab_cpc/model"
)

var (
	UsersService usersServiceInterface = &userService{}
)

type userService struct {
}

type usersServiceInterface interface {
	GetAllUsers(skip, take int) ([]model.User, error)
	GetUserByID(id int) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUsersCount() (int64, error)
}

func (s *userService) GetAllUsers(skip, take int) ([]model.User, error) {
	var users []model.User
	if err := db.Client.Limit(take).Offset(skip).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userService) GetUserByID(id int) (*model.User, error) {
	var user *model.User
	if err := db.Client.First(&user, id).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByEmail(email string) (*model.User, error) {
	var user *model.User
	if err := db.Client.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUsersCount() (int64, error) {
	var count int64
	if err := db.Client.Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
