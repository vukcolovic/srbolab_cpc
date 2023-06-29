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
	DeleteUser(id int) error
	CreateUser(user model.User) (*model.User, error)
	UpdateUser(user model.User) (*model.User, error)
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

func (s *userService) DeleteUser(id int) error {
	return db.Client.Delete(&model.User{}, id).Error
}

func (s *userService) CreateUser(user model.User) (*model.User, error) {
	result := db.Client.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (s *userService) UpdateUser(user model.User) (*model.User, error) {
	result := db.Client.Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
