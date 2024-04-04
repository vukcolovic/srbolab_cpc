package service

import (
	"errors"
	"srbolab_cpc/db"
	"srbolab_cpc/logoped"
	"srbolab_cpc/model"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

const SeminarFolder = "seminars"

var (
	SeminarService seminarServiceInterface = &seminarService{}
)

type seminarService struct {
}

type seminarServiceInterface interface {
	GetAllSeminars(skip, take int, filter model.SeminarFilter) ([]model.Seminar, error)
	GetAllSeminarsWithTrainees(skip, take int, filter model.SeminarFilter) ([]model.Seminar, error)
	GetAllSeminarsWithSeminarDays(skip, take int, filter model.SeminarFilter) ([]model.Seminar, error)
	GetAllSeminarsByStatus(statusCode string) ([]model.Seminar, error)
	GetSeminarByID(id int) (*model.Seminar, error)
	GetSeminarByIDWithClientFiles(id int) (*model.Seminar, error)
	GetSeminarsCount() (int64, error)
	DeleteSeminar(id int) error
	CreateSeminar(seminar model.Seminar) (*model.Seminar, error)
	UpdateSeminar(seminar model.Seminar) (*model.Seminar, error)
	DeleteSeminarClient(clientSeminar model.ClientSeminar) error
	UpdateSeminarStatusIfNeed(seminarID int) error
	GetClientSeminarBySeminarIDAndClientID(seminarID, clientID uint) (*model.ClientSeminar, error)
}

func (c *seminarService) GetAllSeminars(skip, take int, filter model.SeminarFilter) ([]model.Seminar, error) {
	var seminars []model.Seminar
	query := buildQuery(filter)
	if err := db.Client.Where(query).Order("start desc").Limit(take).Offset(skip).Preload("ClassRoom.Location").Joins("ClassRoom").Joins("SeminarTheme").Joins("SeminarTheme.BaseSeminarType").Joins("SeminarStatus").Find(&seminars).Error; err != nil {
		logoped.ErrorLog.Println("Error getting all seminars " + err.Error())
		return nil, err
	}
	return seminars, nil
}

func (c *seminarService) GetAllSeminarsWithTrainees(skip, take int, filter model.SeminarFilter) ([]model.Seminar, error) {
	var seminars []model.Seminar
	query := buildQuery(filter)
	if err := db.Client.Where(query).Order("start desc").Limit(take).Offset(skip).Preload("ClassRoom.Location").Joins("ClassRoom").Joins("SeminarTheme").Joins("SeminarTheme.BaseSeminarType").Joins("SeminarStatus").Preload("Trainees").Find(&seminars).Error; err != nil {
		logoped.ErrorLog.Println("Error getting all seminars " + err.Error())
		return nil, err
	}
	return seminars, nil
}

func (c *seminarService) GetAllSeminarsWithSeminarDays(skip, take int, filter model.SeminarFilter) ([]model.Seminar, error) {
	var seminars []model.Seminar
	query := buildQuery(filter)
	if err := db.Client.Where(query).Order("start desc").Limit(take).Offset(skip).Preload("ClassRoom.Location").Joins("ClassRoom").Joins("SeminarTheme").Joins("SeminarTheme.BaseSeminarType").Joins("SeminarStatus").Preload("Days").Preload("Days.Classes").Preload("Days.Classes.Teacher").Find(&seminars).Error; err != nil {
		logoped.ErrorLog.Println("Error getting all seminars " + err.Error())
		return nil, err
	}
	return seminars, nil
}

func buildQuery(filter model.SeminarFilter) string {
	query := ""
	if filter.LocationID > 0 {
		classRooms, err := ClassRoomService.GetClassRoomsByLocationID(filter.LocationID)
		if err != nil {
			logoped.ErrorLog.Println("Error making seminar query: " + err.Error())
		}
		query = query + "class_room_id IN ("
		for _, c := range classRooms {
			query = query + strconv.Itoa(int(c.ID)) + ", "
		}

		query, _ = strings.CutSuffix(query, ", ")
		query = query + ")"
	}
	if !filter.DateFrom.IsZero() {
		if query != "" {
			query = query + " AND "
		}
		query = query + "start >= " + "'" + filter.DateFrom.Format("2006-01-02") + "'" + "::date"
	}
	if !filter.DateTo.IsZero() {
		if query != "" {
			query = query + " AND "
		}
		query = query + "start < " + "'" + filter.DateTo.Format("2006-01-02") + "'" + "::date + " + "'1 day'::interval"
	}

	return query
}

func (c *seminarService) GetAllSeminarsByStatus(statusCode string) ([]model.Seminar, error) {
	status, err := SeminarStatusService.GetSeminarStatusByCode(statusCode)
	if err != nil {
		logoped.ErrorLog.Println("Error getting all seminars by status, error get seminar status by code " + err.Error())
		return nil, err
	}
	var seminars []model.Seminar
	if err := db.Client.Where("seminar_status_id", status.ID).Order("id desc").Preload("ClassRoom.Location").Joins("ClassRoom").Joins("SeminarTheme").Joins("SeminarTheme.BaseSeminarType").Joins("SeminarStatus").Find(&seminars).Error; err != nil {
		logoped.ErrorLog.Println("Error getting all seminars by status " + err.Error())
		return nil, err
	}
	return seminars, nil
}

func (c *seminarService) GetSeminarByID(id int) (*model.Seminar, error) {
	var seminar *model.Seminar
	if err := db.Client.Preload("Trainees").Preload("Trainees.Client").Preload("Trainees.Client.Company").Preload("Days").Preload("Days.Presence").Preload("Documents").Preload("ClassRoom.Location").Preload("Admin").Joins("SeminarTheme.BaseSeminarType").Joins("SeminarStatus").First(&seminar, id).Error; err != nil {
		logoped.ErrorLog.Println("Error getting seminar by id " + err.Error())
		return nil, err
	}

	return seminar, nil
}

func (c *seminarService) GetSeminarByIDWithClientFiles(id int) (*model.Seminar, error) {
	var seminar *model.Seminar
	if err := db.Client.Preload("Trainees").Preload("Trainees.Client").Preload("Trainees.Client.Company").Preload("Trainees.Client.Documents").Preload("ClassRoom.Location").Joins("SeminarTheme.BaseSeminarType").Joins("SeminarStatus").First(&seminar, id).Error; err != nil {
		logoped.ErrorLog.Println("Error getting seminar by id " + err.Error())
		return nil, err
	}

	return seminar, nil
}

func (c *seminarService) GetSeminarsCount() (int64, error) {
	var count int64
	if err := db.Client.Model(&model.Seminar{}).Count(&count).Error; err != nil {
		logoped.ErrorLog.Println("Error getting seminars count " + err.Error())
		return 0, err
	}

	return count, nil
}

func (c *seminarService) DeleteSeminar(id int) error {
	err := db.Client.Delete(&model.Seminar{}, id).Error
	if err != nil {
		logoped.ErrorLog.Println("Error deleting seminar " + err.Error())
		return err
	}
	return nil
}

func (c *seminarService) CreateSeminar(seminar model.Seminar) (*model.Seminar, error) {
	if seminar.SerialNumberByLocation > 0 {
		countByLocation := 0
		db.Client.Raw("SELECT COUNT(*) FROM seminars s JOIN class_rooms cr ON s.class_room_id = cr.id JOIN seminar_themes st ON st.id = s.seminar_theme_id WHERE s.deleted_at is null AND cr.location_id = ? AND s.serial_number_by_location = ? AND st.base_seminar_type_id = ?", seminar.ClassRoom.LocationID, seminar.SerialNumberByLocation, seminar.SeminarTheme.BaseSeminarTypeID).Scan(&countByLocation)
		if countByLocation > 0 {
			return nil, errors.New("već postoji seminar sa ovim brojem na izabranoj lokaciji")
		}
	}

	result := db.Client.Create(&seminar)
	if result.Error != nil {
		logoped.ErrorLog.Println("Error creating seminar " + result.Error.Error())
		return nil, result.Error
	}

	return &seminar, nil
}

func (c *seminarService) UpdateSeminar(seminar model.Seminar) (*model.Seminar, error) {
	if seminar.SerialNumberByLocation > 0 {
		countByLocation := 0
		db.Client.Raw("SELECT COUNT(*) FROM seminars s JOIN class_rooms cr ON s.class_room_id = cr.id JOIN seminar_themes st ON st.id = s.seminar_theme_id WHERE s.deleted_at is null AND cr.location_id = ? AND s.serial_number_by_location = ? AND s.id <> ? AND st.base_seminar_type_id = ?", seminar.ClassRoom.LocationID, seminar.SerialNumberByLocation, seminar.ID, seminar.SeminarTheme.BaseSeminarTypeID).Scan(&countByLocation)
		if countByLocation > 0 {
			logoped.ErrorLog.Println("Error updating seminar, already exists seminar number on this location.")
			return nil, errors.New("već postoji seminar sa ovim brojem na izabranoj lokaciji")
		}
	}

	oldSeminar, err := c.GetSeminarByID(int(seminar.ID))
	if err != nil {
		logoped.ErrorLog.Println("Error updating seminar, error getting seminar by id.")
		return nil, err
	}

	if seminar.SeminarStatus.ID == model.SEMINAR_STATUS_CLOSED {
		notPassed := make(map[uint]string)
		for _, day := range seminar.Days {
			for _, pr := range day.Presence {
				if !*pr.Presence {
					notPassed[pr.ClientID] = ""
				}
			}
		}

		for i, pr := range seminar.Trainees {
			if _, ok := notPassed[pr.ClientID]; !ok {
				b := true
				seminar.Trainees[i].Pass = &b
				switch seminar.SeminarTheme.Code {
				case "1":
					seminar.Trainees[i].Client.PassedCheckboxes.WorkTimeAndTahografs = &b
				case "2":
					seminar.Trainees[i].Client.PassedCheckboxes.ThemeDocuments = &b
				case "3":
					seminar.Trainees[i].Client.PassedCheckboxes.Burden = &b
				case "4":
					seminar.Trainees[i].Client.PassedCheckboxes.Regulations = &b
				case "5":
					seminar.Trainees[i].Client.PassedCheckboxes.Tahografs2 = &b
				}
			} else {
				b := false
				seminar.Trainees[i].Pass = &b
			}
		}
	}

	for _, od := range oldSeminar.Documents {
		found := false
		for _, nd := range seminar.Documents {
			if od.ID == nd.ID {
				found = true
				break
			}
		}

		if !found {
			if len(od.Content) == 0 {
				FileService.DeleteFile(FileService.GetPath(SeminarFolder, strconv.Itoa(int(seminar.ID)), od.Name))
			}
			result := db.Client.Exec("DELETE FROM seminar_file WHERE seminar_id = ? AND file_id = ?", seminar.ID, od.ID)
			if result.Error != nil {
				logoped.ErrorLog.Println("Error updating seminar, delete from seminar file")
				return nil, result.Error
			}
			result = db.Client.Exec("DELETE FROM files WHERE id = ?", od.ID)
			if result.Error != nil {
				return nil, result.Error
			}
		}
	}

	for _, nd := range seminar.Documents {
		found := false
		for _, od := range oldSeminar.Documents {
			if od.ID == nd.ID {
				found = true
				break
			}
		}

		if !found {
			FileService.WriteFile(FileService.GetPath(SeminarFolder, strconv.Itoa(int(seminar.ID)), nd.Name), nd.Content)
		}
	}

	for _, doc := range seminar.Documents {
		doc.Content = ""
	}

	result := db.Client.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&seminar)
	if result.Error != nil {
		logoped.ErrorLog.Println("Error updating seminar")
		return nil, result.Error
	}

	return &seminar, nil
}

func (c *seminarService) DeleteSeminarClient(clientSeminar model.ClientSeminar) error {
	result := db.Client.Exec("DELETE FROM client_seminars WHERE client_id = ? AND seminar_id = ?", clientSeminar.ClientID, clientSeminar.SeminarID)
	if result.Error != nil {
		logoped.ErrorLog.Println("Error deleting client seminar " + result.Error.Error())
		return result.Error
	}

	if clientSeminar.Seminar.SeminarStatusID == model.SEMINAR_STATUS_IN_PROGRESS {
		result := db.Client.Exec("DELETE FROM client_presences WHERE client_id = ? AND seminar_day_id IN (SELECT id FROM seminar_days WHERE seminar_id = ?)", clientSeminar.ClientID, clientSeminar.SeminarID)
		if result.Error != nil {
			logoped.ErrorLog.Println("Error deleting client seminar, error deleting client presence " + result.Error.Error())
			return result.Error
		}
	}

	return c.UpdateSeminarStatusIfNeed(int(clientSeminar.Seminar.ID))
}

func (c *seminarService) UpdateSeminarStatusIfNeed(seminarID int) error {
	seminar, err := c.GetSeminarByID(seminarID)
	if err != nil {
		logoped.ErrorLog.Println("Error updating seminar status if need, error get seminar by id " + err.Error())
		return err
	}

	if seminar.SeminarStatus.ID == model.SEMINAR_STATUS_FILLED && len(seminar.Trainees) < seminar.ClassRoom.MaxStudents {
		statusOpened, err := SeminarStatusService.GetSeminarStatusByID(model.SEMINAR_STATUS_OPENED)
		if err != nil {
			logoped.ErrorLog.Println("Error updating seminar status if need, error get seminar status by id " + err.Error())
			return err
		}
		seminar.SeminarStatus = *statusOpened
		seminar.SeminarStatusID = model.SEMINAR_STATUS_OPENED
		_, err = c.UpdateSeminar(*seminar)
		if err != nil {
			logoped.ErrorLog.Println("Error updating seminar status if need " + err.Error())
			return err
		}

		return nil
	}

	if seminar.SeminarStatus.ID == model.SEMINAR_STATUS_OPENED && len(seminar.Trainees) >= seminar.ClassRoom.MaxStudents {
		statusFilled, err := SeminarStatusService.GetSeminarStatusByID(model.SEMINAR_STATUS_FILLED)
		if err != nil {
			return err
		}
		seminar.SeminarStatus = *statusFilled
		seminar.SeminarStatusID = model.SEMINAR_STATUS_FILLED
		_, err = c.UpdateSeminar(*seminar)
		if err != nil {
			logoped.ErrorLog.Println("Error updating seminar status if need " + err.Error())
			return err
		}

		return nil
	}

	return nil
}

func (c *seminarService) GetClientSeminarBySeminarIDAndClientID(seminarID, clientID uint) (*model.ClientSeminar, error) {
	var clientSeminar *model.ClientSeminar
	if err := db.Client.Where("client_id = ? AND seminar_id = ?", clientID, seminarID).First(&clientSeminar).Error; err != nil {
		logoped.ErrorLog.Println("Error getting ClientSeminar by SeminarID and ClientID" + err.Error())
		return nil, err
	}

	return clientSeminar, nil
}
