package model

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	Person       Person  `json:"person" gorm:"embedded"`
	Address      Address `json:"address" gorm:"embedded"`
	JMBG         string  `json:"jmbg"`
	DriveLicence string  `json:"drive_licence"`
	PlaceBirth   string  `json:"place_birth"`
	CountryBirth string  `json:"country_birth"`
	CompanyID    *uint   `json:"company-id"`
	Company      Company `json:"company"`
	Verified     bool    `json:"verified"`
	Documents    []File  `json:"documents" gorm:"many2many:client_file;"`
}

type ClientPresence struct {
	ClientID     uint
	Presence     bool
	SeminarDayID uint
}
