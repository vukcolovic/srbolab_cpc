package service

import (
	"gorm.io/gorm"
	"srbolab_cpc/db"
	"srbolab_cpc/model"
)

var (
	ClassRoomService classRoomServiceInterface = &classRoomService{}
)

type classRoomService struct {
}

type classRoomServiceInterface interface {
	GetAllClassRooms() ([]model.ClassRoom, error)
	GetClassRoomByID(id int) (*model.ClassRoom, error)
	CreateClassRoom(classRoom model.ClassRoom) (*model.ClassRoom, error)
	UpdateClassRoom(classRoom model.ClassRoom) (*model.ClassRoom, error)
	GetClassRoomsByLocationID(locationID int) ([]model.ClassRoom, error)
}

func (c *classRoomService) GetAllClassRooms() ([]model.ClassRoom, error) {
	var classRooms []model.ClassRoom
	if err := db.Client.Find(&classRooms).Error; err != nil {
		return nil, err
	}
	return classRooms, nil
}

func (c *classRoomService) GetClassRoomByID(id int) (*model.ClassRoom, error) {
	var classRooms *model.ClassRoom
	if err := db.Client.First(&classRooms, id).Error; err != nil {
		return nil, err
	}

	return classRooms, nil
}

func (c *classRoomService) CreateClassRoom(classRoom model.ClassRoom) (*model.ClassRoom, error) {
	result := db.Client.Create(&classRoom)
	if result.Error != nil {
		return nil, result.Error
	}

	return &classRoom, nil
}

func (c *classRoomService) UpdateClassRoom(classRoom model.ClassRoom) (*model.ClassRoom, error) {
	result := db.Client.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&classRoom)
	if result.Error != nil {
		return nil, result.Error
	}

	return &classRoom, nil
}

func (c *classRoomService) GetClassRoomsByLocationID(locationID int) ([]model.ClassRoom, error) {
	var classRooms []model.ClassRoom
	if err := db.Client.Where("location_id", locationID).Find(&classRooms).Error; err != nil {
		return nil, err
	}
	return classRooms, nil
}
