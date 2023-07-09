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
	GetClassRoomsByLocationID(locationID int) ([]model.ClassRoom, error)
}

func (c *locationService) GetAllLocations() ([]model.Location, error) {
	var locations []model.Location
	if err := db.Client.Find(&locations).Error; err != nil {
		return nil, err
	}
	return locations, nil
}

func (c *locationService) GetClassRoomsByLocationID(locationID int) ([]model.ClassRoom, error) {
	var classRooms []model.ClassRoom
	if err := db.Client.Where("location_id", locationID).Find(&classRooms).Error; err != nil {
		return nil, err
	}
	return classRooms, nil
}
