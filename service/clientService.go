package service

import (
	"gorm.io/gorm"
	"srbolab_cpc/db"
	"srbolab_cpc/model"
)

var (
	ClientService clientServiceInterface = &clientService{}
)

type clientService struct {
}

type clientServiceInterface interface {
	GetAllClients(skip, take int) ([]model.Client, error)
	GetClientByID(id int) (*model.Client, error)
	GetClientsCount() (int64, error)
	DeleteClient(id int) error
	CreateClient(client model.Client) (*model.Client, error)
	UpdateClient(client model.Client) (*model.Client, error)
}

func (c *clientService) GetAllClients(skip, take int) ([]model.Client, error) {
	var clients []model.Client
	if err := db.Client.Order("id desc").Limit(take).Offset(skip).Find(&clients).Error; err != nil {
		return nil, err
	}
	return clients, nil
}

func (c *clientService) GetClientByID(id int) (*model.Client, error) {
	var client *model.Client
	if err := db.Client.Preload("Documents").First(&client, id).Error; err != nil {
		return nil, err
	}

	return client, nil
}

func (c *clientService) GetClientsCount() (int64, error) {
	var count int64
	if err := db.Client.Model(&model.Client{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (c *clientService) DeleteClient(id int) error {
	return db.Client.Delete(&model.Client{}, id).Error
}

func (c *clientService) CreateClient(client model.Client) (*model.Client, error) {
	result := db.Client.Create(&client)
	if result.Error != nil {
		return nil, result.Error
	}

	return &client, nil
}

func (c *clientService) UpdateClient(client model.Client) (*model.Client, error) {
	result := db.Client.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&client)
	if result.Error != nil {
		return nil, result.Error
	}

	return &client, nil
}
