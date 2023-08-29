package service

import (
	"gorm.io/gorm"
	"sort"
	"srbolab_cpc/db"
	"srbolab_cpc/model"
	"srbolab_cpc/util"
	"strconv"
)

var (
	TestService testServiceInterface = &testService{}
)

type testService struct {
}

type testServiceInterface interface {
	GetAllTests() ([]model.Test, error)
	GetAllTestsBySeminarTheme(seminarThemeID int) ([]model.Test, error)
	GetTestByID(id int) (*model.Test, error)
	CreateTest(test model.Test) (*model.Test, error)
	UpdateTest(test model.Test) (*model.Test, error)
	CreateClientTest(clientTest model.ClientTest) (*model.ClientTest, error)
	GetClientTestBySeminarDayIDAndJMBG(seminarDayID int, jmbg string) ([]model.ClientTest, error)
}

func (c *testService) GetAllTests() ([]model.Test, error) {
	var tests []model.Test
	if err := db.Client.Order("id desc").Preload("SeminarTheme").Preload("SeminarTheme.BaseSeminarType").Find(&tests).Error; err != nil {
		return nil, err
	}
	return tests, nil
}

func (c *testService) GetAllTestsBySeminarTheme(seminarThemeID int) ([]model.Test, error) {
	var tests []model.Test
	if err := db.Client.Order("id desc").Where("seminar_theme_id = ?", seminarThemeID).Find(&tests).Error; err != nil {
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

func (c *testService) CreateClientTest(clientTest model.ClientTest) (*model.ClientTest, error) {
	res := ""
	mapAnswers := map[uint]string{}
	sort.Slice(clientTest.QuestionAnswer, func(i, j int) bool {
		return clientTest.QuestionAnswer[i].QuestionID < clientTest.QuestionAnswer[i].QuestionID
	})
	for _, qa := range clientTest.QuestionAnswer {
		mapAnswers[qa.QuestionID] = qa.Answer
		res = res + strconv.Itoa(int(qa.QuestionID)) + ":" + qa.Answer + ","
	}
	res = util.TrimSuffix(res, ",")
	clientTest.ResultStr = res

	correctAnswers := 0

	for _, q := range clientTest.Test.Questions {
		for _, a := range q.Answers {
			if a.Correct != nil && *a.Correct == true {
				v, ok := mapAnswers[q.ID]
				if !ok || v != a.Letter {
					break
				} else {
					correctAnswers++
				}
			}
		}
	}
	clientTest.Result = float64(correctAnswers) / float64(len(clientTest.Test.Questions))

	result := db.Client.Create(&clientTest)
	if result.Error != nil {
		return nil, result.Error
	}

	return &clientTest, nil
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

func (c *testService) GetClientTestBySeminarDayIDAndJMBG(seminarDayID int, jmbg string) ([]model.ClientTest, error) {
	client, err := ClientService.GetClientByJMBG(jmbg)
	if err != nil {
		return []model.ClientTest{}, err
	}

	var clientTests []model.ClientTest
	if err := db.Client.Where("seminar_day_id = ? AND client_id = ?", seminarDayID, client.ID).Find(&clientTests).Error; err != nil {
		return nil, err
	}

	return clientTests, nil
}
