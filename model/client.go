package model

import (
	"time"

	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	Person             Person     `json:"person" gorm:"embedded"`
	Address            Address    `json:"address" gorm:"embedded"`
	JMBG               *string    `json:"jmbg" gorm:"index:idx_jmbg,unique"`
	DriveLicence       *string    `json:"drive_licence"`
	PlaceBirth         *string    `json:"place_birth"`
	CountryBirth       *string    `json:"country_birth"`
	CompanyID          *uint      `json:"company_id"`
	Company            Company    `json:"company"`
	PartnerID          *uint      `json:"partner_id"`
	Partner            Partner    `json:"partner"`
	CompanyPIB         *string    `json:"company_pib"`
	Verified           *bool      `json:"verified"`
	WaitSeminar        *bool      `json:"wait_seminar"`
	Documents          []*File    `json:"documents" gorm:"many2many:client_file;"`
	Resident           *bool      `json:"resident"`
	SecondCitizenship  *string    `json:"second_citizenship"`
	IDCardNumber       *string    `json:"id_card_number"`
	CPCNumber          *string    `json:"cpc_number"`
	CPCDate            *time.Time `json:"cpc_date"`
	EducationalProfile *string    `json:"educational_profile"`
	Comment            *string    `json:"comment"`
	//for now changed purpose, now is number of seminar we currently do
	InitialCompletedSeminars *int             `json:"initial_completed_seminars"`
	Seminars                 []ClientSeminar  `json:"seminars"`
	CreatedBy                User             `json:"created_by"`
	CreatedByID              *uint            `json:"created_by_id"`
	VerifiedBy               User             `json:"verified_by"`
	VerifiedByID             *uint            `json:"verified_by_id"`
	CLicence                 *bool            `json:"c_licence"`
	DLicence                 *bool            `json:"d_licence"`
	PassedCheckboxes         PassedCheckboxes `json:"passed_checkboxes" gorm:"embedded"`
}

type PassedCheckboxes struct {
	WorkTimeAndTahografs *bool `json:"work_time_and_tahografs"` //1
	Tahografs2           *bool `json:"tahografs_2"`             //5
	Regulations          *bool `json:"regulations"`             //4
	Burden               *bool `json:"burden"`                  //3
	ThemeDocuments       *bool `json:"theme_documents"`         //2
	EmergencySituations  *bool `json:"emergency_situations"`    //6
}

func (a Address) GetStreetWithNumber() string {
	return a.Street + " " + a.HouseNumber
}

type ClientFilter struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	JMBG        string `json:"jmbg"`
	WaitSeminar string `json:"wait_seminar"`
	Verified    string `json:"verified"`
	CompanyID   int    `json:"company_id"`
	WaitingRoom bool   `json:"waiting_room"`
}

type ClientSeminar struct {
	gorm.Model
	ClientID           uint       `json:"client_id"`
	Client             Client     `json:"client"`
	SeminarID          uint       `json:"seminar_id"`
	Seminar            Seminar    `json:"seminar"`
	Payed              *bool      `json:"payed"`
	Pass               *bool      `json:"pass"`
	PayedBy            *string    `json:"payed_by"`
	PayDate            *time.Time `json:"pay_date"`
	ConfirmationNumber int        `json:"confirmation_number"`
}

type ClientPresence struct {
	gorm.Model
	ClientID     uint   `json:"client_id"`
	Client       Client `json:"client"`
	Presence     *bool  `json:"presence"`
	SeminarDayID uint   `json:"seminar_day_id"`
}
