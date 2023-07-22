package service

import (
	"srbolab_cpc/db"
	"srbolab_cpc/model"
)

var (
	SeminarStatusService seminarStatusServiceInterface = &seminarStatusService{}
)

type seminarStatusService struct {
}

type seminarStatusServiceInterface interface {
	GetSeminarStatusByID(statusID int) (*model.SeminarStatus, error)
	GetSeminarStatusByCode(code string) (*model.SeminarStatus, error)
	GetAllSeminarStatuses() ([]model.SeminarStatus, error)
}

func (c *seminarStatusService) GetSeminarStatusByID(statusID int) (*model.SeminarStatus, error) {
	var status *model.SeminarStatus
	if err := db.Client.First(&status, statusID).Error; err != nil {
		return nil, err
	}

	return status, nil
}

func (c *seminarStatusService) GetSeminarStatusByCode(code string) (*model.SeminarStatus, error) {
	var status *model.SeminarStatus
	if err := db.Client.Where("code", code).First(&status).Error; err != nil {
		return nil, err
	}

	return status, nil
}

func (c *seminarStatusService) GetAllSeminarStatuses() ([]model.SeminarStatus, error) {
	var seminarStatuses []model.SeminarStatus
	if err := db.Client.Find(&seminarStatuses).Error; err != nil {
		return nil, err
	}
	return seminarStatuses, nil
}
