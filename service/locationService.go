package service

import (
	"gorm.io/gorm"
	"srbolab_cpc/db"
	"srbolab_cpc/model"
)

var (
	LocationService locationServiceInterface = &locationService{}
)

type locationService struct {
}

type locationServiceInterface interface {
	GetAllLocations() ([]model.Location, error)
	GetLocationByID(id int) (*model.Location, error)
	GetLocationsCount() (int64, error)
	CreateLocation(location model.Location) (*model.Location, error)
	UpdateLocation(location model.Location) (*model.Location, error)
}

func (c *locationService) GetAllLocations() ([]model.Location, error) {
	var locations []model.Location
	if err := db.Client.Find(&locations).Error; err != nil {
		return nil, err
	}
	return locations, nil
}

func (c *locationService) GetLocationByID(id int) (*model.Location, error) {
	var location *model.Location
	if err := db.Client.First(&location, id).Error; err != nil {
		return nil, err
	}

	return location, nil
}

func (c *locationService) GetLocationsCount() (int64, error) {
	var count int64
	if err := db.Client.Model(&model.Location{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (c *locationService) CreateLocation(location model.Location) (*model.Location, error) {
	result := db.Client.Create(&location)
	if result.Error != nil {
		return nil, result.Error
	}

	return &location, nil
}

func (c *locationService) UpdateLocation(location model.Location) (*model.Location, error) {
	result := db.Client.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&location)
	if result.Error != nil {
		return nil, result.Error
	}

	return &location, nil
}
