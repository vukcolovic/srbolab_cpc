package service

import (
	"bytes"
	"fmt"
	"github.com/xuri/excelize/v2"
	"srbolab_cpc/model"
	"strconv"
)

var (
	ExcelService excelServiceInterface = &excelService{}
)

type excelService struct {
}

type excelServiceInterface interface {
	CreateClientTestsBySeminarDayReport(tests []model.ClientTest) ([]byte, error)
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

	exc.SetCellValue("Sheet1", "A1", "Ime i prezime")
	exc.SetCellValue("Sheet1", "B1", "Seminar (tema)")
	exc.SetCellValue("Sheet1", "C1", "Dan")
	exc.SetCellValue("Sheet1", "D1", "Vreme testa")
	exc.SetCellValue("Sheet1", "E1", "Rezultat (%)")
	exc.SetCellValue("Sheet1", "F1", "Odgovori")

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
