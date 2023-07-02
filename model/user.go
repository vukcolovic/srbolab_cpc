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
	Person   Person `gorm:"embedded"`
	Password string
	Roles    []Role `gorm:"many2many:user_role;"`
}

type Role struct {
	gorm.Model
	Code string
}
