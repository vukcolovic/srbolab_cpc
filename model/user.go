package model

import "gorm.io/gorm"

type Address struct {
	Place       string
	Street      string
	HouseNumber string
	PostCode    string
}

type Person struct {
	FirstName   string
	MiddleName  string
	LastName    string
	Email       string
	Address     string
	PhoneNumber string
}

type User struct {
	gorm.Model
	Person   Person `gorm:"embedded"`
	Password string
	Roles    []Role `gorm:"many2many:user_roles;"`
}

type Role struct {
	gorm.Model
	Code string
}

type Client struct {
	gorm.Model
	Person       Person  `gorm:"embedded"`
	Address      Address `gorm:"embedded"`
	JMBG         string
	PlaceBirth   string
	CountryBirth string
}
