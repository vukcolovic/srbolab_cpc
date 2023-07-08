package service

import (
	"srbolab_cpc/db"
	"srbolab_cpc/model"
)

var (
	SeminarService seminarServiceInterface = &seminarService{}
)

type seminarService struct {
}

type seminarServiceInterface interface {
	GetAllSeminars(skip, take int) ([]model.Seminar, error)
	GetSeminarByID(id int) (*model.Seminar, error)
	GetSeminarsCount() (int64, error)
	DeleteSeminar(id int) error
	CreateSeminar(seminar model.Seminar) (*model.Seminar, error)
	UpdateSeminar(seminar model.Seminar) (*model.Seminar, error)
}

func (c *seminarService) GetAllSeminars(skip, take int) ([]model.Seminar, error) {
	var seminars []model.Seminar
	if err := db.Client.Order("id desc").Limit(take).Offset(skip).Joins("Location").Joins("SeminarType").Joins("SeminarStatus").Find(&seminars).Error; err != nil {
		return nil, err
	}
	return seminars, nil
}

func (c *seminarService) GetSeminarByID(id int) (*model.Seminar, error) {
	var seminar *model.Seminar
	if err := db.Client.Joins("Location").Joins("SeminarType").Joins("SeminarStatus").First(&seminar, id).Error; err != nil {
		return nil, err
	}

	return seminar, nil
}

func (c *seminarService) GetSeminarsCount() (int64, error) {
	var count int64
	if err := db.Client.Model(&model.Seminar{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (c *seminarService) DeleteSeminar(id int) error {
	return db.Client.Delete(&model.Seminar{}, id).Error
}

func (c *seminarService) CreateSeminar(seminar model.Seminar) (*model.Seminar, error) {
	result := db.Client.Create(&seminar)
	if result.Error != nil {
		return nil, result.Error
	}

	return &seminar, nil
}

func (c *seminarService) UpdateSeminar(seminar model.Seminar) (*model.Seminar, error) {
	result := db.Client.Save(&seminar)
	if result.Error != nil {
		return nil, result.Error
	}

	return &seminar, nil
}
