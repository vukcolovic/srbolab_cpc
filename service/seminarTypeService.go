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
	GetAllBaseSeminarTypes() ([]model.BaseSeminarType, error)
	GetSeminarThemesBypeID(typeID int) ([]model.SeminarTheme, error)
}

func (c *seminarTypeService) GetAllBaseSeminarTypes() ([]model.BaseSeminarType, error) {
	var seminarBaseTypes []model.BaseSeminarType
	if err := db.Client.Find(&seminarBaseTypes).Error; err != nil {
		return nil, err
	}
	return seminarBaseTypes, nil
}

func (c *seminarTypeService) GetSeminarThemesBypeID(typeID int) ([]model.SeminarTheme, error) {
	var seminarThemes []model.SeminarTheme
	if err := db.Client.Where("base_seminar_type_id", typeID).Find(&seminarThemes).Error; err != nil {
		return nil, err
	}
	return seminarThemes, nil
}
