package service

import (
	"srbolab_cpc/db"
	"srbolab_cpc/model"
)

var (
	SurveyService surveyServiceInterface = &surveyService{}
)

type surveyService struct {
}

type surveyServiceInterface interface {
	GetAllSurveys() ([]model.Survey, error)
	GetSurveyByID(id int) (*model.Survey, error)
	CreateSurvey(survey model.Survey) (*model.Survey, error)
	GetAllSurveyQuestions() ([]model.SurveyQuestion, error)
	GetSurveyQuestionByID(id int) (*model.SurveyQuestion, error)
	CreateSurveyQuestion(question model.SurveyQuestion) (*model.SurveyQuestion, error)
}

func (c *surveyService) GetAllSurveys() ([]model.Survey, error) {
	var surveys []model.Survey
	if err := db.Client.Order("id desc").Preload("Questions").Find(&surveys).Error; err != nil {
		return nil, err
	}
	return surveys, nil
}

func (c *surveyService) GetSurveyByID(id int) (*model.Survey, error) {
	var survey *model.Survey
	if err := db.Client.Preload("Questions").First(&survey, id).Error; err != nil {
		return nil, err
	}

	return survey, nil
}

func (c *surveyService) CreateSurvey(survey model.Survey) (*model.Survey, error) {
	result := db.Client.Create(&survey)
	if result.Error != nil {
		return nil, result.Error
	}

	return &survey, nil
}

func (c *surveyService) GetAllSurveyQuestions() ([]model.SurveyQuestion, error) {
	var questions []model.SurveyQuestion
	if err := db.Client.Order("id desc").Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

func (c *surveyService) GetSurveyQuestionByID(id int) (*model.SurveyQuestion, error) {
	var question *model.SurveyQuestion
	if err := db.Client.First(&question, id).Error; err != nil {
		return nil, err
	}

	return question, nil
}

func (c *surveyService) CreateSurveyQuestion(question model.SurveyQuestion) (*model.SurveyQuestion, error) {
	result := db.Client.Create(&question)
	if result.Error != nil {
		return nil, result.Error
	}

	return &question, nil
}
