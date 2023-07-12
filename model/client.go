package model

import (
	"gorm.io/gorm"
	"time"
)

type Client struct {
	gorm.Model
	Person                   Person          `json:"person" gorm:"embedded"`
	Address                  Address         `json:"address" gorm:"embedded"`
	JMBG                     *string         `json:"jmbg"`
	DriveLicence             *string         `json:"drive_licence"`
	PlaceBirth               *string         `json:"place_birth"`
	CountryBirth             *string         `json:"country_birth"`
	CompanyID                *uint           `json:"company_id"`
	Company                  Company         `json:"company"`
	Verified                 *bool           `json:"verified"`
	WaitSeminar              *bool           `json:"wait_seminar"`
	Documents                []*File         `json:"documents" gorm:"many2many:client_file;"`
	Resident                 *bool           `json:"resident"`
	SecondCitizenship        *string         `json:"second_citizenship"`
	IDCardNumber             *string         `json:"id_card_number"`
	CPCNumber                *string         `json:"cpc_number"`
	CPCDate                  time.Time       `json:"cpc_date"`
	EducationalProfile       *string         `json:"educational_profile"`
	Comment                  *string         `json:"comment"`
	InitialCompletedSeminars *int            `json:"initial_completed_seminars"`
	Seminars                 []ClientSeminar `json:"seminars"`
}

type ClientFilter struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	JMBG        string `json:"jmbg"`
	WaitSeminar bool   `json:"wait_seminar"`
	Verified    bool   `json:"verified"`
}

type ClientSeminar struct {
	gorm.Model
	ClientID  uint    `json:"client_id"`
	Client    Client  `json:"client"`
	SeminarID uint    `json:"seminar_id"`
	Seminar   Seminar `json:"seminar"`
	Payed     *bool   `json:"payed"`
	Pass      *bool   `json:"pass"`
}

type ClientPresence struct {
	ClientID     uint
	Presence     bool
	SeminarDayID uint
}
