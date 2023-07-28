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
	GetAllSeminarsByStatus(statusCode string) ([]model.Seminar, error)
	GetSeminarByID(id int) (*model.Seminar, error)
	GetSeminarsCount() (int64, error)
	DeleteSeminar(id int) error
	CreateSeminar(seminar model.Seminar) (*model.Seminar, error)
	UpdateSeminar(seminar model.Seminar) (*model.Seminar, error)
	DeleteSeminarClient(clientSeminar model.ClientSeminar) error
	UpdateSeminarStatusIfNeed(seminarID int) error
}

func (c *seminarService) GetAllSeminars(skip, take int) ([]model.Seminar, error) {
	var seminars []model.Seminar
	if err := db.Client.Order("id desc").Limit(take).Offset(skip).Preload("ClassRoom.Location").Joins("ClassRoom").Joins("SeminarTheme").Joins("SeminarTheme.BaseSeminarType").Joins("SeminarStatus").Find(&seminars).Error; err != nil {
		return nil, err
	}
	return seminars, nil
}

func (c *seminarService) GetAllSeminarsByStatus(statusCode string) ([]model.Seminar, error) {
	status, err := SeminarStatusService.GetSeminarStatusByCode(statusCode)
	if err != nil {
		return nil, err
	}
	var seminars []model.Seminar
	if err := db.Client.Where("seminar_status_id", status.ID).Order("id desc").Preload("ClassRoom.Location").Joins("ClassRoom").Joins("SeminarTheme").Joins("SeminarTheme.BaseSeminarType").Joins("SeminarStatus").Find(&seminars).Error; err != nil {
		return nil, err
	}
	return seminars, nil
}

func (c *seminarService) GetSeminarByID(id int) (*model.Seminar, error) {
	var seminar *model.Seminar
	if err := db.Client.Preload("Trainees").Preload("Trainees.Client").Preload("Trainees.Client.Company").Preload("Trainees.Client.Seminars").Preload("Days").Preload("Documents").Preload("ClassRoom.Location").Joins("ClassRoom").Joins("SeminarTheme.BaseSeminarType").Joins("SeminarTheme").Joins("SeminarStatus").First(&seminar, id).Error; err != nil {
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
	oldSeminar, err := c.GetSeminarByID(int(seminar.ID))
	if err != nil {
		return nil, err
	}

	result := db.Client.Save(&seminar)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, od := range oldSeminar.Documents {
		found := false
		for _, nd := range seminar.Documents {
			if od.ID == nd.ID {
				found = true
				break
			}
		}

		if !found {
			result := db.Client.Exec("DELETE FROM seminar_file WHERE seminar_id = ? AND file_id = ?", seminar.ID, od.ID)
			if result.Error != nil {
				return nil, result.Error
			}
		}
	}

	return &seminar, nil
}

func (c *seminarService) DeleteSeminarClient(clientSeminar model.ClientSeminar) error {
	result := db.Client.Exec("DELETE FROM client_seminars WHERE client_id = ? AND seminar_id = ?", clientSeminar.ClientID, clientSeminar.SeminarID)
	if result.Error != nil {
		return result.Error
	}

	return c.UpdateSeminarStatusIfNeed(int(clientSeminar.Seminar.ID))
}

func (c *seminarService) UpdateSeminarStatusIfNeed(seminarID int) error {
	seminar, err := c.GetSeminarByID(seminarID)
	if err != nil {
		return err
	}

	if seminar.SeminarStatus.ID == model.SEMINAR_STATUS_FILLED && len(seminar.Trainees) < seminar.ClassRoom.MaxStudents {
		statusOpened, err := SeminarStatusService.GetSeminarStatusByID(model.SEMINAR_STATUS_OPENED)
		if err != nil {
			return err
		}
		seminar.SeminarStatus = *statusOpened
		seminar.SeminarStatusID = model.SEMINAR_STATUS_OPENED
		_, err = c.UpdateSeminar(*seminar)
		if err != nil {
			return err
		}

		return nil
	}

	if seminar.SeminarStatus.ID == model.SEMINAR_STATUS_OPENED && len(seminar.Trainees) >= seminar.ClassRoom.MaxStudents {
		statusFilled, err := SeminarStatusService.GetSeminarStatusByID(model.SEMINAR_STATUS_FILLED)
		if err != nil {
			return err
		}
		seminar.SeminarStatus = *statusFilled
		seminar.SeminarStatusID = model.SEMINAR_STATUS_FILLED
		_, err = c.UpdateSeminar(*seminar)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}
