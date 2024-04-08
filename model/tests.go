package model

import (
	"time"

	"gorm.io/gorm"
)

type Test struct {
	gorm.Model
	Name              string        `json:"name"`
	SeminarTheme      *SeminarTheme `json:"seminar_theme"`
	SeminarThemeID    *uint         `json:"seminar_theme_id"`
	Questions         []Question    `json:"questions" gorm:"many2many:test_question;"`
	IncludeMultiTheme *bool         `json:"include_multi_theme"`
	Practice          *bool         `json:"practice"`
	PracticeTime      time.Time     `json:"practice_time"`
}

type ClientTest struct {
	gorm.Model
	Jmbg             string           `json:"jmbg"`
	Client           Client           `json:"client"`
	ClientID         uint             `json:"client_id"`
	Test             Test             `json:"test"`
	TestID           uint             `json:"test_id"`
	SeminarDay       SeminarDay       `json:"seminar_day"`
	SeminarDayID     uint             `json:"seminar_day_id"`
	Result           float64          `json:"result"`
	ResultStr        string           `json:"result_str"`
	QuestionAnswer   []QuestionAnswer `json:"questions_answers" gorm:"-:all"`
	IsSecondFinished bool             `json:"isSecondFinished" gorm:"-:all"`
}

type QuestionAnswer struct {
	QuestionID uint   `json:"question_id"`
	Answer     string `json:"answer"`
}

type Question struct {
	gorm.Model
	Content        string        `json:"content"`
	SeminarTheme   *SeminarTheme `json:"seminar_theme"`
	SeminarThemeID *uint         `json:"seminar_theme_id"`
	Answers        []Answer      `json:"answers"`
	Image          *string       `gorm:"type:bytea" json:"image"`
	MultiTheme     *bool         `json:"multi_theme"`
}

type Answer struct {
	gorm.Model
	Content    string   `json:"content"`
	Question   Question `json:"question"`
	QuestionID uint     `json:"question_id"`
	Correct    *bool    `json:"correct"`
	Letter     string   `json:"letter"`
}
