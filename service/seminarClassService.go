package service

import (
	"gorm.io/gorm"
	"srbolab_cpc/db"
	"srbolab_cpc/model"
)

var (
	SeminarClassService seminarClassServiceInterface = &seminarClassService{}
)

type seminarClassService struct {
}

type seminarClassServiceInterface interface {
	GetSeminarClassNameByID(seminarClassNameID int) (*model.SeminarClassName, error)
	GetAllSeminarClassNames() ([]model.SeminarClassName, error)
	CreateSeminarClassName(seminarClassName model.SeminarClassName) (*model.SeminarClassName, error)
	UpdateSeminarClassName(seminarClassName model.SeminarClassName) (*model.SeminarClassName, error)
	GetSeminarClassesNamesBySeminarThemeAndDayNumberAsMap(seminarThemeID uint, seminarDayNumber int) (map[int]string, error)
}

func (c *seminarClassService) GetSeminarClassNameByID(seminarClassNameID int) (*model.SeminarClassName, error) {
	var classNames *model.SeminarClassName
	if err := db.Client.Preload("SeminarTheme.BaseSeminarType").Find(&classNames, seminarClassNameID).Error; err != nil {
		return nil, err
	}
	return classNames, nil
}

func (c *seminarClassService) GetAllSeminarClassNames() ([]model.SeminarClassName, error) {
	var classNames []model.SeminarClassName
	if err := db.Client.Preload("SeminarTheme.BaseSeminarType").Find(&classNames).Error; err != nil {
		return nil, err
	}

	return classNames, nil
}

func (c *seminarClassService) CreateSeminarClassName(seminarClassName model.SeminarClassName) (*model.SeminarClassName, error) {
	result := db.Client.Create(&seminarClassName)
	if result.Error != nil {
		return nil, result.Error
	}

	return &seminarClassName, nil
}

func (c *seminarClassService) UpdateSeminarClassName(seminarClassName model.SeminarClassName) (*model.SeminarClassName, error) {
	result := db.Client.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&seminarClassName)
	if result.Error != nil {
		return nil, result.Error
	}

	return &seminarClassName, nil
}

func (c *seminarClassService) GetSeminarClassesNamesBySeminarThemeAndDayNumberAsMap(seminarThemeID uint, seminarDayNumber int) (map[int]string, error) {
	var classes []model.SeminarClassName
	if err := db.Client.Where("seminar_theme_id = ? AND day_number = ?", seminarThemeID, seminarDayNumber).Find(&classes).Error; err != nil {
		return nil, err
	}

	mapClasses := map[int]string{}
	for _, c := range classes {
		mapClasses[c.ClassNumber] = c.ClassName
	}

	return mapClasses, nil
}
