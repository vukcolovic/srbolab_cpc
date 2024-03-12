package service

import (
	"errors"
	"srbolab_cpc/db"
	"srbolab_cpc/model"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

const ClientFolder = "clients"

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
	GetAllClientsWithSeminarsAndBasePersonalInfo(skip, take int) ([]model.Client, error)
}

func buildClientQuery(filter model.ClientFilter) string {
	query := ""
	if len(filter.FirstName) > 0 {
		if query != "" {
			query = query + " AND "
		}
		query = query + ` first_name ILIKE ` + "'" + filter.FirstName + "%'"
	}
	if len(filter.LastName) > 0 {
		if query != "" {
			query = query + " AND "
		}
		query = query + ` last_name ILIKE ` + "'" + filter.LastName + "%'"
	}
	if len(filter.JMBG) > 0 {
		if query != "" {
			query = query + " AND "
		}
		query = query + ` jmbg ILIKE ` + "'" + filter.JMBG + "%'"
	}
	if filter.CompanyID > 0 {
		if query != "" {
			query = query + " AND "
		}
		query = query + ` company_id = ` + strconv.Itoa(filter.CompanyID)
	}
	if len(filter.Verified) > 0 {
		if filter.Verified == "true" {
			if query != "" {
				query = query + " AND "
			}
			query = query + ` verified = true`
		} else {
			if query != "" {
				query = query + " AND "
			}
			query = query + ` verified = false`
		}
	}

	if len(filter.WaitSeminar) > 0 {
		if filter.WaitSeminar == "true" {
			if query != "" {
				query = query + " AND "
			}
			query = query + ` wait_seminar = true`
		} else {
			if query != "" {
				query = query + " AND "
			}
			query = query + ` wait_seminar = false`
		}
	}

	return query
}

func (c *clientService) GetAllClients(skip, take int, filter model.ClientFilter) ([]model.Client, error) {
	var clients []model.Client

	if filter.WaitingRoom {
		if err := db.Client.Where("verified = ?", false).Or("wait_seminar = ?", true).Limit(take).Offset(skip).Preload("Company").Find(&clients).Error; err != nil {
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

	query := buildClientQuery(filter)
	if err := db.Client.Where(query).Order("id desc").Limit(take).Offset(skip).Preload("Company").Find(&clients).Error; err != nil {
		return nil, err
	}

	return clients, nil
}

func (c *clientService) GetClientByID(id int) (*model.Client, error) {
	var client *model.Client
	if err := db.Client.Preload("Documents").Preload("Seminars").Preload("Seminars.Seminar.SeminarTheme").Preload("Seminars.Seminar.SeminarStatus").Preload("Seminars.Seminar.SeminarTheme.BaseSeminarType").Preload("Seminars.Seminar.ClassRoom.Location").Preload("Company").First(&client, id).Error; err != nil {
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

	documents := []model.File{}
	for _, doc := range client.Documents {
		documents = append(documents, model.File{Content: doc.Content, Name: doc.Name})
		doc.Content = ""
	}

	result := db.Client.Create(&client)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, doc := range documents {
		err := FileService.WriteFile(FileService.GetPath(ClientFolder, strconv.Itoa(int(client.ID)), doc.Name), doc.Content)
		if err != nil {
			return nil, err
		}
	}

	if len(client.Seminars) == 0 {
		return &client, nil
	}

	for _, s := range client.Seminars {
		err := SeminarService.UpdateSeminarStatusIfNeed(int(s.SeminarID))
		if err != nil {
			return &client, err
		}
		seminar, err := SeminarService.GetSeminarByID(int(s.SeminarID))
		if seminar.SeminarStatusID == model.SEMINAR_STATUS_IN_PROGRESS {
			SeminarDayService.AddClientToInProgressSeminar(s)
		}
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

	for _, od := range oldClient.Documents {
		found := false
		for _, nd := range client.Documents {
			if od.ID == nd.ID {
				found = true
				break
			}
		}

		if !found {
			if len(od.Content) == 0 {
				FileService.DeleteFile(FileService.GetPath(ClientFolder, strconv.Itoa(int(client.ID)), od.Name))
			}
			result := db.Client.Exec("DELETE FROM client_file WHERE client_id = ? AND file_id = ?", client.ID, od.ID)
			if result.Error != nil {
				return nil, result.Error
			}
			result = db.Client.Exec("DELETE FROM files WHERE id = ?", od.ID)
			if result.Error != nil {
				return nil, result.Error
			}
		}
	}

	for _, nd := range client.Documents {
		found := false
		for _, od := range oldClient.Documents {
			if od.ID == nd.ID {
				found = true
				break
			}
		}

		if !found {
			FileService.WriteFile(FileService.GetPath(ClientFolder, strconv.Itoa(int(client.ID)), nd.Name), nd.Content)
		}
	}

	for _, doc := range client.Documents {
		doc.Content = ""
	}

	result := db.Client.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&client)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, os := range oldSeminars {
		if os.Seminar.SeminarStatus.ID != model.SEMINAR_STATUS_OPENED && os.Seminar.SeminarStatus.ID != model.SEMINAR_STATUS_FILLED && os.Seminar.SeminarStatus.ID != model.SEMINAR_STATUS_IN_PROGRESS {
			continue
		}
		found := false
		for _, ns := range client.Seminars {
			if ns.Seminar.SeminarStatus.ID != model.SEMINAR_STATUS_OPENED && os.Seminar.SeminarStatus.ID != model.SEMINAR_STATUS_FILLED && os.Seminar.SeminarStatus.ID != model.SEMINAR_STATUS_IN_PROGRESS {
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
			if ns.Seminar.SeminarStatusID == model.SEMINAR_STATUS_IN_PROGRESS {
				SeminarDayService.AddClientToInProgressSeminar(ns)
			}
		}
	}

	return &client, nil
}

func (c *clientService) GetAllClientsWithSeminarsAndBasePersonalInfo(skip, take int) ([]model.Client, error) {
	var clients []model.Client

	if err := db.Client.Order("id desc").Limit(take).Offset(skip).Preload("Company").Preload("Seminars").Preload("Seminars.Seminar").Preload("Seminars.Seminar.ClassRoom").Preload("Seminars.Seminar.SeminarTheme").Find(&clients).Error; err != nil {
		return nil, err
	}

	return clients, nil
}
