package service

import (
	"gorm.io/gorm"
	"srbolab_cpc/db"
	"srbolab_cpc/model"
)

var (
	ClientSeminarService clientSeminarServiceInterface = &clientSeminarService{}
)

type clientSeminarService struct {
}

type clientSeminarServiceInterface interface {
	UpdateClientSeminar(clientSeminar model.ClientSeminar) (*model.ClientSeminar, error)
	GetMaxConfirmationNumber() (int, error)
	GetNumberOfPassedSeminars(clientID uint) (int, error)
}

func (c *clientSeminarService) UpdateClientSeminar(clientSeminar model.ClientSeminar) (*model.ClientSeminar, error) {
	result := db.Client.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&clientSeminar)
	if result.Error != nil {
		return nil, result.Error
	}

	return &clientSeminar, nil
}

func (c *clientSeminarService) GetMaxConfirmationNumber() (int, error) {
	var max *int
	err := db.Client.Raw("SELECT MAX(confirmation_number) FROM client_seminars").Scan(&max).Error
	if err != nil {
		return 0, err
	}

	if max == nil {
		return 0, nil
	}

	return *max, nil
}

func (c *clientSeminarService) GetNumberOfPassedSeminars(clientID uint) (int, error) {
	var count *int
	err := db.Client.Raw("select count(*) from client_seminars cs where cs.client_id = ? and (cs.confirmation_number  > 0 or cs.pass = true)", clientID).Scan(&count).Error
	if err != nil {
		return 0, err
	}

	if count == nil {
		return 0, nil
	}

	return *count, nil
}
