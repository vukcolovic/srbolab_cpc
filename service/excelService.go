package service

import (
	"bytes"
	"fmt"
	"sort"
	"srbolab_cpc/model"
	"srbolab_cpc/util"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

var (
	ExcelService excelServiceInterface = &excelService{}
)

type excelService struct {
}

type excelServiceInterface interface {
	CreateClientTestsBySeminarDayReport(tests []model.ClientTest) ([]byte, error)
	CreateListClientsBySeminarDayReport(seminarDay *model.SeminarDay) ([]byte, error)
	CreateSeminarsReportOfClients(seminars []model.Seminar) ([]byte, error)
	CreateSeminarsReportOfTeachers(seminars []model.Seminar) ([]byte, error)
	CreateClientsReport() ([]byte, error)
}

func (excelService) CreateClientTestsBySeminarDayReport(tests []model.ClientTest) ([]byte, error) {
	exc := excelize.NewFile()
	firstRowStyle, err := exc.NewStyle(&excelize.Style{
		Border: []excelize.Border{excelize.Border{Type: "left", Style: 1, Color: "#000000"}, excelize.Border{Type: "right", Style: 1, Color: "#000000"}, excelize.Border{Type: "bottom", Style: 1, Color: "#000000"}, excelize.Border{Type: "top", Style: 1, Color: "#000000"}},
		Fill:   excelize.Fill{Type: "pattern", Color: []string{"#ADD8E6"}, Pattern: 1},
	})
	if err != nil {
		return []byte{}, err
	}
	err = exc.SetRowStyle("Sheet1", 1, 1, firstRowStyle)
	if err != nil {
		return []byte{}, err
	}

	allRowStyle, err := exc.NewStyle(&excelize.Style{
		Border: []excelize.Border{excelize.Border{Type: "left", Style: 1, Color: "#000000"}, excelize.Border{Type: "right", Style: 1, Color: "#000000"}, excelize.Border{Type: "bottom", Style: 1, Color: "#000000"}, excelize.Border{Type: "top", Style: 1, Color: "#000000"}},
	})
	if err != nil {
		return []byte{}, err
	}
	err = exc.SetRowStyle("Sheet1", 2, 10000, allRowStyle)
	if err != nil {
		return []byte{}, err
	}
	exc.SetRowHeight("Sheet1", 1, 25.0)

	exc.SetColWidth("Sheet1", "A", "B", 25.0)
	exc.SetColWidth("Sheet1", "C", "D", 25.0)
	exc.SetColWidth("Sheet1", "E", "E", 15.0)
	exc.SetColWidth("Sheet1", "F", "F", 25.0)

	exc.SetCellValue("Sheet1", "A1", "Име и презиме")
	exc.SetCellValue("Sheet1", "B1", "Семинар (тема)")
	exc.SetCellValue("Sheet1", "C1", "Дан")
	exc.SetCellValue("Sheet1", "D1", "Време теста")
	exc.SetCellValue("Sheet1", "E1", "Резултати (%)")
	exc.SetCellValue("Sheet1", "F1", "Одговори")

	for i, t := range tests {
		exc.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), t.Client.Person.FullName())
		exc.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), t.SeminarDay.Seminar.SeminarTheme.GetSeminarThemeWithBaseTheme())
		exc.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), t.SeminarDay.Number)
		exc.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), t.CreatedAt.Format("02.01.2006 15:4"))
		exc.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), fmt.Sprintf("%.2f", t.Result*100))
		exc.SetCellValue("Sheet1", "F"+strconv.Itoa(i+2), t.ResultStr)
	}

	var buf bytes.Buffer
	err = exc.Write(&buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (excelService) CreateListClientsBySeminarDayReport(seminarDay *model.SeminarDay) ([]byte, error) {
	exc := excelize.NewFile()
	firstRowStyle, err := exc.NewStyle(&excelize.Style{
		Border:    []excelize.Border{excelize.Border{Type: "center", Style: 1, Color: "#000000"}, excelize.Border{Type: "right", Style: 1, Color: "#000000"}, excelize.Border{Type: "left", Style: 1, Color: "#000000"}, excelize.Border{Type: "bottom", Style: 1, Color: "#000000"}, excelize.Border{Type: "top", Style: 1, Color: "#000000"}},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#D8D8D8"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	if err != nil {
		return []byte{}, err
	}

	err = exc.SetCellStyle("Sheet1", "B3", "I3", firstRowStyle)
	if err != nil {
		return []byte{}, err
	}

	err = exc.SetCellStyle("Sheet1", "B3", "B"+strconv.Itoa(3+len(seminarDay.Presence)), firstRowStyle)
	if err != nil {
		return []byte{}, err
	}

	allRowStyle, err := exc.NewStyle(&excelize.Style{
		Border:    []excelize.Border{excelize.Border{Type: "center", Style: 1, Color: "#000000"}, excelize.Border{Type: "right", Style: 1, Color: "#000000"}, excelize.Border{Type: "left", Style: 1, Color: "#000000"}, excelize.Border{Type: "bottom", Style: 1, Color: "#000000"}, excelize.Border{Type: "top", Style: 1, Color: "#000000"}},
		Alignment: &excelize.Alignment{Horizontal: "center", WrapText: true},
	})
	if err != nil {
		return []byte{}, err
	}

	err = exc.SetCellStyle("Sheet1", "C4", "I"+strconv.Itoa(3+len(seminarDay.Presence)), allRowStyle)
	if err != nil {
		return []byte{}, err
	}

	exc.SetRowHeight("Sheet1", 1, 25.0)

	exc.SetColWidth("Sheet1", "B", "B", 7.0)
	exc.SetColWidth("Sheet1", "C", "C", 42.0)
	exc.SetColWidth("Sheet1", "D", "D", 18.0)
	exc.SetColWidth("Sheet1", "E", "E", 25.0)
	exc.SetColWidth("Sheet1", "F", "F", 55.0)
	exc.SetColWidth("Sheet1", "G", "G", 22.0)
	exc.SetColWidth("Sheet1", "H", "H", 15.0)
	exc.SetColWidth("Sheet1", "I", "I", 30.0)

	exc.SetCellValue("Sheet1", "B3", "Ред. бр.")
	exc.SetCellValue("Sheet1", "C3", "Име (име једног родитеља) презиме")
	exc.SetCellValue("Sheet1", "D3", "ЈМБГ")
	exc.SetCellValue("Sheet1", "E3", "Назив Центра за обуку")
	exc.SetCellValue("Sheet1", "F3", "Назив теме")
	exc.SetCellValue("Sheet1", "G3", "Датум семинара")
	exc.SetCellValue("Sheet1", "H3", "Датум уплате")
	exc.SetCellValue("Sheet1", "I3", "Правно лице/физичко лице")

	sort.Slice(seminarDay.Presence, func(i, j int) bool {
		return *seminarDay.Presence[i].Client.JMBG < *seminarDay.Presence[j].Client.JMBG
	})

	i := 0
	for _, t := range seminarDay.Presence {
		if t.Presence == nil || !*t.Presence {
			continue
		}
		exc.SetCellValue("Sheet1", "B"+strconv.Itoa(i+4), strconv.Itoa(i+1))
		exc.SetCellValue("Sheet1", "C"+strconv.Itoa(i+4), t.Client.Person.FullNameWithMiddleName())
		exc.SetCellValue("Sheet1", "D"+strconv.Itoa(i+4), *t.Client.JMBG)
		exc.SetCellValue("Sheet1", "E"+strconv.Itoa(i+4), util.CentarForEducationName)
		exc.SetCellValue("Sheet1", "F"+strconv.Itoa(i+4), seminarDay.Name)
		exc.SetCellValue("Sheet1", "G"+strconv.Itoa(i+4), seminarDay.Date.Format("02.01.2006"))
		clientSeminar, _ := SeminarService.GetClientSeminarBySeminarIDAndClientID(seminarDay.SeminarID, t.ClientID)
		if clientSeminar == nil {
			continue
		}
		payDate := ""
		if clientSeminar != nil && clientSeminar.PayDate != nil && !clientSeminar.PayDate.Before(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)) {
			payDate = clientSeminar.PayDate.Format("02.01.2006")
		}
		exc.SetCellValue("Sheet1", "H"+strconv.Itoa(i+4), payDate)
		payedBy := ""
		if clientSeminar.Payed != nil && *clientSeminar.Payed == true && clientSeminar.PayedBy != nil {
			payedBy = *clientSeminar.PayedBy
		}
		exc.SetCellValue("Sheet1", "I"+strconv.Itoa(i+4), payedBy)
		i++
	}

	var buf bytes.Buffer
	err = exc.Write(&buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (excelService) CreateSeminarsReportOfClients(seminars []model.Seminar) ([]byte, error) {
	exc := excelize.NewFile()
	firstRowStyle, err := exc.NewStyle(&excelize.Style{
		Border:    []excelize.Border{excelize.Border{Type: "center", Style: 1, Color: "#000000"}, excelize.Border{Type: "right", Style: 1, Color: "#000000"}, excelize.Border{Type: "left", Style: 1, Color: "#000000"}, excelize.Border{Type: "bottom", Style: 1, Color: "#000000"}, excelize.Border{Type: "top", Style: 1, Color: "#000000"}},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#D8D8D8"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	if err != nil {
		return []byte{}, err
	}

	err = exc.SetCellStyle("Sheet1", "B3", "H3", firstRowStyle)
	if err != nil {
		return []byte{}, err
	}

	allRowStyle, err := exc.NewStyle(&excelize.Style{
		Border:    []excelize.Border{excelize.Border{Type: "center", Style: 1, Color: "#000000"}, excelize.Border{Type: "right", Style: 1, Color: "#000000"}, excelize.Border{Type: "left", Style: 1, Color: "#000000"}, excelize.Border{Type: "bottom", Style: 1, Color: "#000000"}, excelize.Border{Type: "top", Style: 1, Color: "#000000"}},
		Alignment: &excelize.Alignment{Horizontal: "center", WrapText: true},
	})
	if err != nil {
		return []byte{}, err
	}

	err = exc.SetCellStyle("Sheet1", "C4", "I"+strconv.Itoa(3+len(seminars)), allRowStyle)
	if err != nil {
		return []byte{}, err
	}

	exc.SetRowHeight("Sheet1", 1, 25.0)

	exc.SetColWidth("Sheet1", "B", "B", 12.0)
	exc.SetColWidth("Sheet1", "C", "C", 20.0)
	exc.SetColWidth("Sheet1", "D", "D", 22.0)
	exc.SetColWidth("Sheet1", "E", "E", 15.0)
	exc.SetColWidth("Sheet1", "F", "F", 25.0)
	exc.SetColWidth("Sheet1", "G", "G", 15.0)
	exc.SetColWidth("Sheet1", "H", "H", 15.0)

	exc.SetCellValue("Sheet1", "B3", "Број сем.")
	exc.SetCellValue("Sheet1", "C3", "Шифра")
	exc.SetCellValue("Sheet1", "D3", "Локација")
	exc.SetCellValue("Sheet1", "E3", "Почетак")
	exc.SetCellValue("Sheet1", "F3", "Назив теме")
	exc.SetCellValue("Sheet1", "G3", "Статус")
	exc.SetCellValue("Sheet1", "H3", "Број полазника")

	sort.Slice(seminars, func(i, j int) bool {
		return seminars[j].ID < seminars[i].ID
	})

	totalTrainees := 0
	rowNum := 4
	for _, s := range seminars {
		if s.SeminarStatusID == model.SEMINAR_STATUS_OPENED || s.SeminarStatusID == model.SEMINAR_STATUS_FILLED {
			continue
		}
		exc.SetCellValue("Sheet1", "B"+strconv.Itoa(rowNum), strconv.Itoa(s.SerialNumberByLocation))
		exc.SetCellValue("Sheet1", "C"+strconv.Itoa(rowNum), s.GetCode())
		exc.SetCellValue("Sheet1", "D"+strconv.Itoa(rowNum), s.ClassRoom.Location.Address.Place)
		exc.SetCellValue("Sheet1", "E"+strconv.Itoa(rowNum), s.Start.Format("02.01.2006"))
		exc.SetCellValue("Sheet1", "F"+strconv.Itoa(rowNum), s.SeminarTheme.GetSeminarThemeWithBaseTheme())
		exc.SetCellValue("Sheet1", "G"+strconv.Itoa(rowNum), s.SeminarStatus.Name)
		exc.SetCellValue("Sheet1", "H"+strconv.Itoa(rowNum), strconv.Itoa(len(s.Trainees)))
		totalTrainees = totalTrainees + len(s.Trainees)
		rowNum++
	}

	exc.SetCellValue("Sheet1", "G"+strconv.Itoa(rowNum), "Укупно:")
	exc.SetCellValue("Sheet1", "H"+strconv.Itoa(rowNum), strconv.Itoa(totalTrainees))

	var buf bytes.Buffer
	err = exc.Write(&buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (excelService) CreateSeminarsReportOfTeachers(seminars []model.Seminar) ([]byte, error) {
	exc := excelize.NewFile()
	firstRowStyle, err := exc.NewStyle(&excelize.Style{
		Border:    []excelize.Border{excelize.Border{Type: "center", Style: 1, Color: "#000000"}, excelize.Border{Type: "right", Style: 1, Color: "#000000"}, excelize.Border{Type: "left", Style: 1, Color: "#000000"}, excelize.Border{Type: "bottom", Style: 1, Color: "#000000"}, excelize.Border{Type: "top", Style: 1, Color: "#000000"}},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#D8D8D8"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	if err != nil {
		return []byte{}, err
	}

	exc.SetRowHeight("Sheet1", 1, 25.0)

	exc.SetColWidth("Sheet1", "B", "u", 30.0)

	//seminar id [teacher] number of classes
	seminarTeachersMap := map[string]map[string]int{}
	columns := []string{"C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "R", "S", "T", "U"}
	columnCounter := 0
	teacherColumnMap := map[string]string{}

	for _, seminar := range seminars {
		if seminar.SeminarStatusID == model.SEMINAR_STATUS_OPENED || seminar.SeminarStatusID == model.SEMINAR_STATUS_FILLED {
			continue
		}

		teacherMap := map[string]int{}
		for _, day := range seminar.Days {
			for _, class := range day.Classes {
				if class.Teacher == nil {
					continue
				}
				if _, ok := teacherMap[class.Teacher.Person.FullName()]; ok {
					teacherMap[class.Teacher.Person.FullName()]++
				} else {
					teacherMap[class.Teacher.Person.FullName()] = 1
					if _, ok = teacherColumnMap[class.Teacher.Person.FullName()]; !ok {
						teacherColumnMap[class.Teacher.Person.FullName()] = columns[columnCounter]
						columnCounter++
					}
				}
			}
		}
		seminarTeachersMap[seminar.GetCode()] = teacherMap
	}

	for k, v := range teacherColumnMap {
		exc.SetCellValue("Sheet1", v+"3", k)
	}

	i := 4
	for s, m := range seminarTeachersMap {
		exc.SetCellValue("Sheet1", "B"+strconv.Itoa(i), s)
		for k, v := range m {
			exc.SetCellValue("Sheet1", teacherColumnMap[k]+strconv.Itoa(i), strconv.Itoa(v))
		}
		i++
	}

	err = exc.SetCellStyle("Sheet1", "B3", "B"+strconv.Itoa(i-1), firstRowStyle)
	if err != nil {
		return []byte{}, err
	}
	if columnCounter > 0 {
		err = exc.SetCellStyle("Sheet1", "B3", columns[columnCounter-1]+"3", firstRowStyle)
		if err != nil {
			return []byte{}, err
		}
	}

	var buf bytes.Buffer
	err = exc.Write(&buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (excelService) CreateClientsReport() ([]byte, error) {
	exc := excelize.NewFile()
	firstRowStyle, err := exc.NewStyle(&excelize.Style{
		Border:    []excelize.Border{excelize.Border{Type: "center", Style: 1, Color: "#000000"}, excelize.Border{Type: "right", Style: 1, Color: "#000000"}, excelize.Border{Type: "left", Style: 1, Color: "#000000"}, excelize.Border{Type: "bottom", Style: 1, Color: "#000000"}, excelize.Border{Type: "top", Style: 1, Color: "#000000"}},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#D8D8D8"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	if err != nil {
		return []byte{}, err
	}

	sheetRowMap := map[uint]int{}

	locationSheetMap := map[uint]string{}
	locationSheetMap[0] = "Sheet1"
	sheetRowMap[0] = 2

	locations, err := LocationService.GetAllLocations()
	if err != nil {
		return []byte{}, err
	}

	for _, l := range locations {
		sheet := l.Code + "-" + l.Address.Place
		locationSheetMap[l.ID] = sheet
		sheetRowMap[l.ID] = 2
		exc.NewSheet(sheet)
	}

	for _, sheet := range locationSheetMap {
		err = exc.SetCellStyle(sheet, "B1", "J1", firstRowStyle)
		if err != nil {
			return []byte{}, err
		}
		exc.SetRowHeight(sheet, 1, 25.0)

		exc.SetColWidth(sheet, "B", "B", 20.0)
		exc.SetColWidth(sheet, "C", "C", 23.0)
		exc.SetColWidth(sheet, "D", "D", 12.0)
		exc.SetColWidth(sheet, "E", "E", 15.0)
		exc.SetColWidth(sheet, "F", "F", 15.0)
		exc.SetColWidth(sheet, "G", "G", 15.0)
		exc.SetColWidth(sheet, "H", "H", 15.0)
		exc.SetColWidth(sheet, "I", "I", 15.0)

		exc.SetCellValue(sheet, "B1", "Ime i prezime")
		exc.SetCellValue(sheet, "C1", "Firma")
		exc.SetCellValue(sheet, "D1", "Broj telefona")
		exc.SetCellValue(sheet, "E1", "Radno vreme")
		exc.SetCellValue(sheet, "F1", "Dokumenta")
		exc.SetCellValue(sheet, "G1", "Teret")
		exc.SetCellValue(sheet, "H1", "Propisi")
		exc.SetCellValue(sheet, "I1", "Tahografi 2")
	}

	clientCount, err := ClientService.GetClientsCount()
	if err != nil {
		return []byte{}, err
	}

	var batch int64 = 200
	pages := clientCount/batch + 1

	for i := 0; i < int(pages); i++ {
		clients, err := ClientService.GetAllClientsWithSeminarsAndBasePersonalInfo(i*int(batch), int(batch))
		if err != nil {
			return []byte{}, err
		}

		for _, c := range clients {
			var location uint
			if len(c.Seminars) > 0 {
				location = c.Seminars[0].Seminar.ClassRoom.LocationID
			}
			exc.SetCellValue(locationSheetMap[location], "B"+strconv.Itoa(sheetRowMap[location]), c.Person.FullName())
			company := c.Company.Name
			if len(company) > 23 {
				company = company[0:20] + ".."
			}
			exc.SetCellValue(locationSheetMap[location], "C"+strconv.Itoa(sheetRowMap[location]), company)
			exc.SetCellValue(locationSheetMap[location], "D"+strconv.Itoa(sheetRowMap[location]), c.Person.PhoneNumber)

			if c.PassedCheckboxes.WorkTimeAndTahografs != nil && *c.PassedCheckboxes.WorkTimeAndTahografs {
				exc.SetCellValue(locationSheetMap[location], "E"+strconv.Itoa(sheetRowMap[location]), "+")
			}
			if c.PassedCheckboxes.ThemeDocuments != nil && *c.PassedCheckboxes.ThemeDocuments {
				exc.SetCellValue(locationSheetMap[location], "F"+strconv.Itoa(sheetRowMap[location]), "+")
			}
			if c.PassedCheckboxes.Burden != nil && *c.PassedCheckboxes.Burden {
				exc.SetCellValue(locationSheetMap[location], "G"+strconv.Itoa(sheetRowMap[location]), "+")
			}
			if c.PassedCheckboxes.Regulations != nil && *c.PassedCheckboxes.Regulations {
				exc.SetCellValue(locationSheetMap[location], "H"+strconv.Itoa(sheetRowMap[location]), "+")
			}
			if c.PassedCheckboxes.Tahografs2 != nil && *c.PassedCheckboxes.Tahografs2 {
				exc.SetCellValue(locationSheetMap[location], "I"+strconv.Itoa(sheetRowMap[location]), "+")
			}

			for _, s := range c.Seminars {
				switch s.Seminar.SeminarTheme.Code {
				case "1":
					exc.SetCellValue(locationSheetMap[location], "E"+strconv.Itoa(sheetRowMap[location]), s.Seminar.Start.Format("02.01.2006."))
				case "2":
					exc.SetCellValue(locationSheetMap[location], "F"+strconv.Itoa(sheetRowMap[location]), s.Seminar.Start.Format("02.01.2006."))
				case "3":
					exc.SetCellValue(locationSheetMap[location], "G"+strconv.Itoa(sheetRowMap[location]), s.Seminar.Start.Format("02.01.2006."))
				case "4":
					exc.SetCellValue(locationSheetMap[location], "H"+strconv.Itoa(sheetRowMap[location]), s.Seminar.Start.Format("02.01.2006."))
				case "5":
					exc.SetCellValue(locationSheetMap[location], "I"+strconv.Itoa(sheetRowMap[location]), s.Seminar.Start.Format("02.01.2006."))
				}
			}

			sheetRowMap[location] = sheetRowMap[location] + 1
		}
	}

	var buf bytes.Buffer
	err = exc.Write(&buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}
