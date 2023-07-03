package model

import "gorm.io/gorm"

type Person struct {
	FirstName   string `json:"first_name"`
	MiddleName  string `json:"middle_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type User struct {
	gorm.Model
	Person   Person `json:"person" gorm:"embedded"`
	Password string `json:"password"`
	Roles    []Role `json:"roles" gorm:"many2many:user_role;"`
}

func (u *User) AfterFind(tx *gorm.DB) (err error) {
	u.Password = ""
	return nil
}

type Role struct {
	gorm.Model
	Code string `json:"code"`
}
