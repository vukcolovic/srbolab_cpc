package service

import (
	"bytes"
	"github.com/go-pdf/fpdf"
	"os"
	"path/filepath"
	"srbolab_cpc/logoped"
	"srbolab_cpc/model"
	"srbolab_cpc/util"
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
	PrintConfirmations(seminar *model.Seminar) ([]byte, error)
	PrintConfirmationReceives(seminar *model.Seminar) ([]byte, error)
	PrintMuster(day *model.SeminarDay) ([]byte, error)
}

func (p *printService) PrintSeminarStudentList(seminar *model.Seminar) ([]byte, error) {
	pwd, err := os.Getwd()
	if err != nil {
		logoped.ErrorLog.Println("Error getting pwd: ", err)
		return []byte{}, err
	}
	pdf := fpdf.New("L", "mm", "A4", filepath.Join(pwd, "font"))
	pdf.AddFont("Arimo-Regular", "", "Arimo-Regular.json")
	//pdf.AddFont("Arimo-Bold", "", "Arimo-Bold.json")
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
		pdf.Text(15, pdf.GetY(), latTr("davanje i obradu podataka o ličnosti, gde je rukovalac obrade Srbolab."))
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
		pdf.Text(15, pdf.GetY(), latTr("pristajem na davanje i obradu sledećih svojih podataka o ličnosti: podaci iz lične karte/pasoša, podaci iz"))
		pdf.Ln(5)
		pdf.Text(15, pdf.GetY(), latTr("nacionalne vozačke dozvole, podaci iz kvalifikacione kartice vozača, elektronsku adresu, kontakt telefon"))
		pdf.Ln(5)
		pdf.Text(15, pdf.GetY(), latTr("vozača, podaci o stručnoj spremi, za potrebe slanja obaveštenja i informacija."))
		pdf.Ln(15)

		pdf.Text(15, pdf.GetY(), latTr("Takođe izjavljujem da sam od AMSS-CMV primio/la sva neophodna obaveštenja, predviđena članom 23"))
		pdf.Ln(5)
		pdf.Text(15, pdf.GetY(), latTr("Zakona o zaštiti podataka o ličnosti, kao i obaveštenje da u svakom trenutku mogu opozvati dat"))
		pdf.Ln(5)
		pdf.Text(15, pdf.GetY(), latTr("pristanak, s tim da opoziv pristanka ne utiče na dopuštenost obrade koja je vršena na osnovu pristanka"))
		pdf.Ln(5)
		pdf.Text(15, pdf.GetY(), latTr("pre opoziva, kao i da nisam u obavezi da dam podatke o ličnosti koji nisu predviđeni kao obavezni"))
		pdf.Ln(5)
		pdf.Text(15, pdf.GetY(), latTr("zakonskim i podzakonskim aktima i da isto neće biti od uticaja na pružanje usluga od strane rukovaoca."))
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

func (p *printService) PrintConfirmations(seminar *model.Seminar) ([]byte, error) {
	pwd, err := os.Getwd()
	if err != nil {
		logoped.ErrorLog.Println("Error getting pwd: ", err)
		return []byte{}, err
	}
	pdf := fpdf.New("P", "mm", "A4", filepath.Join(pwd, "font"))
	pdf.AddFont("Arimo-Regular", "", "Arimo-Regular.json")
	pdf.AddFont("Arimo-Bold", "", "Arimo-Bold.json")
	latTr := pdf.UnicodeTranslatorFromDescriptor("iso-8859-16")

	pdf.SetMargins(15, 40, marginRight)

	startSeminar := seminar.Start
	endSeminar := seminar.Start
	for _, day := range seminar.Days {
		if day.Date.After(endSeminar) {
			endSeminar = day.Date
		}
	}

	for _, client := range seminar.Trainees {
		pdf.AddPage()

		pdf.SetFont("Arimo-Regular", "", 11)

		pdf.Ln(5)
		pdf.Text(15, pdf.GetY(), "Broj:")
		pdf.Ln(5)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(15, pdf.GetY(), "Dana:")
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(30, pdf.GetY(), time.Now().Format("02.01.2006."))
		pdf.Ln(5)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(15, pdf.GetY(), "Mesto:")
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(30, pdf.GetY(), latTr(seminar.ClassRoom.Location.Address.Place))
		pdf.Ln(15)

		pdf.Text(80, pdf.GetY(), "POTVRDA")
		pdf.Ln(10)

		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(25, pdf.GetY(), latTr("o završenoj periodičnoj obuci na obaveynim seminarima unapređenja znanja"))
		pdf.Ln(10)

		ch := 9.0
		wl := 80.0
		wr := 100.0

		pdf.CellFormat(wl, ch-1, "Ime, ime jednog", "LRT", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.CellFormat(wr, ch-1, latTr(client.Client.Person.FirstName+" ("+client.Client.Person.MiddleName+") "+client.Client.Person.LastName), "LRT", 0, "L", false, 0, "")
		pdf.Ln(ch - 1)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.CellFormat(wl, ch-1, "roditelja, prezime", "LRB", 0, "L", false, 0, "")
		pdf.CellFormat(wr, ch-1, "", "LRB", 0, "L", false, 0, "")
		pdf.Ln(ch - 1)
		pdf.CellFormat(wl, ch, "JMBG", "1", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", 11)
		for _, l := range *client.Client.JMBG {
			pdf.CellFormat(wr/13, ch, string(l), "1", 0, "C", false, 0, "")
		}
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.CellFormat(wl, ch, latTr("Mesto prebivališta"), "1", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.CellFormat(wr, ch, latTr(client.Client.Address.Place), "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.CellFormat(wl, ch, latTr("Adresa prebivališta"), "1", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.CellFormat(wr, ch, latTr(client.Client.Address.Street+" "+client.Client.Address.HouseNumber), "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.CellFormat(wl, ch, "Redni broj seminara", "1", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", 11)
		//fixme just closed seminars and passed
		completedSeminars := 0
		if client.Client.InitialCompletedSeminars != nil {
			completedSeminars = *client.Client.InitialCompletedSeminars
		}
		seminarNumber := completedSeminars + len(client.Client.Seminars)
		cx := 87 + float64(seminarNumber*8)
		pdf.Circle(cx, 137, 3, "")
		pdf.CellFormat(wr, ch, " I    II    III    IV    V", "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.CellFormat(wl, ch-1, latTr("Datum pohađanja"), "LRT", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.CellFormat(wr, ch-1, "od "+startSeminar.Format("02.01.2006"), "LRT", 0, "L", false, 0, "")
		pdf.Ln(ch - 1)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.CellFormat(wl, ch-1, latTr("periodične obuke"), "LRB", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.CellFormat(wr, ch-1, "do "+endSeminar.Format("02.01.2006"), "LRB", 0, "L", false, 0, "")
		pdf.Ln(ch - 1)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.CellFormat(wl, ch-1, latTr("Mesto pohađanja"), "LRT", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.CellFormat(wr, ch-1, latTr(client.Client.Address.Place+", "+client.Client.Address.Street+" "+client.Client.Address.HouseNumber), "LRT", 0, "L", false, 0, "")
		pdf.Ln(ch - 1)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.CellFormat(wl, ch-1, latTr("periodične obuke"), "LRB", 0, "L", false, 0, "")
		pdf.CellFormat(wr, ch-1, "", "LRB", 0, "L", false, 0, "")
		pdf.Ln(ch - 1)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.CellFormat(wl, ch, "Vrsta CPC", "1", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Circle(97.5, 178.5, 2.5, "")
		pdf.CellFormat(wr, ch, "1. prevoz tereta  2. prevoz putnika", "1", 0, "L", false, 0, "")
		pdf.Ln(20)

		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(15, pdf.GetY(), "NAPOMENA:")
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(50, pdf.GetY(), latTr("Ova potvrda se izdaje na osnovu odslušane obavezne periodične obuke"))
		pdf.Ln(5)
		pdf.Text(15, pdf.GetY(), latTr("za potrebe sticanja periodičnog CPC i ne može se koristiti u druge svrhe."))
	}

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (p *printService) PrintConfirmationReceives(seminar *model.Seminar) ([]byte, error) {
	pwd, err := os.Getwd()
	if err != nil {
		logoped.ErrorLog.Println("Error getting pwd: ", err)
		return []byte{}, err
	}
	pdf := fpdf.New("P", "mm", "A4", filepath.Join(pwd, "font"))
	pdf.AddFont("Arimo-Regular", "", "Arimo-Regular.json")
	pdf.AddFont("Arimo-Bold", "", "Arimo-Bold.json")
	latTr := pdf.UnicodeTranslatorFromDescriptor("iso-8859-16")

	pdf.SetMargins(marginLeft, 20, marginRight)

	for _, client := range seminar.Trainees {
		pdf.AddPage()

		createSimpleHeader(pdf, latTr)

		pdf.SetFont("Arimo-Bold", "", 12)
		pdf.Text(40, pdf.GetY(), latTr("Izjava o preuzimanju potvrde i završenoj periodičnoj obuci"))
		pdf.Ln(5)
		pdf.Text(60, pdf.GetY(), latTr("na obaveznim seminarima unapređenja znanja"))
		pdf.Ln(20)

		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(15, pdf.GetY(), "Dana")
		pdf.Line(27, pdf.GetY(), 57, pdf.GetY())
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(30, pdf.GetY()-1, latTr(time.Now().Format("02.01.2006.")))
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(60, pdf.GetY(), "godine, ")
		pdf.Line(75, pdf.GetY(), 135, pdf.GetY())
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(80, pdf.GetY()-1, latTr(client.Client.Person.FullName()))
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(135, pdf.GetY(), "(ime i prezime), JMBG")
		pdf.Ln(6)
		pdf.Line(15, pdf.GetY(), 50, pdf.GetY())
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(20, pdf.GetY()-1, *client.Client.JMBG)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(55, pdf.GetY(), latTr("je preuzeo potvrdu o završenoj periodičnoj obuci na"))
		pdf.Ln(6)
		pdf.Text(15, pdf.GetY(), latTr("obaveznim seminarima unapređenja znanja."))

		pdf.Ln(80)

		pdf.Text(15, pdf.GetY(), "Potvrdu preuzeo: ")
		pdf.Ln(10)
		pdf.Line(15, pdf.GetY(), 60, pdf.GetY())
		pdf.Ln(8)
		pdf.Text(15, pdf.GetY(), "Dana: ")
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Line(26, pdf.GetY(), 48, pdf.GetY())
		pdf.Text(27, pdf.GetY()-1, time.Now().Format("02.01.2006"))
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(50, pdf.GetY(), "godine.")
	}

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (p *printService) PrintMuster(day *model.SeminarDay) ([]byte, error) {
	pwd, err := os.Getwd()
	if err != nil {
		logoped.ErrorLog.Println("Error getting pwd: ", err)
		return []byte{}, err
	}
	pdf := fpdf.New("L", "mm", "A4", filepath.Join(pwd, "font"))
	pdf.AddFont("Arimo-Regular", "", "Arimo-Regular.json")
	//pdf.AddFont("Arimo-Bold", "", "Arimo-Bold.json")
	latTr := pdf.UnicodeTranslatorFromDescriptor("iso-8859-16")

	pdf.SetMargins(marginLeft, marginTop, marginRight)
	pdf.AddPage()

	pdf.SetFont("Arimo-Regular", "", 8)
	createSimpleHeader(pdf, latTr)

	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), "Mesto: ")
	pdf.Text(27, pdf.GetY(), day.Seminar.ClassRoom.Location.Address.Place)
	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), latTr("Šifra obuke: "))
	pdf.Text(30, pdf.GetY(), "")
	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), "Datum: ")
	dayInWeek := util.GetDaySerbian(day.Date)
	pdf.Text(27, pdf.GetY(), dayInWeek+" "+day.Date.Format("02.01.2006."))
	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), latTr("Prozivnik polaznika seminara unapređenja znanja na periodičnoj obuci profesionalnih vozača"))

	ch := 14.0
	pdf.Ln(2)
	pdf.CellFormat(10, ch, "R.B.", "1", 0, "C", false, 0, "")
	pdf.CellFormat(45, ch, "Ime i prezime / JMBG", "1", 0, "C", false, 0, "")
	pdf.CellFormat(27, ch, latTr("1. čas"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(27, ch, latTr("2. čas"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(27, ch, latTr("3. čas"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(27, ch, latTr("4. čas"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(27, ch, latTr("5. čas"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(27, ch, latTr("6. čas"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(27, ch, latTr("7. čas"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(31, ch, "Napomena", "1", 0, "C", false, 0, "")

	for i, cs := range day.Presence {
		pdf.Ln(ch)
		pdf.CellFormat(10, ch, strconv.Itoa(i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(45, ch, "", "1", 0, "C", false, 0, "")
		pdf.Text(25, 71+float64(i*14), cs.Client.Person.FullName())
		pdf.Text(25, 76+float64(i*14), *cs.Client.JMBG)
		pdf.CellFormat(27, ch, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(27, ch, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(27, ch, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(27, ch, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(27, ch, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(27, ch, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(27, ch, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(31, ch, "", "1", 0, "C", false, 0, "")
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
	pdf.Text(60, 18, tr("SEKTOR ZA STRUČNO USAVRŠAVANJE, RAZVOJ I BEZBEDNOST SAOBRAĆAJA"))
	pdf.Ln(15)
}
