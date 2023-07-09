package model

import (
	"gorm.io/gorm"
	"time"
)

type Seminar struct {
	gorm.Model
	Start           time.Time     `json:"start_date"`
	LocationID      uint          `json:"location_id"`
	Location        Location      `json:"location"`
	Trainees        []Client      `json:"trainees" gorm:"many2many:seminar_client;"`
	SeminarTypeID   uint          `json:"seminar_type_id"`
	SeminarType     SeminarType   `json:"seminar_type"`
	SeminarStatusID uint          `json:"seminar_status_id"`
	SeminarStatus   SeminarStatus `json:"seminar_status"`
	Days            []SeminarDay  `json:"days"`
}

type SeminarType struct {
	gorm.Model
	Code string `json:"code"`
	Name string `json:"name"`
}

type SeminarStatus struct {
	gorm.Model
	Code string `json:"code"`
	Name string `json:"name"`
}

type SeminarDay struct {
	gorm.Model
	Date      time.Time        `json:"date"`
	Number    int              `json:"number"`
	Name      string           `json:"name"`
	Classes   []SeminarClass   `json:"classes"`
	SeminarID uint             `json:"seminar_id"`
	Presence  []ClientPresence `json:"presence"`
}

type SeminarClass struct {
	gorm.Model
	Name         string
	Teachers     []User `gorm:"many2many:class_teacher;"`
	SeminarDayID uint
}
