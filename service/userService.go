package service

import (
	"srbolab_cpc/db"
	"srbolab_cpc/logoped"
	"srbolab_cpc/model"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	UsersService usersServiceInterface = &userService{}
)

type userService struct {
}

type usersServiceInterface interface {
	GetUserIDByToken(token string) (int, error)
	GetAllUsers(skip, take int) ([]model.User, error)
	GetUserByID(id int) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUsersCount() (int64, error)
	DeleteUser(id int) error
	CreateUser(user model.User) (*model.User, error)
	UpdateUser(user model.User) (*model.User, error)
	GetAllTeachers() ([]model.User, error)
}

func (s *userService) GetUserIDByToken(token string) (int, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		logoped.ErrorLog.Println("Error getting user by token, error parse claims: ", err)
		return 0, err
	}

	id, err := strconv.Atoi(claims["Id"].(string))
	if err != nil {
		logoped.ErrorLog.Println("Error getting user by token: ", err)
		return 0, err
	}

	return id, nil
}

func (s *userService) GetAllUsers(skip, take int) ([]model.User, error) {
	var users []model.User
	if err := db.Client.Limit(take).Offset(skip).Find(&users).Error; err != nil {
		return nil, err
	}

	for _, u := range users {
		u.Password = ""
	}

	return users, nil
}

func (s *userService) GetAllTeachers() ([]model.User, error) {
	var users []model.User
	if err := db.Client.Where("is_teacher = ?", true).Find(&users).Error; err != nil {
		return nil, err
	}

	for _, u := range users {
		u.Password = ""
	}

	return users, nil
}

func (s *userService) GetUserByID(id int) (*model.User, error) {
	var user *model.User
	if err := db.Client.Preload("Roles").First(&user, id).Error; err != nil {
		return nil, err
	}

	user.Password = ""

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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		logoped.ErrorLog.Println("Error creating user, hashing password error: ", err)
		return nil, err
	}

	user.Password = string(hashedPassword)
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

	user.Password = oldUser.Password

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

	user.Password = ""
	return &user, nil
}
