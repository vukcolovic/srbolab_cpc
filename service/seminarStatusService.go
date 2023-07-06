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
	GetAllSeminarStatuses() ([]model.SeminarStatus, error)
}

func (c *seminarStatusService) GetAllSeminarStatuses() ([]model.SeminarStatus, error) {
	var seminarStatuses []model.SeminarStatus
	if err := db.Client.Find(&seminarStatuses).Error; err != nil {
		return nil, err
	}
	return seminarStatuses, nil
}
