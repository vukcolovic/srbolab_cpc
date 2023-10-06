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
	Address Address `json:"address" gorm:"embedded"`
	Code    string  `json:"code"`
}

type ClassRoom struct {
	gorm.Model
	Name        string   `json:"name"`
	MaxStudents int      `json:"max_students"`
	LocationID  uint     `json:"location_id"`
	Location    Location `json:"location"`
}
