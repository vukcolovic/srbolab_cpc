package service

import (
	"srbolab_cpc/db"
	"srbolab_cpc/model"
)

var (
	PartnerService partnerServiceInterface = &partnerService{}
)

type partnerService struct {
}

type partnerServiceInterface interface {
	GetAllPartners(skip, take int) ([]model.Partner, error)
	GetPartnerByID(id int) (*model.Partner, error)
	GetPartnersCount() (int64, error)
	DeletePartner(id int) error
	CreatePartner(partner model.Partner) (*model.Partner, error)
	UpdatePartner(partner model.Partner) (*model.Partner, error)
}

func (c *partnerService) GetAllPartners(skip, take int) ([]model.Partner, error) {
	var partners []model.Partner
	if err := db.Client.Order("id desc").Limit(take).Offset(skip).Find(&partners).Error; err != nil {
		return nil, err
	}
	return partners, nil
}

func (c *partnerService) GetPartnerByID(id int) (*model.Partner, error) {
	var partner *model.Partner
	if err := db.Client.First(&partner, id).Error; err != nil {
		return nil, err
	}

	return partner, nil
}

func (c *partnerService) GetPartnersCount() (int64, error) {
	var count int64
	if err := db.Client.Model(&model.Partner{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (c *partnerService) DeletePartner(id int) error {
	return db.Client.Delete(&model.Partner{}, id).Error
}

func (c *partnerService) CreatePartner(partner model.Partner) (*model.Partner, error) {
	result := db.Client.Create(&partner)
	if result.Error != nil {
		return nil, result.Error
	}

	return &partner, nil
}

func (c *partnerService) UpdatePartner(partner model.Partner) (*model.Partner, error) {
	result := db.Client.Updates(&partner)
	if result.Error != nil {
		return nil, result.Error
	}

	return &partner, nil
}
