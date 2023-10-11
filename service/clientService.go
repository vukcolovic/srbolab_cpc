package service

import (
	"errors"
	"gorm.io/gorm"
	"srbolab_cpc/db"
	"srbolab_cpc/model"
	"strings"
)

var (
	ClientService clientServiceInterface = &clientService{}
)

type clientService struct {
}

type clientServiceInterface interface {
	GetAllClients(skip, take int, filter model.ClientFilter) ([]model.Client, error)
	GetClientByID(id int) (*model.Client, error)
	GetClientByJMBG(jmbg string) (*model.Client, error)
	GetClientByJMBGWithSeminars(jmbg string) (*model.Client, error)
	GetClientsCount() (int64, error)
	DeleteClient(id int) error
	CreateClient(client model.Client, userID int) (*model.Client, error)
	CreateClientNotVerified(client model.Client) (*model.Client, error)
	UpdateClient(client model.Client, userID int) (*model.Client, error)
}

func makeFilterMap(filter model.ClientFilter) map[string]interface{} {
	filterMap := make(map[string]interface{})
	if len(filter.FirstName) > 0 {
		filterMap["first_name"] = filter.FirstName
	}
	if len(filter.LastName) > 0 {
		filterMap["last_name"] = filter.LastName
	}
	if len(filter.JMBG) > 0 {
		filterMap["jmbg"] = filter.JMBG
	}

	if len(filter.Verified) > 0 {
		if filter.Verified == "true" {
			filterMap["verified"] = true
		} else {
			filterMap["verified"] = false
		}
	}

	if len(filter.WaitSeminar) > 0 {
		if filter.WaitSeminar == "true" {
			filterMap["wait_seminar"] = true
		} else {
			filterMap["wait_seminar"] = false
		}
	}

	return filterMap
}

func (c *clientService) GetAllClients(skip, take int, filter model.ClientFilter) ([]model.Client, error) {
	var clients []model.Client

	if filter.WaitingRoom {
		if err := db.Client.Where("verified = ?", false).Or("wait_seminar = ?", true).Limit(take).Offset(skip).Find(&clients).Error; err != nil {
			return nil, err
		}

		for i, c := range clients {
			if filter.FirstName != "" && !strings.Contains(c.Person.FirstName, filter.FirstName) {
				clients = append(clients[:i], clients[i+1:]...)
			}
			if filter.LastName != "" && !strings.Contains(c.Person.LastName, filter.LastName) {
				clients = append(clients[:i], clients[i+1:]...)
			}
			if filter.JMBG != "" && !strings.Contains(*c.JMBG, filter.JMBG) {
				clients = append(clients[:i], clients[i+1:]...)
			}
		}
		return clients, nil
	}

	filterMap := makeFilterMap(filter)
	if err := db.Client.Where(filterMap).Order("id desc").Limit(take).Offset(skip).Find(&clients).Error; err != nil {
		return nil, err
	}

	return clients, nil
}

func (c *clientService) GetClientByID(id int) (*model.Client, error) {
	var client *model.Client
	if err := db.Client.Preload("Documents").Preload("Seminars").Preload("Seminars.Seminar.SeminarTheme").Preload("Seminars.Seminar.SeminarStatus").Preload("Seminars.Seminar.SeminarTheme.BaseSeminarType").Preload("Company").First(&client, id).Error; err != nil {
		return nil, err
	}

	return client, nil
}

func (c *clientService) GetClientByJMBG(jmbg string) (*model.Client, error) {
	var client *model.Client
	if err := db.Client.Where("jmbg", jmbg).First(&client).Error; err != nil {
		return nil, err
	}

	return client, nil
}

func (c *clientService) GetClientByJMBGWithSeminars(jmbg string) (*model.Client, error) {
	var client *model.Client
	if err := db.Client.Preload("Seminars").Preload("Seminars.Seminar.Days").Where("jmbg", jmbg).First(&client).Error; err != nil {
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

func (c *clientService) CreateClient(client model.Client, userID int) (*model.Client, error) {
	userIDUint := uint(userID)
	client.CreatedByID = &userIDUint
	if *client.Verified {
		client.VerifiedByID = &userIDUint
	}
	result := db.Client.Create(&client)
	if result.Error != nil {
		return nil, result.Error
	}

	if len(client.Seminars) > 0 {
		err := SeminarService.UpdateSeminarStatusIfNeed(int(client.Seminars[0].SeminarID))
		if err != nil {
			return &client, err
		}
	}

	return &client, nil
}

func (c *clientService) CreateClientNotVerified(client model.Client) (*model.Client, error) {
	cl, _ := c.GetClientByJMBG(*client.JMBG)
	if cl != nil {
		return nil, errors.New("Osoba sa datim JMBG-om već postoji u sistemu")
	}
	result := db.Client.Create(&client)
	if result.Error != nil {
		return nil, result.Error
	}

	return &client, nil
}

func (c *clientService) UpdateClient(client model.Client, userID int) (*model.Client, error) {
	oldClient, err := c.GetClientByID(int(client.ID))
	if err != nil {
		return nil, err
	}

	userIDUint := uint(userID)

	if !*oldClient.Verified && *client.Verified {
		client.VerifiedByID = &userIDUint
	}
	oldSeminars := oldClient.Seminars

	result := db.Client.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&client)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, od := range oldClient.Documents {
		found := false
		for _, nd := range client.Documents {
			if od.ID == nd.ID {
				found = true
				break
			}
		}

		if !found {
			result := db.Client.Exec("DELETE FROM client_file WHERE client_id = ? AND file_id = ?", client.ID, od.ID)
			if result.Error != nil {
				return nil, result.Error
			}
		}
	}

	for _, os := range oldSeminars {
		if os.Seminar.SeminarStatus.ID != model.SEMINAR_STATUS_OPENED && os.Seminar.SeminarStatus.ID != model.SEMINAR_STATUS_FILLED {
			continue
		}
		found := false
		for _, ns := range client.Seminars {
			if ns.Seminar.SeminarStatus.ID != model.SEMINAR_STATUS_OPENED && os.Seminar.SeminarStatus.ID != model.SEMINAR_STATUS_FILLED {
				continue
			}
			if os.ClientID == ns.ClientID && os.SeminarID == ns.SeminarID {
				found = true
				break
			}
		}

		if !found {
			SeminarService.DeleteSeminarClient(os)
		}
	}

	if len(client.Seminars) > len(oldSeminars) {
		updatedClient, err := c.GetClientByID(int(client.ID))
		if err != nil {
			return nil, err
		}
		for _, ns := range updatedClient.Seminars {
			found := false
			for _, os := range oldSeminars {
				if os.ClientID == ns.ClientID && os.SeminarID == ns.SeminarID {
					found = true
					break
				}
			}

			if !found {
				err = SeminarService.UpdateSeminarStatusIfNeed(int(ns.SeminarID))

			}
		}
	}

	return &client, nil
}
