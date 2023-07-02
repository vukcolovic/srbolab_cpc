package model

import "gorm.io/gorm"

type Address struct {
	Place       string `json:"place"`
	Street      string `json:"street"`
	HouseNumber string `json:"house_number"`
	PostCode    string `json:"post_code"`
}

type Location struct {
	gorm.Model
	Address Address `gorm:"embedded"`
}
