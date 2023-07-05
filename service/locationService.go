package service

import (
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
}

func (c *locationService) GetAllLocations() ([]model.Location, error) {
	var locations []model.Location
	if err := db.Client.Find(&locations).Error; err != nil {
		return nil, err
	}
	return locations, nil
}
