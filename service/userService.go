package service

import (
	"gorm.io/gorm"
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
	if err := db.Client.Preload("Roles").First(&user, id).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByEmail(email string) (*model.User, error) {
	var user *model.User
	if err := db.Client.Where("email = ?", email).Preload("Roles").First(&user).Error; err != nil {
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
	oldUser, err := s.GetUserByID(int(user.ID))
	if err != nil {
		return nil, err
	}

	result := db.Client.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, or := range oldUser.Roles {
		found := false
		for _, nr := range user.Roles {
			if or.ID == nr.ID {
				found = true
				break
			}
		}

		if !found {
			result := db.Client.Exec("DELETE FROM user_role WHERE user_id = ? AND role_id = ?", user.ID, or.ID)
			if result.Error != nil {
				return nil, result.Error
			}
		}
	}

	return &user, nil
}
