package model

import "gorm.io/gorm"

type Address struct {
	Place       string
	Street      string
	HouseNumber string
	PostCode    string
}

type Location struct {
	gorm.Model
	Address Address `gorm:"embedded"`
}
