package service

import (
	"gorm.io/gorm"
	"srbolab_cpc/db"
	"srbolab_cpc/model"
)

var (
	TestService testServiceInterface = &testService{}
)

type testService struct {
}

type testServiceInterface interface {
	GetAllTests() ([]model.Test, error)
	GetTestByID(id int) (*model.Test, error)
	CreateTest(test model.Test) (*model.Test, error)
	UpdateTest(test model.Test) (*model.Test, error)
}

func (c *testService) GetAllTests() ([]model.Test, error) {
	var tests []model.Test
	if err := db.Client.Order("id desc").Preload("SeminarTheme").Preload("SeminarTheme.BaseSeminarType").Find(&tests).Error; err != nil {
		return nil, err
	}
	return tests, nil
}

func (c *testService) GetTestByID(id int) (*model.Test, error) {
	var test *model.Test
	if err := db.Client.Preload("SeminarTheme").Preload("SeminarTheme.BaseSeminarType").Preload("Questions").First(&test, id).Error; err != nil {
		return nil, err
	}

	return test, nil
}

func (c *testService) CreateTest(test model.Test) (*model.Test, error) {
	result := db.Client.Create(&test)
	if result.Error != nil {
		return nil, result.Error
	}

	return &test, nil
}

func (c *testService) UpdateTest(test model.Test) (*model.Test, error) {
	oldTest, err := c.GetTestByID(int(test.ID))
	if err != nil {
		return nil, err
	}

	result := db.Client.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&test)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, od := range oldTest.Questions {
		found := false
		for _, nd := range test.Questions {
			if od.ID == nd.ID {
				found = true
				break
			}
		}

		if !found {
			result := db.Client.Exec("DELETE FROM test_question WHERE test_id = ? AND question_id = ?", test.ID, od.ID)
			if result.Error != nil {
				return nil, result.Error
			}
		}
	}

	return &test, nil
}
