package service

import (
	"srbolab_cpc/db"
	"srbolab_cpc/model"
)

var (
	SeminarTypeService seminarTypeServiceInterface = &seminarTypeService{}
)

type seminarTypeService struct {
}

type seminarTypeServiceInterface interface {
	GetAllSeminarTypes() ([]model.SeminarType, error)
}

func (c *seminarTypeService) GetAllSeminarTypes() ([]model.SeminarType, error) {
	var seminarTypes []model.SeminarType
	if err := db.Client.Find(&seminarTypes).Error; err != nil {
		return nil, err
	}
	return seminarTypes, nil
}
