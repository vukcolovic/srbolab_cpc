package service

import (
	"gorm.io/gorm"
	"srbolab_cpc/db"
	"srbolab_cpc/model"
)

var (
	QuestionService questionServiceInterface = &questionService{}
)

type questionService struct {
}

type questionServiceInterface interface {
	GetAllQuestions(skip, take int) ([]model.Question, error)
	GetAllQuestionsBySeminarTheme(seminarThemeID int) ([]model.Question, error)
	GetQuestionByID(id int) (*model.Question, error)
	GetQuestionsCount() (int64, error)
	DeleteQuestion(id int) error
	CreateQuestion(question model.Question) (*model.Question, error)
	UpdateQuestion(question model.Question) (*model.Question, error)
}

func (c *questionService) GetAllQuestions(skip, take int) ([]model.Question, error) {
	var questions []model.Question
	if err := db.Client.Order("id desc").Limit(take).Offset(skip).Preload("SeminarTheme").Preload("SeminarTheme.BaseSeminarType").Preload("Answers").Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

func (c *questionService) GetAllQuestionsBySeminarTheme(seminarThemeID int) ([]model.Question, error) {
	var questions []model.Question
	if err := db.Client.Order("id desc").Where("seminar_theme_id = ?", seminarThemeID).Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

func (c *questionService) GetQuestionByID(id int) (*model.Question, error) {
	var question *model.Question
	if err := db.Client.Preload("SeminarTheme").Preload("SeminarTheme.BaseSeminarType").Preload("Answers").First(&question, id).Error; err != nil {
		return nil, err
	}

	return question, nil
}

func (c *questionService) GetQuestionsCount() (int64, error) {
	var count int64
	if err := db.Client.Model(&model.Question{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (c *questionService) DeleteQuestion(id int) error {
	return db.Client.Delete(&model.Question{}, id).Error
}

func (c *questionService) CreateQuestion(question model.Question) (*model.Question, error) {
	result := db.Client.Create(&question)
	if result.Error != nil {
		return nil, result.Error
	}

	return &question, nil
}

func (c *questionService) UpdateQuestion(question model.Question) (*model.Question, error) {
	result := db.Client.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&question)
	if result.Error != nil {
		return nil, result.Error
	}

	return &question, nil
}
