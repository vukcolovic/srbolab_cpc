package service

import (
	"errors"
	"srbolab_cpc/db"
	"srbolab_cpc/model"

	"golang.org/x/exp/slices"

	"gorm.io/gorm"
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
	CreateBulk(seminarID int, clientIDs []int) error
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

func (c *clientSeminarService) CreateBulk(seminarID int, clientIDs []int) error {
	if seminarID == 0 {
		return errors.New("seminar ne može biti prazan")
	}
	if len(clientIDs) == 0 {
		return errors.New("lista klijenata je prazna")
	}

	alreadyAddedClients := []uint{}
	err := db.Client.Raw("select client_id from client_seminars cs where cs.seminar_id = ? and client_id in (?)", seminarID, clientIDs).Scan(&alreadyAddedClients).Error
	if err != nil {
		return err
	}

	tx := db.Client.Begin()

	for _, id := range clientIDs {
		if slices.Contains(alreadyAddedClients, uint(id)) {
			continue
		}
		err = tx.Raw("INSERT INTO client_seminars (seminar_id, client_id, created_at) VALUES (?, ?, now())", seminarID, id).Scan(&alreadyAddedClients).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}
