package service

import (
	"errors"
	"gorm.io/gorm"
	"srbolab_cpc/db"
	"srbolab_cpc/model"
	"srbolab_cpc/util"
	"time"
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
	if err := db.Client.Preload("Seminar").Preload("Documents").Preload("Seminar.SeminarTheme").Preload("Seminar.SeminarTheme.BaseSeminarType").Preload("Seminar.ClassRoom.Location").Preload("Presence").Preload("Presence.Client").Preload("Classes").Preload("Classes.Teacher").Find(&seminarDay, seminarDayID).Error; err != nil {
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
	oldSeminarDay, err := c.GetSeminarDayByID(int(seminarDay.ID))
	if err != nil {
		return nil, err
	}

	result := db.Client.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&seminarDay)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, od := range oldSeminarDay.Documents {
		found := false
		for _, nd := range seminarDay.Documents {
			if od.ID == nd.ID {
				found = true
				break
			}
		}

		if !found {
			result := db.Client.Exec("DELETE FROM seminarday_file WHERE seminar_day_id = ? AND file_id = ?", seminarDay.ID, od.ID)
			if result.Error != nil {
				return nil, result.Error
			}
		}
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
	dateForDay = time.Date(dateForDay.Year(), dateForDay.Month(), dateForDay.Day(), 8, 0, 0, 0, dateForDay.Location())
	for i := 1; i <= seminar.SeminarTheme.NumberOfDays; i++ {
		presences := []model.ClientPresence{}
		for _, client := range seminar.Trainees {
			presence := false
			presences = append(presences, model.ClientPresence{ClientID: client.ClientID, Presence: &presence})
		}

		classes := []model.SeminarClass{}
		for j := 1; j <= 7; j++ {
			classes = append(classes, model.SeminarClass{Number: j})
		}
		day := model.SeminarDay{SeminarID: seminar.ID, Date: dateForDay, Number: i, Presence: presences, Classes: classes}
		result := db.Client.Create(&day)
		if result.Error != nil {
			return []model.SeminarDay{}, result.Error
		}

		seminarDays = append(seminarDays)
		dateForDay = util.IfWeekendGetFirstWorkDay(dateForDay.AddDate(0, 0, 1))
	}

	return seminarDays, nil
}
