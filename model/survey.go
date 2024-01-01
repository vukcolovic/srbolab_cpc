package model

import "gorm.io/gorm"

type Survey struct {
	gorm.Model
	Name      string           `json:"name"`
	Questions []SurveyQuestion `json:"questions" gorm:"many2many:survey_question;"`
	Active    *bool            `json:"active"`
}

type ClientSurvey struct {
	gorm.Model
	Client                Client                 `json:"client"`
	ClientID              uint                   `json:"client_id"`
	Survey                Survey                 `json:"survey"`
	SurveyID              uint                   `json:"survey_id"`
	SeminarDay            SeminarDay             `json:"seminar_day"`
	SeminarDayID          uint                   `json:"seminar_day_id"`
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
