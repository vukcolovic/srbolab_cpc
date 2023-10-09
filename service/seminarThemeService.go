package service

import (
	"gorm.io/gorm"
	"srbolab_cpc/db"
	"srbolab_cpc/model"
)

var (
	SeminarThemeService seminarThemeServiceInterface = &seminarThemeService{}
)

type seminarThemeService struct {
}

type seminarThemeServiceInterface interface {
	GetSeminarDayThemeNameByID(seminarDayThemeNameID int) (*model.SeminarDayThemeName, error)
	GetAllSeminarDayThemeNames() ([]model.SeminarDayThemeName, error)
	CreateSeminarDayThemeName(seminarDayThemeName model.SeminarDayThemeName) (*model.SeminarDayThemeName, error)
	UpdateSeminarDayThemeName(seminarDayThemeName model.SeminarDayThemeName) (*model.SeminarDayThemeName, error)
	GetSeminarThemeNamesBySeminarThemeAsMap(seminarThemeID uint) (map[int]string, error)
}

func (c *seminarThemeService) GetSeminarDayThemeNameByID(seminarDayThemeNameID int) (*model.SeminarDayThemeName, error) {
	var themeName *model.SeminarDayThemeName
	if err := db.Client.Preload("SeminarTheme.BaseSeminarType").Find(&themeName, seminarDayThemeNameID).Error; err != nil {
		return nil, err
	}
	return themeName, nil
}

func (c *seminarThemeService) GetAllSeminarDayThemeNames() ([]model.SeminarDayThemeName, error) {
	var themeNames []model.SeminarDayThemeName
	if err := db.Client.Preload("SeminarTheme.BaseSeminarType").Find(&themeNames).Error; err != nil {
		return nil, err
	}

	return themeNames, nil
}

func (c *seminarThemeService) CreateSeminarDayThemeName(seminarDayThemeName model.SeminarDayThemeName) (*model.SeminarDayThemeName, error) {
	result := db.Client.Create(&seminarDayThemeName)
	if result.Error != nil {
		return nil, result.Error
	}

	return &seminarDayThemeName, nil
}

func (c *seminarThemeService) UpdateSeminarDayThemeName(seminarDayThemeName model.SeminarDayThemeName) (*model.SeminarDayThemeName, error) {
	result := db.Client.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&seminarDayThemeName)
	if result.Error != nil {
		return nil, result.Error
	}

	return &seminarDayThemeName, nil
}

func (c *seminarThemeService) GetSeminarThemeNamesBySeminarThemeAsMap(seminarThemeID uint) (map[int]string, error) {
	var themes []model.SeminarDayThemeName
	if err := db.Client.Where("seminar_theme_id", seminarThemeID).Find(&themes).Error; err != nil {
		return nil, err
	}

	mapThemes := map[int]string{}
	for _, t := range themes {
		mapThemes[t.DayNumber] = t.ThemeName
	}

	return mapThemes, nil
}
