package model

import "gorm.io/gorm"

const (
	GENERAL = iota
	TEACHER
)

type Survey struct {
	gorm.Model
	Name      string           `json:"name"`
	Questions []SurveyQuestion `json:"questions" gorm:"many2many:survey_question;"`
	Active    *bool            `json:"active"`
	Type      int              `json:"type" gorm:"default:0"`
}

type ClientSurvey struct {
	gorm.Model
	JMBG                  string                 `json:"jmbg" gorm:"-:all"`
	Client                Client                 `json:"client"`
	ClientID              uint                   `json:"client_id"`
	Survey                Survey                 `json:"survey"`
	SurveyID              uint                   `json:"survey_id"`
	SeminarDay            SeminarDay             `json:"seminar_day"`
	SeminarDayID          uint                   `json:"seminar_day_id"`
	TeacherID             *int                   `json:"teacher_id" gorm:"default:null"`
	Teacher               *User                  `json:"teacher"`
	SurveyQuestionAnswers []SurveyQuestionAnswer `json:"survey_questions_answers"`
}

type SurveyQuestion struct {
	gorm.Model
	Content string `json:"content"`
}

type SurveyQuestionAnswer struct {
	gorm.Model
	SurveyQuestion   SurveyQuestion `json:"survey_question"`
	SurveyQuestionID uint           `json:"survey_question_id"`
	ClientSurvey     ClientSurvey   `json:"client_survey"`
	ClientSurveyID   uint           `json:"client_survey_id"`
	Grade            int            `json:"grade"`
}
