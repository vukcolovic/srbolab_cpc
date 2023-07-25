package service

import (
	"bytes"
	"github.com/go-pdf/fpdf"
	"os"
	"path/filepath"
	"srbolab_cpc/logoped"
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
	PrintConfirmationStatements(seminar *model.Seminar) ([]byte, error)
}

func (p *printService) PrintSeminarStudentList(seminar *model.Seminar) ([]byte, error) {
	pwd, err := os.Getwd()
	if err != nil {
		logoped.ErrorLog.Println("Error getting pwd: ", err)
		return []byte{}, err
	}
	pdf := fpdf.New("L", "mm", "A4", filepath.Join(pwd, "font"))
	pdf.AddFont("Arimo-Regular", "", "Arimo-Regular.json")
	pdf.AddFont("Arimo-Bold", "", "Arimo-Bold.json")
	latTr := pdf.UnicodeTranslatorFromDescriptor("iso-8859-16")

	pdf.SetMargins(marginLeft, marginTop, marginRight)
	pdf.AddPage()

	pdf.SetFont("Arimo-Regular", "", 8)
	createSimpleHeader(pdf, latTr)

	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), "Broj dokumenta: ")
	pdf.Text(35, pdf.GetY(), "")
	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), latTr("Šifra obuke: "))
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
	err = pdf.Output(&buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (p *printService) PrintConfirmationStatements(seminar *model.Seminar) ([]byte, error) {
	pwd, err := os.Getwd()
	if err != nil {
		logoped.ErrorLog.Println("Error getting pwd: ", err)
		return []byte{}, err
	}
	pdf := fpdf.New("P", "mm", "A4", filepath.Join(pwd, "font"))
	pdf.AddFont("Arimo-Regular", "", "Arimo-Regular.json")
	pdf.AddFont("Arimo-Bold", "", "Arimo-Bold.json")
	latTr := pdf.UnicodeTranslatorFromDescriptor("iso-8859-16")

	pdf.SetMargins(marginLeft, 40, marginRight)

	for _, client := range seminar.Trainees {
		pdf.AddPage()

		pdf.SetFont("Arimo-Regular", "", 11)

		pdf.Ln(5)
		pdf.Text(15, pdf.GetY(), latTr("U skladu sa članom 6. i članom 7. uredbe EU 2016/679 od 27. aprila 2016. godine i članom 12. i članom"))
		pdf.Ln(5)
		pdf.Text(15, pdf.GetY(), latTr("15. Zakona o zaštiti podataka o ličnosti (Sl. Glasnik RS“, br. 87/2018 od 13/11/2018) dajem pristanak za"))
		pdf.Ln(5)
		pdf.Text(15, pdf.GetY(), "davanje i obradu podataka o ličnosti, gde je rukovalac obrade Srbolab.")
		pdf.Ln(20)

		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(80, pdf.GetY(), "IZJAVA O PRISTANKU")
		pdf.Ln(20)

		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(15, pdf.GetY(), "Ja")
		pdf.Line(20, pdf.GetY(), 80, pdf.GetY())
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(25, pdf.GetY()-1, latTr(client.Client.Person.FullName()))
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(82, pdf.GetY(), "(ime i prezime),")
		pdf.Line(110, pdf.GetY(), 170, pdf.GetY())
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(115, pdf.GetY()-1, latTr(*client.Client.JMBG))
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(172, pdf.GetY(), "(JMBG),")

		pdf.Ln(5)
		pdf.Text(15, pdf.GetY(), "pristajem na davanje i obradu sledećih svojih podataka o ličnosti: podaci iz lične karte/pasoša, podaci iz")
		pdf.Ln(5)
		pdf.Text(15, pdf.GetY(), "nacionalne vozačke dozvole, podaci iz kvalifikacione kartice vozača, elektronsku adresu, kontakt telefon")
		pdf.Ln(5)
		pdf.Text(15, pdf.GetY(), "vozača, podaci o stručnoj spremi, za potrebe slanja obaveštenja i informacija.")
		pdf.Ln(15)

		pdf.Text(15, pdf.GetY(), "Takođe izjavljujem da sam od AMSS-CMV primio/la sva neophodna obaveštenja, predviđena članom 23")
		pdf.Ln(5)
		pdf.Text(15, pdf.GetY(), "Zakona o zaštiti podataka o ličnosti, kao i obaveštenje da u svakom trenutku mogu opozvati dat")
		pdf.Ln(5)
		pdf.Text(15, pdf.GetY(), "pristanak, s tim da opoziv pristanka ne utiče na dopuštenost obrade koja je vršena na osnovu pristanka")
		pdf.Ln(5)
		pdf.Text(15, pdf.GetY(), "pre opoziva, kao i da nisam u obavezi da dam podatke o ličnosti koji nisu predviđeni kao obavezni")
		pdf.Ln(5)
		pdf.Text(15, pdf.GetY(), "zakonskim i podzakonskim aktima i da isto neće biti od uticaja na pružanje usluga od strane rukovaoca.")
		pdf.Ln(25)

		pdf.Text(15, pdf.GetY(), "Datum: ")
		pdf.Line(30, pdf.GetY(), 70, pdf.GetY())
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(35, pdf.GetY()-1, time.Now().Format("02.01.2006"))
		pdf.SetFont("Arimo-Regular", "", 11)

		pdf.Ln(15)

		pdf.Text(15, pdf.GetY(), "Potpis: ")
		pdf.Line(30, pdf.GetY(), 70, pdf.GetY())
	}

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func createSimpleHeader(pdf *fpdf.Fpdf, tr func(string) string) {
	pdf.SetFont("Arimo-Regular", "", 10)
	pdf.Image("./images/srbolab_logo.png", 15, 10, 30, 10, false, "png", 0, "")
	pdf.CellFormat(35, 10, "", "0", 0, "C", false, 0, "")
	pdf.Text(100, 14, "SRBOLAB D.O.O.")
	pdf.SetFont("Arial", "", 10)
	pdf.Text(60, 18, tr("SEKTOR ZA STRUČNO USAVRŠAVANJE, RAZVOJ I BEZBEDNOST SAOBRAĆAJA"))
	pdf.Ln(15)
}
