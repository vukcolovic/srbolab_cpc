package model

import "gorm.io/gorm"

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
	Roles    []Role `gorm:"many2many:user_role;"`
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
	DriveLicence string
	PlaceBirth   string
	CountryBirth string
	CompanyID    uint
	Company      Company
	Verified     bool
	Documents    []File `gorm:"many2many:client_file;"`
}

type ClientPresence struct {
	ClientID     uint
	Presence     bool
	SeminarDayID uint
}
