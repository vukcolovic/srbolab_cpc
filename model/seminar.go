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
	Trainees        []Client      `gorm:"many2many:seminar_client;"`
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
	Date      time.Time
	Name      string
	Classes   []SeminarClass
	SeminarID uint
	Presence  []ClientPresence
}

type SeminarClass struct {
	gorm.Model
	Name         string
	Teachers     []User `gorm:"many2many:class_teacher;"`
	SeminarDayID uint
}
