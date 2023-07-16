package service

import (
	"errors"
	"srbolab_cpc/db"
	"srbolab_cpc/model"
	"srbolab_cpc/util"
)

var (
	SeminarDayService seminarDayServiceInterface = &seminarDayService{}
)

type seminarDayService struct {
}

type seminarDayServiceInterface interface {
	GetSeminarDayByID(seminarDayID int) (*model.SeminarDay, error)
	GetSeminarDaysBySeminarID(seminarID int) ([]model.SeminarDay, error)
	CreateSeminarDay(seminarDay model.SeminarDay) (*model.SeminarDay, error)
	UpdateSeminarDay(seminarDay model.SeminarDay) (*model.SeminarDay, error)
	CreateAllSeminarDaysForSeminar(seminarID int) ([]model.SeminarDay, error)
}

func (c *seminarDayService) GetSeminarDaysBySeminarID(seminarID int) ([]model.SeminarDay, error) {
	var days []model.SeminarDay
	if err := db.Client.Where("seminar_id", seminarID).Find(&days).Error; err != nil {
		return nil, err
	}
	return days, nil
}

func (c *seminarDayService) GetSeminarDayByID(seminarDayID int) (*model.SeminarDay, error) {
	var seminarDay *model.SeminarDay
	if err := db.Client.Find(&seminarDay, seminarDayID).Error; err != nil {
		return nil, err
	}

	return seminarDay, nil
}

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

func (c *seminarDayService) CreateAllSeminarDaysForSeminar(seminarID int) ([]model.SeminarDay, error) {
	seminar, err := SeminarService.GetSeminarByID(seminarID)
	if err != nil {
		return []model.SeminarDay{}, err
	}

	if len(seminar.Days) > 0 {
		return []model.SeminarDay{}, errors.New("Greška prilikom pravljenja dana za seminar, seminar već ima dane!")
	}

	seminarDays := []model.SeminarDay{}
	dateForDay := util.IfWeekendGetFirstWorkDay(seminar.Start)
	for i := 1; i <= seminar.SeminarTheme.NumberOfDays; i++ {
		day := model.SeminarDay{SeminarID: seminar.ID, Date: dateForDay, Number: i}
		dateForDay = util.IfWeekendGetFirstWorkDay(dateForDay.AddDate(0, 0, 1))
		result := db.Client.Create(&day)
		if result.Error != nil {
			return []model.SeminarDay{}, result.Error
		}
		if err != nil {
			return []model.SeminarDay{day}, err
		}
		seminarDays = append(seminarDays)
	}

	return seminarDays, nil
}
