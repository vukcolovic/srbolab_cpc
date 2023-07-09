package service

import (
	"srbolab_cpc/db"
	"srbolab_cpc/model"
)

var (
	CompanyService companyServiceInterface = &companyService{}
)

type companyService struct {
}

type companyServiceInterface interface {
	GetAllCompanies(skip, take int) ([]model.Company, error)
	GetCompanyByID(id int) (*model.Company, error)
	GetCompaniesCount() (int64, error)
	DeleteCompany(id int) error
	CreateCompany(company model.Company) (*model.Company, error)
	UpdateCompany(company model.Company) (*model.Company, error)
}

func (c *companyService) GetAllCompanies(skip, take int) ([]model.Company, error) {
	var companies []model.Company
	if err := db.Client.Order("id desc").Limit(take).Offset(skip).Find(&companies).Error; err != nil {
		return nil, err
	}
	return companies, nil
}

func (c *companyService) GetCompanyByID(id int) (*model.Company, error) {
	var companies *model.Company
	if err := db.Client.First(&companies, id).Error; err != nil {
		return nil, err
	}

	return companies, nil
}

func (c *companyService) GetCompaniesCount() (int64, error) {
	var count int64
	if err := db.Client.Model(&model.Company{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (c *companyService) DeleteCompany(id int) error {
	return db.Client.Delete(&model.Company{}, id).Error
}

func (c *companyService) CreateCompany(company model.Company) (*model.Company, error) {
	result := db.Client.Create(&company)
	if result.Error != nil {
		return nil, result.Error
	}

	return &company, nil
}

func (c *companyService) UpdateCompany(company model.Company) (*model.Company, error) {
	result := db.Client.Updates(&company)
	if result.Error != nil {
		return nil, result.Error
	}

	return &company, nil
}
