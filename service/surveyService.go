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
	GetActiveSurveyByType(surveyType int) (*model.Survey, error)
	GetActiveSurveys() ([]model.Survey, error)
	CreateClientSurvey(cs model.ClientSurvey) (*model.ClientSurvey, error)
	GetClientSurveysBySeminarDayIDAndType(seminarDayID int, surveyType int) ([]model.ClientSurvey, error)
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

func (c *surveyService) GetActiveSurveyByType(surveyType int) (*model.Survey, error) {
	var survey *model.Survey
	if err := db.Client.Where("active = ? AND type = ?", true, surveyType).Preload("Questions").Last(&survey).Error; err != nil {
		return nil, err
	}

	return survey, nil
}

func (c *surveyService) GetActiveSurveys() ([]model.Survey, error) {
	result := []model.Survey{}
	generalSurvey, err := c.GetActiveSurveyByType(model.GENERAL)
	if err != nil {
		return []model.Survey{}, err
	}
	result = append(result, *generalSurvey)

	teacherSurvey, err := c.GetActiveSurveyByType(model.TEACHER)
	if err != nil {
		return []model.Survey{}, err
	}
	result = append(result, *teacherSurvey)

	return result, nil
}

func (c *surveyService) CreateClientSurvey(cs model.ClientSurvey) (*model.ClientSurvey, error) {
	result := db.Client.Create(&cs)
	if result.Error != nil {
		return nil, result.Error
	}

	return &cs, nil
}

func (c *surveyService) GetClientSurveysBySeminarDayIDAndType(seminarDayID int, surveyType int) ([]model.ClientSurvey, error) {
	var clientSurveys []model.ClientSurvey
	if err := db.Client.Where("seminar_day_id = ?", seminarDayID).Preload("SurveyQuestionAnswers.SurveyQuestion").Preload("Client").Preload("Survey.Questions").Preload("Teacher").Find(&clientSurveys).Error; err != nil {
		return nil, err
	}

	clientSurveysByType := []model.ClientSurvey{}

	for _, cs := range clientSurveys {
		if cs.Survey.Type == surveyType {
			clientSurveysByType = append(clientSurveysByType, cs)
		}
	}

	return clientSurveysByType, nil
}
