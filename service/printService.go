package service

import (
	"bytes"
	"github.com/go-pdf/fpdf"
	"srbolab_cpc/model"
	"strconv"
	"time"
)

var (
	PrintService printServiceInterface = &printService{}
)

const (
	marginLeft  = 10.0
	marginTop   = 15.0
	marginRight = 10.0
)

type printService struct {
}

type printServiceInterface interface {
	PrintSeminarStudentList(seminar *model.Seminar) ([]byte, error)
}

func (p *printService) PrintSeminarStudentList(seminar *model.Seminar) ([]byte, error) {
	pdf := fpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 8)
	pdf.SetMargins(marginLeft, marginTop, marginRight)

	createSimpleHeader(pdf)

	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), "Broj dokumenta: ")
	pdf.Text(35, pdf.GetY(), "")
	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), "Šifra obuke: ")
	pdf.Text(30, pdf.GetY(), "")
	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), "Mesto: ")
	pdf.Text(27, pdf.GetY(), seminar.ClassRoom.Location.Address.Place)
	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), "Registracioni list - spisak polaznika za "+seminar.SeminarTheme.Name)
	pdf.Text(140, pdf.GetY(), "Datum: "+time.Now().Format("01.02.2006"))

	ch := 8.0
	pdf.Ln(ch)
	pdf.CellFormat(20, ch, "Redni broj", "1", 0, "C", false, 0, "")
	pdf.CellFormat(75, ch, "Ime i prezime", "1", 0, "C", false, 0, "")
	pdf.CellFormat(35, ch, "JMBG", "1", 0, "C", false, 0, "")
	pdf.CellFormat(90, ch, "Firma u kojoj ste zaposleni/Telefon", "1", 0, "C", false, 0, "")
	pdf.CellFormat(60, ch, "Potpis", "1", 0, "C", false, 0, "")

	for i, cs := range seminar.Trainees {
		pdf.Ln(ch)
		pdf.CellFormat(20, ch, strconv.Itoa(i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(75, ch, cs.Client.Person.FullName(), "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, ch, *cs.Client.JMBG, "1", 0, "C", false, 0, "")
		pdf.CellFormat(90, ch, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(60, ch, "", "1", 0, "C", false, 0, "")
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func createSimpleHeader(pdf *fpdf.Fpdf) {
	pdf.SetFont("Arial", "", 10)
	pdf.Image("./images/srbolab_logo.png", 15, 10, 30, 10, false, "png", 0, "")
	pdf.CellFormat(35, 10, "", "0", 0, "C", false, 0, "")
	pdf.Text(100, 14, "SRBOLAB D.O.O.")
	pdf.SetFont("Arial", "", 10)
	pdf.Text(60, 18, "SEKTOR ZA STRUČNO USAVRŠAVANJE, RAZVOJ I BEZBEDNOST SAOBRAĆAJA")
	pdf.Ln(15)
}
