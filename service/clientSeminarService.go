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
}

func (c *clientSeminarService) UpdateClientSeminar(clientSeminar model.ClientSeminar) (*model.ClientSeminar, error) {
	result := db.Client.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&clientSeminar)
	if result.Error != nil {
		return nil, result.Error
	}

	return &clientSeminar, nil
}
