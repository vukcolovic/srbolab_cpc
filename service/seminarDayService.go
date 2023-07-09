package service

import (
	"srbolab_cpc/db"
	"srbolab_cpc/model"
)

var (
	SeminarDayService seminarDayServiceInterface = &seminarDayService{}
)

type seminarDayService struct {
}

type seminarDayServiceInterface interface {
	GetSeminarDaysBySeminarID(seminarID int) ([]model.SeminarDay, error)
	CreateSeminarDay(seminarDay model.SeminarDay) (*model.SeminarDay, error)
	UpdateSeminarDay(seminarDay model.SeminarDay) (*model.SeminarDay, error)
}

func (c *seminarDayService) GetSeminarDaysBySeminarID(seminarID int) ([]model.SeminarDay, error) {
	var days []model.SeminarDay
	if err := db.Client.Where("seminar_id", seminarID).Find(&days).Error; err != nil {
		return nil, err
	}
	return days, nil
}

//
//func (c *seminarService) GetSeminarByID(id int) (*model.Seminar, error) {
//	var seminar *model.Seminar
//	if err := db.Client.Joins("Location").Joins("SeminarType").Joins("SeminarStatus").First(&seminar, id).Error; err != nil {
//		return nil, err
//	}
//
//	return seminar, nil
//}

func (c *seminarDayService) CreateSeminarDay(seminarDay model.SeminarDay) (*model.SeminarDay, error) {
	result := db.Client.Create(&seminarDay)
	if result.Error != nil {
		return nil, result.Error
	}

	return &seminarDay, nil
}

func (c *seminarDayService) UpdateSeminarDay(seminarDay model.SeminarDay) (*model.SeminarDay, error) {
	result := db.Client.Save(&seminarDay)
	if result.Error != nil {
		return nil, result.Error
	}

	return &seminarDay, nil
}
