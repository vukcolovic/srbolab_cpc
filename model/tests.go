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
