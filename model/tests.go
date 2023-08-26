package model

import (
	"gorm.io/gorm"
)

type Test struct {
	gorm.Model
	Name           string        `json:"name"`
	SeminarTheme   *SeminarTheme `json:"seminar_theme"`
	SeminarThemeID *uint         `json:"seminar_theme_id"`
	Questions      []Question    `json:"questions" gorm:"many2many:test_question;"`
}

type ClientTest struct {
	gorm.Model
	Name      string  `json:"name"`
	Client    Client  `json:"seminar_theme"`
	ClientID  uint    `json:"seminar_theme_id"`
	Test      Test    `json:"test"`
	TestID    uint    `json:"test_id"`
	Result    float64 `json:"result"`
	ResultStr string  `json:"result_str"`
}

type Question struct {
	gorm.Model
	Content        string        `json:"content"`
	SeminarTheme   *SeminarTheme `json:"seminar_theme"`
	SeminarThemeID *uint         `json:"seminar_theme_id"`
	Answers        []Answer      `json:"answers"`
}

type Answer struct {
	gorm.Model
	Content    string   `json:"content"`
	Question   Question `json:"question"`
	QuestionID uint     `json:"question_id"`
	Correct    *bool    `json:"correct"`
	Letter     string   `json:"letter"`
}
