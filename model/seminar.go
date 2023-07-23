package model

import (
	"gorm.io/gorm"
	"time"
)

const (
	SEMINAR_STATUS_OPENED      = 1
	SEMINAR_STATUS_FILLED      = 2
	SEMINAR_STATUS_IN_PROGRESS = 3
	SEMINAR_STATUS_CLOSED      = 4
)

type Seminar struct {
	gorm.Model
	Start           time.Time       `json:"start_date"`
	ClassRoomID     uint            `json:"class_room_id"`
	ClassRoom       ClassRoom       `json:"class_room"`
	Trainees        []ClientSeminar `json:"trainees"`
	SeminarThemeID  uint            `json:"seminar_theme_id"`
	SeminarTheme    SeminarTheme    `json:"seminar_theme"`
	SeminarStatusID uint            `json:"seminar_status_id"`
	SeminarStatus   SeminarStatus   `json:"seminar_status"`
	Days            []SeminarDay    `json:"days"`
}

type BaseSeminarType struct {
	gorm.Model
	Code string `json:"code"`
	Name string `json:"name"`
}

type SeminarTheme struct {
	gorm.Model
	BaseSeminarTypeID uint            `json:"base_seminar_type_id"`
	BaseSeminarType   BaseSeminarType `json:"base_seminar_type"`
	Code              string          `json:"code"`
	Name              string          `json:"name"`
	NumberOfDays      int             `json:"number_of_days"`
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
	Seminar   Seminar          `json:"seminar"`
	Presence  []ClientPresence `json:"presence"`
}

type SeminarClass struct {
	gorm.Model
	Name         string
	Teachers     []User `gorm:"many2many:class_teacher;"`
	SeminarDayID uint
}
