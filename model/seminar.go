package model

import (
	"gorm.io/gorm"
	"time"
)

type Seminar struct {
	gorm.Model
	Start           time.Time
	LocationID      uint
	Location        Location
	Trainees        []Client `gorm:"many2many:seminar_client;"`
	SeminarTypeID   uint
	SeminarType     SeminarType
	SeminarStatusID uint
	SeminarStatus   SeminarStatus
	Days            []SeminarDay
}

type SeminarType struct {
	gorm.Model
	Code string
	Name string
}

type SeminarStatus struct {
	gorm.Model
	Code string
	Name string
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
