package service

import (
	"bytes"
	"fmt"
	"github.com/xuri/excelize/v2"
	"srbolab_cpc/model"
	"srbolab_cpc/util"
	"strconv"
	"time"
)

var (
	ExcelService excelServiceInterface = &excelService{}
)

type excelService struct {
}

type excelServiceInterface interface {
	CreateClientTestsBySeminarDayReport(tests []model.ClientTest) ([]byte, error)
	CreateListClientsBySeminarDayReport(seminarDay *model.SeminarDay) ([]byte, error)
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

	for i, t := range seminarDay.Presence {
		exc.SetCellValue("Sheet1", "B"+strconv.Itoa(i+4), strconv.Itoa(i+1))
		exc.SetCellValue("Sheet1", "C"+strconv.Itoa(i+4), t.Client.Person.FullName())
		exc.SetCellValue("Sheet1", "D"+strconv.Itoa(i+4), *t.Client.JMBG)
		exc.SetCellValue("Sheet1", "E"+strconv.Itoa(i+4), util.CentarForEducationName)
		exc.SetCellValue("Sheet1", "F"+strconv.Itoa(i+4), seminarDay.Name)
		exc.SetCellValue("Sheet1", "G"+strconv.Itoa(i+4), t.CreatedAt.Format("02.01.2006"))
		clientSeminar, _ := SeminarService.GetClientSeminarBySeminarIDAndClientID(seminarDay.SeminarID, t.ClientID)
		payDate := ""
		if clientSeminar != nil && clientSeminar.PayDate != nil && !clientSeminar.PayDate.Before(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)) {
			payDate = clientSeminar.PayDate.Format("02.01.2006")
		}
		exc.SetCellValue("Sheet1", "H"+strconv.Itoa(i+4), payDate)
		company := ""
		if t.Client.Company.ID > 0 {
			company = t.Client.Company.Name
		}
		exc.SetCellValue("Sheet1", "I"+strconv.Itoa(i+4), company)
	}

	var buf bytes.Buffer
	err = exc.Write(&buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}
