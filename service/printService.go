package service

import (
	"bytes"
	"fmt"
	"github.com/go-pdf/fpdf"
	"github.com/skip2/go-qrcode"
	"os"
	"path/filepath"
	"sort"
	"srbolab_cpc/logoped"
	"srbolab_cpc/model"
	"srbolab_cpc/util"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var (
	PrintService printServiceInterface = &printService{}
	Domain       string
)

const (
	marginLeft  = 10.0
	marginTop   = 15.0
	marginRight = 10.0
)

type translationDetails struct {
	pdf             *fpdf.Fpdf
	cirilicFont     string
	latinicFont     string
	defaultFontSize float64
	latinicFunc     func(string) string
	cirilicFunc     func(string) string
}

func newTranslationDetails(pdf *fpdf.Fpdf, cirilicFont, latinicFont string, defaultFontSize float64, latinicFunc, cirilicFunc func(string) string) *translationDetails {
	return &translationDetails{
		pdf:             pdf,
		cirilicFont:     cirilicFont,
		latinicFont:     latinicFont,
		defaultFontSize: defaultFontSize,
		latinicFunc:     latinicFunc,
		cirilicFunc:     cirilicFunc,
	}
}

// method for automatic choose translate on latinic or cirilic
// font size is nessesary because of custom solution of changing font family on changing translations
// if font size is same as default set number less then 1
func (t *translationDetails) translate(input string, fontSize float64) string {
	if fontSize < 1 {
		fontSize = t.defaultFontSize
	}
	for _, r := range input {
		if unicode.Is(unicode.Cyrillic, r) {
			t.pdf.SetFont(t.cirilicFont, "", fontSize)
			return t.cirilicFunc(input)
		}
	}
	t.pdf.SetFont(t.latinicFont, "", fontSize)
	return t.latinicFunc(input)
}

func (t *translationDetails) translDef(input string) string {
	return t.translate(input, -1)
}

type printService struct {
}

type printServiceInterface interface {
	PrintSeminarStudentList(seminar *model.Seminar) ([]byte, error)
	PrintConfirmations(seminar *model.Seminar) ([]byte, error)
	PrintConfirmationReceives(seminar *model.Seminar) ([]byte, error)
	PrintMuster(day *model.SeminarDay) ([]byte, error)
	PrintCheckIn(seminar *model.Seminar) ([]byte, error)
	PrintSeminarEvidence(day *model.SeminarDay) ([]byte, error)
	PrintTestBarcode() ([]byte, error)
	PrintPlanTreningRealization(day *model.SeminarDay) ([]byte, error)
	PrintPayments(seminar *model.Seminar) ([]byte, error)
	PrintSeminarReport(seminar *model.Seminar) ([]byte, error)
	PrintSeminarReport2(seminar *model.Seminar) ([]byte, error)
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
	pdf.AddFont("Helvetica", "", "helvetica_1251.json")
	latTr := pdf.UnicodeTranslatorFromDescriptor("iso-8859-16")
	cirTr := pdf.UnicodeTranslatorFromDescriptor("cp1251")
	fontSize := 10.0
	trObj := newTranslationDetails(pdf, "Helvetica", "Arimo-Regular", fontSize, latTr, cirTr)

	pdf.SetMargins(marginLeft, marginTop, marginRight)
	pdf.AddPage()

	pdf.SetFont("Helvetica", "", fontSize)
	createSimpleHeaderForLandscape(pdf, cirTr)

	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), trObj.translDef("Број документа: "))
	pdf.Text(37, pdf.GetY(), "")
	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), trObj.translDef(("Шифра обуке: ")))
	pdf.Text(39, pdf.GetY(), trObj.translDef(seminar.GetCode()))
	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), trObj.translDef("Место: "))
	pdf.Text(28, pdf.GetY(), trObj.translDef(seminar.ClassRoom.Location.Address.Place))
	pdf.Ln(5)
	seminarType := "периодичну"
	if seminar.SeminarTheme.BaseSeminarType.Code == "ADDITIONAL" {
		seminarType = "додатну"
	}
	if seminar.SeminarTheme.BaseSeminarType.Code == "BASE" {
		seminarType = "основну"
	}
	pdf.Text(15, pdf.GetY(), trObj.translDef(fmt.Sprintf("Регистрациони лист - списак полазника за %s обуку", seminarType)))
	pdf.Text(140, pdf.GetY(), trObj.translDef("Датум: "+seminar.Start.Format("02.01.2006")))

	ch := 8.0
	pdf.Ln(ch)
	pdf.CellFormat(20, ch, trObj.translDef("Редни број"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(70, ch, trObj.translDef("Име и презиме"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, ch, trObj.translDef("ЈМБГ"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(80, ch, trObj.translDef("Фирма у којој сте запослени"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(35, ch, trObj.translDef("Телефон"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(45, ch, trObj.translDef("Потпис"), "1", 0, "C", false, 0, "")

	sort.Slice(seminar.Trainees, func(i, j int) bool {
		return *seminar.Trainees[i].Client.JMBG < *seminar.Trainees[j].Client.JMBG
	})

	pdf.Ln(ch)
	for i, cs := range seminar.Trainees {
		lines, num := splitLine(cs.Client.Company.Name, 40)

		chc := ch
		if num > 1 {
			chc = chc * num * 0.7
		}
		pdf.CellFormat(20, chc, strconv.Itoa(i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(70, chc, trObj.translDef(cs.Client.Person.FullName()), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, chc, *cs.Client.JMBG, "1", 0, "C", false, 0, "")
		pdf.CellFormat(80, chc, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, chc, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(45, chc, "", "1", 0, "C", false, 0, "")
		current := pdf.GetY() + 4.5
		for i, line := range lines {
			pdf.Text(130, current+float64(i)*4.0, trObj.translDef(line))
		}
		pdf.Ln(chc)
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
	pdf.AddFont("Helvetica", "", "helvetica_1251.json")
	latTr := pdf.UnicodeTranslatorFromDescriptor("iso-8859-16")
	cirTr := pdf.UnicodeTranslatorFromDescriptor("cp1251")
	trObj := newTranslationDetails(pdf, "Helvetica", "Arimo-Regular", 11, latTr, cirTr)

	pdf.SetMargins(15, 40, marginRight)

	startSeminar := seminar.Start
	endSeminar := seminar.Start
	notPassedClientIds := make(map[uint]string)
	for _, day := range seminar.Days {
		if day.Date.After(endSeminar) {
			endSeminar = day.Date
		}

		for _, p := range day.Presence {
			if !*p.Presence && !day.Date.After(time.Now()) {
				notPassedClientIds[p.ClientID] = ""
			}
		}
	}

	sort.Slice(seminar.Trainees, func(i, j int) bool {
		return *seminar.Trainees[i].Client.JMBG < *seminar.Trainees[j].Client.JMBG
	})

	maxConfNum, err := ClientSeminarService.GetMaxConfirmationNumber()
	if err != nil {
		logoped.ErrorLog.Println("Error getting max confirmation number: ", err)
		return []byte{}, err
	}

	for i, client := range seminar.Trainees {
		if _, exist := notPassedClientIds[client.ClientID]; exist {
			continue
		}
		pdf.AddPage()

		pdf.SetFont("Arimo-Regular", "", 11)

		if client.ConfirmationNumber == 0 {
			maxConfNum++
			client.ConfirmationNumber = maxConfNum
			//var t = true
			//switch client.Seminar.SeminarTheme.Code {
			//case "1":
			//	client.Client.PassedCheckboxes.Burden = &t
			//case "2":
			//	client.Client.PassedCheckboxes.Burden = &t
			//case "3":
			//	client.Client.PassedCheckboxes.Burden = &t
			//case "4":
			//	client.Client.PassedCheckboxes.Burden = &t
			//case "5":
			//	client.Client.PassedCheckboxes.Burden = &t
			//}
			_, err = ClientSeminarService.UpdateClientSeminar(client)
			if err != nil {
				return nil, err
			}
		}

		pdf.Ln(25)
		pdf.Text(15, pdf.GetY(), trObj.translDef("Број:"))
		pdf.Text(30, pdf.GetY(), trObj.translDef(seminar.GetCode()+"/"+strconv.Itoa(i+1)))
		pdf.Ln(5)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(15, pdf.GetY(), trObj.translDef("Дана:"))
		//pdf.SetFont("Arimo-Bold", "", 11)
		dateOfLastDay := seminar.Start
		maxNumber := 0
		for _, day := range seminar.Days {
			if day.Number > maxNumber {
				dateOfLastDay = day.Date
				maxNumber = day.Number
			}
		}
		pdf.Text(30, pdf.GetY(), dateOfLastDay.Format("02.01.2006."))
		pdf.Ln(5)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(15, pdf.GetY(), trObj.translDef("Место:"))
		//pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(30, pdf.GetY(), trObj.translDef(seminar.ClassRoom.Location.Address.Place))
		pdf.Ln(15)

		pdf.Text(95, pdf.GetY(), trObj.translDef("ПОТВРДА"))
		pdf.Ln(10)

		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(25, pdf.GetY(), trObj.translDef(fmt.Sprintf("о завршеној %s обуци на обавезним семинарима унапређења знања", seminar.SeminarTheme.BaseSeminarType.GetSeminarTypeForSentence())))
		pdf.Ln(10)

		ch := 9.0
		wl := 80.0
		wr := 100.0

		pdf.CellFormat(wl, ch-1, trObj.translDef("Име, име једног"), "LRT", 0, "L", false, 0, "")
		//pdf.SetFont("Arimo-Bold", "", 11)
		pdf.CellFormat(wr, ch-1, trObj.translDef(client.Client.Person.FirstName+" ("+client.Client.Person.MiddleName+") "+client.Client.Person.LastName), "LRT", 0, "L", false, 0, "")
		pdf.Ln(ch - 1)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.CellFormat(wl, ch-1, trObj.translDef("родитеља, презиме"), "LRB", 0, "L", false, 0, "")
		pdf.CellFormat(wr, ch-1, "", "LRB", 0, "L", false, 0, "")
		pdf.Ln(ch - 1)
		pdf.CellFormat(wl, ch, trObj.translDef("ЈМБГ"), "1", 0, "L", false, 0, "")
		//pdf.SetFont("Arimo-Bold", "", 11)
		for _, l := range trObj.translDef(*client.Client.JMBG) {
			pdf.CellFormat(wr/13, ch, string(l), "1", 0, "C", false, 0, "")
		}
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.CellFormat(wl, ch, trObj.translDef("Место пребивалишта"), "1", 0, "L", false, 0, "")
		//pdf.SetFont("Arimo-Bold", "", 11)
		pdf.CellFormat(wr, ch, trObj.translDef(client.Client.Address.Place), "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.CellFormat(wl, ch, trObj.translDef("Адреса пребивалишта"), "1", 0, "L", false, 0, "")
		//pdf.SetFont("Arimo-Bold", "", 11)
		pdf.CellFormat(wr, ch, trObj.translDef(client.Client.Address.Street+" "+client.Client.Address.HouseNumber), "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.CellFormat(wl, ch, trObj.translDef("Редни број семинара"), "1", 0, "L", false, 0, "")
		//pdf.SetFont("Arimo-Bold", "", 11)
		//fixme just closed seminars and passed
		completedSeminarsBeforeSrbolab := 0
		if client.Client.InitialCompletedSeminars != nil {
			completedSeminarsBeforeSrbolab = *client.Client.InitialCompletedSeminars
		}

		completedInSrbolab, err := ClientSeminarService.GetNumberOfPassedSeminars(client.ClientID)
		if err != nil {
			return nil, err
		}
		seminarNumber := completedInSrbolab + completedSeminarsBeforeSrbolab
		if seminarNumber > 5 {
			seminarNumber = seminarNumber - 5
		}
		cx := 89.5 + float64(seminarNumber)*7
		if seminarNumber == 1 {
			cx = cx + 1
		}
		if seminarNumber == 5 {
			cx = cx + 1
		}
		if seminarNumber > 0 {
			pdf.Circle(cx, 157, 3, "")
		}
		pdf.CellFormat(wr, ch, " I    II    III    IV    V", "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.CellFormat(wl, ch-1, trObj.translDef("Датум похађања"), "LRT", 0, "L", false, 0, "")
		//pdf.SetFont("Arimo-Bold", "", 11)
		pdf.CellFormat(wr, ch-1, trObj.translDef("од ")+startSeminar.Format("02.01.2006"), "LRT", 0, "L", false, 0, "")
		pdf.Ln(ch - 1)
		pdf.SetFont("Arimo-Regular", "", 11)
		seminarType := "периодичне"
		if seminar.SeminarTheme.BaseSeminarType.Code == "ADDITIONAL" {
			seminarType = "додатне"
		}
		if seminar.SeminarTheme.BaseSeminarType.Code == "BASE" {
			seminarType = "основне"
		}
		pdf.CellFormat(wl, ch-1, trObj.translDef(fmt.Sprintf("%s обуке", seminarType)), "LRB", 0, "L", false, 0, "")
		//pdf.SetFont("Arimo-Bold", "", 11)
		pdf.CellFormat(wr, ch-1, trObj.translDef("до ")+endSeminar.Format("02.01.2006"), "LRB", 0, "L", false, 0, "")
		pdf.Ln(ch - 1)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.CellFormat(wl, ch-1, trObj.translDef("Место похађања"), "LRT", 0, "L", false, 0, "")
		//pdf.SetFont("Arimo-Bold", "", 11)
		pdf.CellFormat(wr, ch-1, trObj.translDef(seminar.ClassRoom.Location.Address.Place+", "+seminar.ClassRoom.Location.Address.Street+" "+seminar.ClassRoom.Location.Address.HouseNumber), "LRT", 0, "L", false, 0, "")
		pdf.Ln(ch - 1)
		pdf.SetFont("Arimo-Regular", "", 11)

		pdf.CellFormat(wl, ch-1, trObj.translDef(fmt.Sprintf("%s обуке", seminarType)), "LRB", 0, "L", false, 0, "")
		pdf.CellFormat(wr, ch-1, "", "LRB", 0, "L", false, 0, "")
		pdf.Ln(ch - 1)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.CellFormat(wl, ch, trObj.translDef("Врста ЦПЦ"), "1", 0, "L", false, 0, "")
		//pdf.SetFont("Arimo-Bold", "", 11)
		if client.Client.CLicence != nil && *client.Client.CLicence {
			pdf.Circle(97.5, 198.5, 2.5, "")
		}
		if client.Client.DLicence != nil && *client.Client.DLicence {
			pdf.Circle(129.5, 198.5, 2.5, "")
		}
		pdf.CellFormat(wr, ch, trObj.translDef("1. превоз терета  2. превоз путника"), "1", 0, "L", false, 0, "")
		pdf.Ln(20)

		//pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(15, pdf.GetY(), trObj.translDef("НАПОМЕНА:"))
		pdf.SetFont("Arimo-Regular", "", 11)
		seminarType = "периодичне"
		if seminar.SeminarTheme.BaseSeminarType.Code == "ADDITIONAL" {
			seminarType = "додатне"
		}
		if seminar.SeminarTheme.BaseSeminarType.Code == "BASE" {
			seminarType = "основне"
		}
		pdf.Text(50, pdf.GetY(), trObj.translDef(fmt.Sprintf("Ова потврда се издаје на основу одслушане обавезне %s обуке", seminarType)))
		pdf.Ln(5)
		seminarType = "периодичног"
		if seminar.SeminarTheme.BaseSeminarType.Code == "ADDITIONAL" {
			seminarType = "додатног"
		}
		if seminar.SeminarTheme.BaseSeminarType.Code == "BASE" {
			seminarType = "основног"
		}
		pdf.Text(15, pdf.GetY(), trObj.translDef(fmt.Sprintf("за потребе стицања %s ЦПЦ и не може се користити у друге сврхе.", seminarType)))

		pdf.Ln(40)
		pdf.Text(15, pdf.GetY(), strconv.Itoa(client.ConfirmationNumber))
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
	pdf.AddFont("Helvetica", "", "helvetica_1251.json")
	latTr := pdf.UnicodeTranslatorFromDescriptor("iso-8859-16")
	cirTr := pdf.UnicodeTranslatorFromDescriptor("cp1251")
	trObj := newTranslationDetails(pdf, "Helvetica", "Arimo-Regular", 11, latTr, cirTr)

	pdf.SetMargins(15.0, 20, 15.0)

	dateOfLastDay := seminar.Start
	maxNumber := 0
	for _, day := range seminar.Days {
		if day.Number > maxNumber {
			dateOfLastDay = day.Date
			maxNumber = day.Number
		}
	}

	companyClientsMap := map[string][]model.ClientSeminar{}

	for _, client := range seminar.Trainees {
		if client.Client.Company.ID > 0 {
			val, ok := companyClientsMap[client.Client.Company.Name]
			if ok {
				val = append(val, client)
			} else {
				val = []model.ClientSeminar{client}
			}
			companyClientsMap[client.Client.Company.Name] = val

			continue
		}

		pdf.AddPage()

		createSimpleHeader(pdf, cirTr)

		//pdf.SetFont("Arimo-Bold", "", 12)
		pdf.Text(40, pdf.GetY(), trObj.translate(fmt.Sprintf("Изјава о преузимању потврде и завршеној %s обуци", seminar.SeminarTheme.BaseSeminarType.GetSeminarTypeForSentence()), 12))
		pdf.Ln(5)
		pdf.Text(60, pdf.GetY(), trObj.translate("на обавезним семинарима унапређења знања", 12))
		pdf.Ln(20)

		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(20, pdf.GetY(), trObj.translDef("Дана"))
		pdf.Line(33, pdf.GetY(), 60, pdf.GetY())
		//pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(33, pdf.GetY()-1, latTr(dateOfLastDay.Format("02.01.2006.")))
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(62, pdf.GetY(), trObj.translDef("године, "))
		pdf.Line(77, pdf.GetY(), 135, pdf.GetY())
		//pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(80, pdf.GetY()-1, trObj.translDef(client.Client.Person.FullName()))
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(137, pdf.GetY(), trObj.translDef("(име и презиме), ЈМБГ"))
		pdf.Ln(6)
		pdf.Line(20, pdf.GetY(), 55, pdf.GetY())
		//pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(25, pdf.GetY()-1, *client.Client.JMBG)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(55, pdf.GetY(), trObj.translDef(fmt.Sprintf("је преузео потврду о завршеној %s обуци на обавезним", seminar.SeminarTheme.BaseSeminarType.GetSeminarTypeForSentence())))
		pdf.Ln(6)
		pdf.Text(20, pdf.GetY(), trObj.translDef("семинарима унапређења знања."))

		pdf.Ln(80)

		pdf.Text(20, pdf.GetY(), trObj.translDef("Потврду преузео: "))
		pdf.Ln(10)
		pdf.Line(20, pdf.GetY(), 60, pdf.GetY())
		pdf.Ln(8)
		pdf.Text(20, pdf.GetY(), trObj.translDef("Дана: "))
		//pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Line(31, pdf.GetY(), 58, pdf.GetY())
		pdf.Text(32, pdf.GetY()-1, dateOfLastDay.Format("02.01.2006"))
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(60, pdf.GetY(), trObj.translDef("године."))
	}

	for company, clients := range companyClientsMap {
		pdf.AddPage()

		createSimpleHeader(pdf, cirTr)

		//pdf.SetFont("Arimo-Bold", "", 12)
		pdf.Text(40, pdf.GetY(), trObj.translate(fmt.Sprintf("Изјава о преузимању потврде и завршеној %s обуци", seminar.SeminarTheme.BaseSeminarType.GetSeminarTypeForSentence()), 12))
		pdf.Ln(5)
		pdf.Text(60, pdf.GetY(), trObj.translate("на обавезним семинарима унапређења знања", 12))
		pdf.Ln(20)

		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(15, pdf.GetY(), trObj.translDef("Дана"))
		pdf.Line(27, pdf.GetY(), 57, pdf.GetY())
		//pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(30, pdf.GetY()-1, trObj.translDef(dateOfLastDay.Format("02.01.2006.")))
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(60, pdf.GetY(), trObj.translDef("године, "))
		pdf.Line(75, pdf.GetY(), 135, pdf.GetY())
		pdf.Text(135, pdf.GetY(), trObj.translDef("(име и презиме), запослен у"))
		pdf.Ln(6)
		pdf.Text(15, pdf.GetY(), trObj.translDef("у фирми"))
		//pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(33, pdf.GetY()-1, trObj.translDef(company))
		pdf.Line(31.5, pdf.GetY(), 179, pdf.GetY())
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(180, pdf.GetY(), ",")
		pdf.Ln(6)
		pdf.Text(15, pdf.GetY(), trObj.translDef(("ЈМБГ")))
		pdf.Line(28, pdf.GetY(), 70, pdf.GetY())
		pdf.Text(70, pdf.GetY(), trObj.translDef(fmt.Sprintf(", је преузео потврде о завршеној %s обуци на", seminar.SeminarTheme.BaseSeminarType.GetSeminarTypeForSentence())))
		pdf.Ln(6)
		pdf.Text(15, pdf.GetY(), trObj.translDef("обавезним семинарима унапређења знања за следећа лица:"))
		pdf.Ln(10)

		ch := 5.0
		pdf.CellFormat(10, ch, trObj.translDef("РБ"), "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, ch, trObj.translDef("Име"), "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, ch, trObj.translDef("Презиме"), "1", 0, "L", false, 0, "")
		pdf.CellFormat(35, ch, trObj.translDef("ЈМБГ"), "1", 0, "L", false, 0, "")
		pdf.CellFormat(50, ch, trObj.translDef("Број потврде"), "1", 0, "L", false, 0, "")

		sort.Slice(clients, func(i, j int) bool {
			return *clients[i].Client.JMBG < *clients[j].Client.JMBG
		})

		for i, client := range clients {
			pdf.Ln(ch)
			pdf.CellFormat(10, ch, strconv.Itoa(i+1), "1", 0, "L", false, 0, "")
			pdf.CellFormat(40, ch, trObj.translDef(client.Client.Person.FirstName), "1", 0, "L", false, 0, "")
			pdf.CellFormat(40, ch, trObj.translDef(client.Client.Person.LastName), "1", 0, "L", false, 0, "")
			pdf.CellFormat(35, ch, *client.Client.JMBG, "1", 0, "L", false, 0, "")
			pdf.CellFormat(50, ch, trObj.translDef(seminar.GetCode())+"/"+strconv.Itoa(i+1), "1", 0, "L", false, 0, "")
		}

		pdf.Ln(20)
		pdf.Text(15, pdf.GetY(), trObj.translDef("Потврду преузео: "))
		pdf.Ln(10)
		pdf.Line(15, pdf.GetY(), 70, pdf.GetY())
		pdf.Ln(6)
		pdf.Text(15, pdf.GetY(), trObj.translDef("Број ЛК:"))
		pdf.Line(30, pdf.GetY(), 70, pdf.GetY())
		pdf.Ln(6)
		pdf.Text(15, pdf.GetY(), trObj.translDef("Дана: "))
		//pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Line(26, pdf.GetY(), 48, pdf.GetY())
		pdf.Text(27, pdf.GetY()-1, dateOfLastDay.Format("02.01.2006"))
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(50, pdf.GetY(), trObj.translDef("година."))
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
	pdf.AddFont("Arimo-Bold", "", "Arimo-Bold.json")
	pdf.AddFont("Helvetica", "", "helvetica_1251.json")
	latTr := pdf.UnicodeTranslatorFromDescriptor("iso-8859-16")
	cirTr := pdf.UnicodeTranslatorFromDescriptor("cp1251")
	fontSize := 10.0
	trObj := newTranslationDetails(pdf, "Helvetica", "Arimo-Regular", fontSize, latTr, cirTr)

	pdf.SetMargins(marginLeft, marginTop, marginRight)
	pdf.AddPage()

	pdf.SetFont("Arimo-Regular", "", fontSize)
	createSimpleHeaderForLandscape(pdf, cirTr)

	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), trObj.translDef("Место: "))
	pdf.Text(27, pdf.GetY(), trObj.translDef(day.Seminar.ClassRoom.Location.Address.Place))
	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), trObj.translDef("Шифра обуке: "))
	pdf.Text(39, pdf.GetY(), trObj.translDef(day.Seminar.GetCode()))
	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), trObj.translDef("Датум: "))
	dayInWeek := util.GetDaySerbian(day.Date)
	pdf.Text(27, pdf.GetY(), trObj.translDef(dayInWeek)+" "+day.Date.Format("02.01.2006."))
	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), trObj.translDef(fmt.Sprintf("Прозивник полазника семинара унапређења знања на %s обуци професионалних возача", day.Seminar.SeminarTheme.BaseSeminarType.GetSeminarTypeForSentence())))

	ch := 14.0
	pdf.Ln(2)
	pdf.CellFormat(10, ch, trObj.translDef("Р.Б."), "1", 0, "C", false, 0, "")
	pdf.CellFormat(45, ch, trObj.translDef("Име и презиме / ЈМБГ"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(27, ch, trObj.translDef("1. час"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(27, ch, trObj.translDef("2. час"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(27, ch, trObj.translDef("3. час"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(27, ch, trObj.translDef("4. час"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(27, ch, trObj.translDef("5. час"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(27, ch, trObj.translDef("6. час"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(27, ch, trObj.translDef("7. час"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(31, ch, trObj.translDef("Напомена"), "1", 0, "C", false, 0, "")

	sort.Slice(day.Presence, func(i, j int) bool {
		return *day.Presence[i].Client.JMBG < *day.Presence[j].Client.JMBG
	})

	for i, cs := range day.Presence {
		pdf.Ln(ch)
		if pdf.GetY() > 177 {
			pdf.Text(270, 200, strconv.Itoa(pdf.PageCount()))
			pdf.CellFormat(10, ch, trObj.translDef("Р.Б."), "1", 0, "C", false, 0, "")
			pdf.CellFormat(45, ch, trObj.translDef("Име и презиме / ЈМБГ"), "1", 0, "C", false, 0, "")
			pdf.CellFormat(27, ch, trObj.translDef("1. час"), "1", 0, "C", false, 0, "")
			pdf.CellFormat(27, ch, trObj.translDef("2. час"), "1", 0, "C", false, 0, "")
			pdf.CellFormat(27, ch, trObj.translDef("3. час"), "1", 0, "C", false, 0, "")
			pdf.CellFormat(27, ch, trObj.translDef("4. час"), "1", 0, "C", false, 0, "")
			pdf.CellFormat(27, ch, trObj.translDef("5. час"), "1", 0, "C", false, 0, "")
			pdf.CellFormat(27, ch, trObj.translDef("6. час"), "1", 0, "C", false, 0, "")
			pdf.CellFormat(27, ch, trObj.translDef("7. час"), "1", 0, "C", false, 0, "")
			pdf.CellFormat(31, ch, trObj.translDef("Напомена"), "1", 0, "C", false, 0, "")
			pdf.Ln(ch)
		}
		pdf.CellFormat(10, ch, strconv.Itoa(i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(45, ch, "", "1", 0, "C", false, 0, "")
		pdf.Text(22, pdf.GetY()+5, trObj.translDef(cs.Client.Person.FullName()))
		pdf.Text(22, pdf.GetY()+10, *cs.Client.JMBG)
		pdf.CellFormat(27, ch, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(27, ch, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(27, ch, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(27, ch, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(27, ch, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(27, ch, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(27, ch, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(31, ch, "", "1", 0, "C", false, 0, "")
	}

	if pdf.GetY() < 177 {
		pdf.Text(270, 200, strconv.Itoa(pdf.PageCount()))
	}

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (p *printService) PrintCheckIn(seminar *model.Seminar) ([]byte, error) {
	pwd, err := os.Getwd()
	if err != nil {
		logoped.ErrorLog.Println("Error getting pwd: ", err)
		return []byte{}, err
	}
	pdf := fpdf.New("P", "mm", "A4", filepath.Join(pwd, "font"))
	pdf.AddFont("Arimo-Regular", "", "Arimo-Regular.json")
	pdf.AddFont("Arimo-Bold", "", "Arimo-Bold.json")
	pdf.AddFont("Helvetica", "", "helvetica_1251.json")
	latTr := pdf.UnicodeTranslatorFromDescriptor("iso-8859-16")
	cirTr := pdf.UnicodeTranslatorFromDescriptor("cp1251")
	trObj := newTranslationDetails(pdf, "Helvetica", "Arimo-Regular", 11, latTr, cirTr)

	pdf.SetMargins(15, 20, marginRight)
	fontSize := 11.0
	ch := 6.0

	for _, client := range seminar.Trainees {
		fmt.Println(client.ClientID)
		pdf.AddPage()

		//pdf.SetFont("Arimo-Bold", "", 15)
		pdf.Text(85, pdf.GetY(), trObj.translate("П Р И Ј А В А*", 15))

		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.Ln(10)
		pdf.Text(70, pdf.GetY(), trObj.translDef("за похађање обавезног семинара"))
		pdf.Ln(5)
		pdf.Text(84, pdf.GetY(), trObj.translDef("унапређења знања"))
		pdf.Ln(10)

		//pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.Text(15, pdf.GetY(), trObj.translDef("ЛИЧНИ ПОДАЦИ"))
		pdf.Ln(1)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("Име (име једног родитеља) презиме:"), "1", 0, "L", false, 0, "")
		//pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.CellFormat(110, ch, trObj.translDef(client.Client.Person.FullNameWithMiddleName()), "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("ЈМБГ:"), "1", 0, "L", false, 0, "")
		//pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.CellFormat(110, ch, *client.Client.JMBG, "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("Датум рођења:"), "1", 0, "L", false, 0, "")
		//pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.CellFormat(110, ch, client.Client.GetBirthDate(), "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("Место рођења, држава:"), "1", 0, "L", false, 0, "")
		//pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.CellFormat(110, ch, trObj.translDef(*client.Client.PlaceBirth+", "+*client.Client.CountryBirth), "1", 0, "L", false, 0, "")

		pdf.Ln(17)
		//pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.Text(15, pdf.GetY(), trObj.translDef("ПОДАЦИ О ПРЕБИВАЛИШТУ/БОРАВИШТУ"))
		pdf.Ln(1)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("Место:"), "1", 0, "L", false, 0, "")
		//pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.CellFormat(110, ch, trObj.translDef(client.Client.Address.Place), "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("Поштански број:"), "1", 0, "L", false, 0, "")
		//pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.CellFormat(110, ch, trObj.translDef(client.Client.Address.PostCode), "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("Улица и кућни број:"), "1", 0, "L", false, 0, "")
		//pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.CellFormat(110, ch, trObj.translDef(client.Client.Address.GetStreetWithNumber()), "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("Телефон/Мобилни:"), "1", 0, "L", false, 0, "")
		//pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.CellFormat(110, ch, trObj.translDef(client.Client.Person.PhoneNumber), "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("е-маил:"), "1", 0, "L", false, 0, "")
		//pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.CellFormat(110, ch, trObj.translDef(client.Client.Person.Email), "1", 0, "L", false, 0, "")

		pdf.Ln(17)
		//pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.Text(15, pdf.GetY(), trObj.translDef("ПОДАЦИ О КВАЛИФИКАЦИОНОЈ КАРТИЦИ ВОЗАЧА"))
		pdf.Ln(1)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("Број картице(СРБ број)*:"), "1", 0, "L", false, 0, "")
		//pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.CellFormat(110, ch, *client.Client.CPCNumber, "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("Серијски број картице(СРБ број):"), "1", 0, "L", false, 0, "")
		//pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.CellFormat(110, ch, "", "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("Рок важења картице:"), "1", 0, "L", false, 0, "")
		//pdf.SetFont("Arimo-Bold", "", fontSize)
		cpcDate := ""
		if client.Client.CPCDate != nil {
			cpcDate = client.Client.CPCDate.Format("02.01.2006.")
		}
		pdf.CellFormat(110, ch, cpcDate, "1", 0, "L", false, 0, "")

		pdf.Ln(17)
		//pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.Text(15, pdf.GetY(), trObj.translDef("ВРСТА ПРЕВОЗА (заокружени број испред врсте превоза)"))
		pdf.Ln(1)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		if client.Client.CLicence != nil && *client.Client.CLicence {
			pdf.Circle(18, 157, 2.5, "")
		}
		if client.Client.DLicence != nil && *client.Client.DLicence {
			pdf.Circle(107.5, 157, 2.5, "")
		}

		pdf.CellFormat(90, ch, trObj.translDef("1. Превоз терета"), "1", 0, "L", false, 0, "")
		pdf.CellFormat(90, ch, trObj.translDef("2. Превоз путника"), "1", 0, "L", false, 0, "")

		pdf.Ln(15)
		pdf.Text(15, pdf.GetY(), trObj.translDef("Уз попуњен образац Пријаве за похађање семинара, прилаже се:"))
		pdf.Ln(5)
		pdf.Text(20, pdf.GetY(), trObj.translDef("- доказ о уплати трошкова за похађање семинара, по важећој тарифи."))

		pdf.Ln(20)
		pdf.Text(18, pdf.GetY(), trObj.translDef("У "))
		//pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.Text(23, pdf.GetY(), trObj.translDef(seminar.ClassRoom.Location.Address.Place))
		pdf.Line(23, pdf.GetY()+1, 65, pdf.GetY()+1)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.Text(65.5, pdf.GetY(), trObj.translDef(", дана"))
		//pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.Line(78, pdf.GetY()+1, 110, pdf.GetY()+1)
		pdf.Text(80, pdf.GetY(), seminar.Start.Format("02.01.2006."))
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.Text(110.5, pdf.GetY(), trObj.translDef(", године"))

		pdf.Ln(30)
		pdf.Text(135, pdf.GetY(), trObj.translDef("Потпис подносиоца пријаве: "))
		pdf.Ln(15)
		pdf.Line(135, pdf.GetY(), 190, pdf.GetY())
		pdf.Ln(10)

		pdf.SetFont("Arimo-Regular", "", 9)
		pdf.Text(15, pdf.GetY(), trObj.translate("*Образац пријаве попунити читко штампаним словима", 9))
		pdf.Ln(5)
		pdf.Text(15, pdf.GetY(), trObj.translate("*Уписати број картице (СРБ број) или број „Потврде о пријему захтева за издавање сертификата о стручној", 9))
		pdf.Ln(4)
		pdf.Text(17, pdf.GetY(), trObj.translate("компетенцији и квалификационе картице возача, уколико сте покренули поступак издавања", 9))
		pdf.Ln(4)
		pdf.Text(17, pdf.GetY(), trObj.translate("квалификаионе картице и сертификата", 9))

		//---------------------------------------------------
		//confirmation statment
		pdf.AddPage()
		pdf.SetFont("Arimo-Regular", "", 11)

		pdf.Ln(5)
		pdf.Text(10, pdf.GetY(), trObj.translDef("У складу са чланом 6. и чланом 7. уредбе ЕУ 2016/679 од 27. априла 2016. године и чланом 12. и чланом"))
		pdf.Ln(5)
		pdf.Text(10, pdf.GetY(), trObj.translDef("15. Закона о заштити пдатака о личности (СЛ. Гласник РС“, бр. 87/2018 од 13/11/2018) дајем пристанак за"))
		pdf.Ln(5)
		pdf.Text(10, pdf.GetY(), trObj.translDef("давање и обраду података о личности, где је руковалац обраде Срболаб."))
		pdf.Ln(20)

		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(80, pdf.GetY(), trObj.translDef("ИЗЈАВА О ПРИСТАНКУ"))
		pdf.Ln(20)

		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(10, pdf.GetY(), trObj.translDef("Ја"))
		pdf.Line(17, pdf.GetY(), 77, pdf.GetY())
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(18, pdf.GetY()-1, trObj.translDef(client.Client.Person.FullName()))
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(77, pdf.GetY(), trObj.translDef("(име и презиме),"))
		pdf.Line(108, pdf.GetY(), 165, pdf.GetY())
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(110, pdf.GetY()-1, trObj.translDef(*client.Client.JMBG))
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(167, pdf.GetY(), trObj.translDef("(ЈМБГ),"))

		pdf.Ln(5)
		pdf.SetFont("Arimo-Regular", "", 10)
		pdf.Text(10, pdf.GetY(), trObj.translate("пристајем на давање и обраду следећих својих података о личности: подаци из личне карте/пасоша, подаци из", 10))
		pdf.Ln(5)
		pdf.Text(10, pdf.GetY(), trObj.translate("националне возачке дозволе, подаци из квалификационе картице возача, електронску адресу, контакт телефон", 10))
		pdf.Ln(5)
		pdf.Text(10, pdf.GetY(), trObj.translate("возача, подаци о стручној спреми, за потребе слања обавештења и информација.", 10))
		pdf.Ln(15)

		pdf.Text(10, pdf.GetY(), trObj.translate("Такође изјављујем да сам од Срболаб д.о.о. примио/ла сва неопходна обавештења, предвиђена чланом 23", 10))
		pdf.Ln(5)
		pdf.Text(10, pdf.GetY(), trObj.translate("Закона о заштити података о личности, као и обавештење да у сваком тренутку могу опозвати дат", 10))
		pdf.Ln(5)
		pdf.Text(10, pdf.GetY(), trObj.translate("пристанак, с тим да опозив пристанка не утиче на допуштеност обраде која је вршена на основу пристанка", 10))
		pdf.Ln(5)
		pdf.Text(10, pdf.GetY(), trObj.translate("пре опозива, као и да нисам у обавези да дам податке о личности који нису предвиђени као обавезни", 10))
		pdf.Ln(5)
		pdf.Text(10, pdf.GetY(), trObj.translate("законским и подзаконским актима и да исто неће бити од утицаја на прижање услуга од стране руковаоца.", 10))
		pdf.Ln(25)

		pdf.Text(15, pdf.GetY(), trObj.translDef("Датум: "))
		pdf.Line(30, pdf.GetY(), 70, pdf.GetY())
		//pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(35, pdf.GetY()-1, seminar.Start.Format("02.01.2006"))
		pdf.SetFont("Arimo-Regular", "", 11)

		pdf.Ln(15)

		pdf.Text(15, pdf.GetY(), trObj.translDef("Потпис: "))
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
	pdf.SetFont("Helvetica", "", 10)
	pdf.Image("./images/srbolab_logo.png", 15, 10, 30, 10, false, "png", 0, "")
	pdf.CellFormat(35, 10, "", "0", 0, "C", false, 0, "")
	pdf.Text(95, 14, tr("СРБОЛАБ Д.О.О."))
	pdf.Text(80, 18, tr("Центар за едукацију и развој Срболаб"))
	pdf.Image("./images/cers_logo.png", 170, 10, 20, 10, false, "png", 0, "")
	pdf.Ln(15)
}

func createSimpleHeaderForLandscape(pdf *fpdf.Fpdf, tr func(string) string) {
	pdf.SetFont("Helvetica", "", 10)
	pdf.Image("./images/srbolab_logo.png", 15, 10, 30, 10, false, "png", 0, "")
	pdf.CellFormat(35, 10, "", "0", 0, "C", false, 0, "")
	pdf.Text(135, 14, tr("СРБОЛАБ Д.О.О."))
	pdf.Text(120, 18, tr("Центар за едукацију и развој Срболаб"))
	pdf.Image("./images/cers_logo.png", 260, 10, 20, 10, false, "png", 0, "")
	pdf.Ln(15)
}

func (p *printService) PrintSeminarEvidence(day *model.SeminarDay) ([]byte, error) {
	pwd, err := os.Getwd()
	if err != nil {
		logoped.ErrorLog.Println("Error getting pwd: ", err)
		return []byte{}, err
	}
	pdf := fpdf.New("L", "mm", "A4", filepath.Join(pwd, "font"))
	pdf.AddFont("Arimo-Regular", "", "Arimo-Regular.json")
	pdf.AddFont("Arimo-Bold", "", "Arimo-Bold.json")
	pdf.AddFont("Helvetica", "", "helvetica_1251.json")
	latTr := pdf.UnicodeTranslatorFromDescriptor("iso-8859-16")
	cirTr := pdf.UnicodeTranslatorFromDescriptor("cp1251")
	trObj := newTranslationDetails(pdf, "Helvetica", "Arimo-Regular", 9, latTr, cirTr)

	pdf.SetMargins(marginLeft, marginTop, marginRight)
	pdf.AddPage()

	pdf.SetFont("Arimo-Regular", "", 9)
	createSimpleHeaderForLandscape(pdf, cirTr)

	pdf.Ln(5)
	pdf.SetFont("Arimo-Bold", "", 9)
	pdf.Text(15, pdf.GetY(), trObj.translDef(fmt.Sprintf("Дневник предавача семинара унапређења знања на %s обуци професионалних возача", day.Seminar.SeminarTheme.BaseSeminarType.GetSeminarTypeForSentence())))
	pdf.Ln(5)
	pdf.SetFont("Arimo-Regular", "", 9)
	pdf.Text(15, pdf.GetY(), trObj.translDef("Датум одржавања семинара"))
	pdf.Text(60, pdf.GetY(), day.Date.Format("02.01.2006."))
	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), trObj.translDef(day.Seminar.ClassRoom.Location.Address.Place))
	pdf.Ln(2)

	ch := 10.0

	pdf.CellFormat(20, ch, "", "1", 0, "TC", false, 0, "")
	pdf.Text(13, pdf.GetY()+3.5, trObj.translDef("Редни број"))
	pdf.Text(17, pdf.GetY()+8.5, trObj.translDef("часа"))
	pdf.CellFormat(40, ch, trObj.translDef("Време одржавања часа"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(90, ch, trObj.translDef("Наставни час"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(70, ch, "", "1", 0, "TC", false, 0, "")
	pdf.Text(178, pdf.GetY()+3.5, trObj.translDef("Име и презиме предавача"))
	pdf.Line(160, pdf.GetY()+5, 230, pdf.GetY()+5)
	pdf.Text(182, pdf.GetY()+8.5, trObj.translDef("Потпис предавача"))
	pdf.CellFormat(50, ch, trObj.translDef("Напомена"), "1", 0, "C", false, 0, "")

	sort.Slice(day.Classes, func(i, j int) bool {
		return day.Classes[i].Number < day.Classes[j].Number
	})

	for i := 0; i < len(day.Classes); i++ {
		pdf.Ln(ch)
		lines, _ := splitLine(day.Classes[i].Name, 55)
		current := pdf.GetY() + 3.5
		for i, line := range lines {
			pdf.Text(72, current+float64(i)*4.0, trObj.translDef(line))
		}
		trObj.translDef(day.Classes[i].Name)
		pdf.CellFormat(20, ch, strconv.Itoa(i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, ch, getClassTime(day.Date, i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(90, ch, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(70, ch, "", "1", 0, "C", false, 0, "")
		teacher := ""
		if day.Classes[i].Teacher != nil {
			teacher = day.Classes[i].Teacher.Person.FullName()
		}
		pdf.Text(178, pdf.GetY()+3.5, trObj.translDef(teacher))
		pdf.Line(160, pdf.GetY()+5, 230, pdf.GetY()+5)
		pdf.Text(182, pdf.GetY()+8.5, "")
		pdf.CellFormat(50, ch, "", "1", 0, "C", false, 0, "")
	}

	pdf.Ln(15)
	pdf.Text(100, pdf.GetY(), trObj.translDef("Шифра обуке:"))
	pdf.Text(122, pdf.GetY(), trObj.translDef(day.Seminar.GetCode()))

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (p *printService) PrintTestBarcode() ([]byte, error) {
	url := fmt.Sprintf("%s/do-test", Domain)
	qrCode, _ := qrcode.New(url, qrcode.Medium)
	var buf bytes.Buffer

	err := qrCode.Write(500, &buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (p *printService) PrintPlanTreningRealization(day *model.SeminarDay) ([]byte, error) {
	pwd, err := os.Getwd()
	if err != nil {
		logoped.ErrorLog.Println("Error getting pwd: ", err)
		return []byte{}, err
	}
	pdf := fpdf.New("P", "mm", "A4", filepath.Join(pwd, "font"))
	pdf.AddFont("Arimo-Regular", "", "Arimo-Regular.json")
	pdf.AddFont("Arimo-Bold", "", "Arimo-Bold.json")
	pdf.AddFont("Helvetica", "", "helvetica_1251.json")
	latTr := pdf.UnicodeTranslatorFromDescriptor("iso-8859-16")
	cirTr := pdf.UnicodeTranslatorFromDescriptor("cp1251")
	trObj := newTranslationDetails(pdf, "Helvetica", "Arimo-Regular", 11, latTr, cirTr)

	pdf.SetMargins(15.0, marginTop, marginRight)
	pdf.AddPage()

	pdf.Ln(30)
	pdf.SetFont("Arimo-Bold", "", 14)
	seminarTypeSentence := "периодичну"
	if day.Seminar.SeminarTheme.BaseSeminarType.Code == "ADDITIONAL" {
		seminarTypeSentence = "додатну"
	}
	if day.Seminar.SeminarTheme.BaseSeminarType.Code == "BASE" {
		seminarTypeSentence = "основну"
	}
	pdf.Text(38, pdf.GetY(), trObj.translate(fmt.Sprintf("План реализације наставе за %s обуку - 7 часова", seminarTypeSentence), 14))
	pdf.Ln(5)

	pdf.SetFont("Arimo-Bold", "", 11)
	ch := 10.0
	chs := 7.0
	cw1 := 12.0
	cw2 := 81.0
	cw3 := 45.0
	cw4 := 21.0
	cw5 := 21.0
	pdf.CellFormat(180, ch, trObj.translDef(util.GetDaySerbian(day.Date))+" "+day.Date.Format("02.01.2006"), "1", 0, "C", false, 0, "")
	pdf.Ln(ch)
	pdf.SetFillColor(180, 197, 231)
	pdf.CellFormat(cw1, ch, trObj.translDef("Р.бр"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw2, ch, trObj.translDef("Назив часа"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw3, ch, trObj.translDef("Предавач"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw4, ch, trObj.translDef("Почетак"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw5, ch, trObj.translDef("Крај"), "1", 0, "C", true, 0, "")
	pdf.SetFont("Arimo-Regular", "", 11)
	pdf.Ln(ch)
	splitWidth := 40
	lines, num := splitLine("Пријава и евидентирање полазника обуке", splitWidth)
	current := pdf.GetY() + 4.5
	for i, line := range lines {
		pdf.Text(30, current+float64(i)*6.0, trObj.translDef(line))
	}
	pdf.CellFormat(cw1, chs*num, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw2, chs*num, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw3, chs*num, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw4, chs*num, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw5, chs*num, "", "1", 0, "C", false, 0, "")

	pdf.Ln(chs * num)
	class1 := model.GetSeminarClassByNumber(day.Classes, 1)
	className := ""
	teacher := ""
	timeStart := ""
	if class1 != nil {
		className = class1.Name
		if class1.Teacher != nil {
			teacher = class1.Teacher.Person.FullName()
		}
		timeStart = getClassTime(day.Date, 1)
	}
	lines, num = splitLine(className, splitWidth)
	current = pdf.GetY() + 4.5
	for i, line := range lines {
		pdf.Text(30, current+float64(i)*6.0, trObj.translDef(line))
	}
	pdf.CellFormat(cw1, chs*num, "1", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw2, chs*num, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw3, chs*num, trObj.translDef(teacher), "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw4, chs*num, strings.Split(timeStart, "-")[0], "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw5, chs*num, strings.Split(timeStart, "-")[1], "1", 0, "C", false, 0, "")

	//pdf.SetFont("Arimo-Bold", "", 11)
	pdf.Ln(chs * num)
	pdf.SetFillColor(180, 197, 231)
	pdf.CellFormat(cw1, ch, latTr(""), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw2, ch, trObj.translDef("Пауза за кафу"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw3, ch, trObj.translDef(""), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw4, ch, latTr("10"), "LTB", 0, "R", true, 0, "")
	pdf.CellFormat(cw5, ch, trObj.translDef("минута"), "TBR", 0, "L", true, 0, "")

	pdf.SetFont("Arimo-Regular", "", 11)
	pdf.Ln(ch)
	class2 := model.GetSeminarClassByNumber(day.Classes, 2)
	className = ""
	teacher = ""
	timeStart = ""
	if class2 != nil {
		className = class2.Name
		if class2.Teacher != nil {
			teacher = class2.Teacher.Person.FullName()
		}
		timeStart = getClassTime(day.Date, 2)
	}
	lines, num = splitLine(className, splitWidth)
	current = pdf.GetY() + 4.5
	for i, line := range lines {
		pdf.Text(30, current+float64(i)*6.0, trObj.translDef(line))
	}
	pdf.CellFormat(cw1, chs*num, "2", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw2, chs*num, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw3, chs*num, trObj.translDef(teacher), "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw4, chs*num, strings.Split(timeStart, "-")[0], "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw5, chs*num, strings.Split(timeStart, "-")[1], "1", 0, "C", false, 0, "")

	//pdf.SetFont("Arimo-Bold", "", 11)
	pdf.Ln(chs * num)
	pdf.SetFillColor(180, 197, 231)
	pdf.CellFormat(cw1, ch, trObj.translDef(""), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw2, ch, trObj.translDef("Пауза"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw3, ch, trObj.translDef(""), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw4, ch, trObj.translDef("5"), "LTB", 0, "R", true, 0, "")
	pdf.CellFormat(cw5, ch, trObj.translDef("минута"), "TBR", 0, "L", true, 0, "")

	pdf.SetFont("Arimo-Regular", "", 11)
	pdf.Ln(ch)
	class3 := model.GetSeminarClassByNumber(day.Classes, 3)
	className = ""
	teacher = ""
	timeStart = ""
	if class3 != nil {
		className = class3.Name
		if class3.Teacher != nil {
			teacher = class3.Teacher.Person.FullName()
		}
		timeStart = getClassTime(day.Date, 3)
	}
	lines, num = splitLine(className, splitWidth)
	current = pdf.GetY() + 4.5
	for i, line := range lines {
		pdf.Text(30, current+float64(i)*6.0, trObj.translDef(line))
	}
	pdf.CellFormat(cw1, chs*num, "3", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw2, chs*num, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw3, chs*num, trObj.translDef(teacher), "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw4, chs*num, strings.Split(timeStart, "-")[0], "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw5, chs*num, strings.Split(timeStart, "-")[1], "1", 0, "C", false, 0, "")

	//pdf.SetFont("Arimo-Bold", "", 11)
	pdf.Ln(chs * num)
	pdf.SetFillColor(180, 197, 231)
	pdf.CellFormat(cw1, ch, trObj.translDef(""), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw2, ch, trObj.translDef("Пауза за доручак"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw3, ch, trObj.translDef(""), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw4, ch, trObj.translDef("25"), "LTB", 0, "R", true, 0, "")
	pdf.CellFormat(cw5, ch, trObj.translDef("минута"), "TBR", 0, "L", true, 0, "")

	pdf.SetFont("Arimo-Regular", "", 11)
	pdf.Ln(ch)
	class4 := model.GetSeminarClassByNumber(day.Classes, 4)
	className = ""
	teacher = ""
	timeStart = ""
	if class4 != nil {
		className = class4.Name
		if class4.Teacher != nil {
			teacher = class4.Teacher.Person.FullName()
		}
		timeStart = getClassTime(day.Date, 4)
	}
	lines, num = splitLine(className, splitWidth)
	current = pdf.GetY() + 4.5
	for i, line := range lines {
		pdf.Text(30, current+float64(i)*6.0, trObj.translDef(line))
	}
	pdf.CellFormat(cw1, chs*num, "4", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw2, chs*num, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw3, chs*num, trObj.translDef(teacher), "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw4, chs*num, strings.Split(timeStart, "-")[0], "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw5, chs*num, strings.Split(timeStart, "-")[1], "1", 0, "C", false, 0, "")

	//pdf.SetFont("Arimo-Bold", "", 11)
	pdf.Ln(chs * num)
	pdf.SetFillColor(180, 197, 231)
	pdf.CellFormat(cw1, ch, trObj.translDef(""), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw2, ch, trObj.translDef("Пауза"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw3, ch, trObj.translDef(""), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw4, ch, trObj.translDef("5"), "LTB", 0, "R", true, 0, "")
	pdf.CellFormat(cw5, ch, trObj.translDef("минута"), "TBR", 0, "L", true, 0, "")

	pdf.SetFont("Arimo-Regular", "", 11)
	pdf.Ln(ch)
	class5 := model.GetSeminarClassByNumber(day.Classes, 5)
	className = ""
	teacher = ""
	timeStart = ""
	if class5 != nil {
		className = class5.Name
		if class5.Teacher != nil {
			teacher = class5.Teacher.Person.FullName()
		}
		timeStart = getClassTime(day.Date, 5)
	}
	lines, num = splitLine(className, splitWidth)
	current = pdf.GetY() + 4.5
	for i, line := range lines {
		pdf.Text(30, current+float64(i)*6.0, trObj.translDef(line))
	}
	pdf.CellFormat(cw1, chs*num, "5", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw2, chs*num, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw3, chs*num, trObj.translDef(teacher), "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw4, chs*num, strings.Split(timeStart, "-")[0], "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw5, chs*num, strings.Split(timeStart, "-")[1], "1", 0, "C", false, 0, "")

	//pdf.SetFont("Arimo-Bold", "", 11)
	pdf.Ln(chs * num)
	pdf.SetFillColor(180, 197, 231)
	pdf.CellFormat(cw1, ch, trObj.translDef(""), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw2, ch, trObj.translDef("Пауза за кафу"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw3, ch, trObj.translDef(""), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw4, ch, trObj.translDef("15"), "LTB", 0, "R", true, 0, "")
	pdf.CellFormat(cw5, ch, trObj.translDef("минута"), "TBR", 0, "L", true, 0, "")

	pdf.SetFont("Arimo-Regular", "", 11)
	pdf.Ln(ch)
	class6 := model.GetSeminarClassByNumber(day.Classes, 6)
	className = ""
	teacher = ""
	timeStart = ""
	if class6 != nil {
		className = class6.Name
		if class6.Teacher != nil {
			teacher = class6.Teacher.Person.FullName()
		}
		timeStart = getClassTime(day.Date, 6)
	}
	lines, num = splitLine(className, splitWidth)
	current = pdf.GetY() + 4.5
	for i, line := range lines {
		pdf.Text(30, current+float64(i)*6.0, trObj.translDef(line))
	}
	pdf.CellFormat(cw1, chs*num, "6", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw2, chs*num, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw3, chs*num, trObj.translDef(teacher), "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw4, chs*num, strings.Split(timeStart, "-")[0], "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw5, chs*num, strings.Split(timeStart, "-")[1], "1", 0, "C", false, 0, "")

	//pdf.SetFont("Arimo-Bold", "", 11)
	pdf.Ln(chs * num)
	pdf.SetFillColor(180, 197, 231)
	pdf.CellFormat(cw1, ch, trObj.translDef(""), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw2, ch, trObj.translDef("Пауза"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw3, ch, trObj.translDef(""), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw4, ch, trObj.translDef("5"), "LTB", 0, "R", true, 0, "")
	pdf.CellFormat(cw5, ch, trObj.translDef("минута"), "TBR", 0, "L", true, 0, "")

	pdf.SetFont("Arimo-Regular", "", 10)
	pdf.Ln(ch)
	class7 := model.GetSeminarClassByNumber(day.Classes, 7)
	className = ""
	teacher = ""
	timeStart = ""
	if class7 != nil {
		className = class7.Name
		if class7.Teacher != nil {
			teacher = class7.Teacher.Person.FullName()
		}
		timeStart = getClassTime(day.Date, 7)
	}
	lines, num = splitLine(className, splitWidth)
	current = pdf.GetY() + 4.5
	for i, line := range lines {
		pdf.Text(30, current+float64(i)*6.0, trObj.translDef(line))
	}
	pdf.CellFormat(cw1, chs*num, "7", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw2, chs*num, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw3, chs*num, trObj.translDef(teacher), "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw4, chs*num, strings.Split(timeStart, "-")[0], "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw5, chs*num, strings.Split(timeStart, "-")[1], "1", 0, "C", false, 0, "")

	pdf.Ln(chs * num)
	lines, num = splitLine("Евалуација наставног процеса и пријем документације", splitWidth)
	current = pdf.GetY() + 4.5
	for i, line := range lines {
		pdf.Text(30, current+float64(i)*6.0, trObj.translDef(line))
	}
	start := day.Date.Add(375 * time.Minute).Format("15:04")
	end := day.Date.Add(390 * time.Minute).Format("15:04")
	pdf.CellFormat(cw1, chs*num, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw2, chs*num, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw3, chs*num, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw4, chs*num, start, "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw5, chs*num, end, "1", 0, "C", false, 0, "")

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (p *printService) PrintPayments(seminar *model.Seminar) ([]byte, error) {
	pwd, err := os.Getwd()
	if err != nil {
		logoped.ErrorLog.Println("Error getting pwd: ", err)
		return []byte{}, err
	}
	pdf := fpdf.New("P", "mm", "A4", filepath.Join(pwd, "font"))

	pdf.SetMargins(15.0, marginTop, marginRight)

	h := 70.0
	for _, t := range seminar.Trainees {
		for _, doc := range t.Client.Documents {
			if (len(doc.Name) == 1 && doc.Name == strconv.Itoa(int(seminar.ID))) || strings.HasPrefix(doc.Name, strconv.Itoa(int(seminar.ID))+".") {
				pdf.AddPage()
				pdf.Ln(30)
				splitedName := strings.Split(doc.Name, ".")

				info := pdf.RegisterImage(FileService.GetFullPath(ClientFolder, strconv.Itoa(int(t.ClientID)), doc.Name), splitedName[len(splitedName)-1])
				if info != nil && info.Width() > 0 && info.Height() > 0 {
					scale := info.Width() / info.Height()
					h = 180 / scale
				}

				pdf.Image(FileService.GetFullPath(ClientFolder, strconv.Itoa(int(t.ClientID)), doc.Name), 10, 10, 180, h, false, strings.Split(doc.Name, ".")[1], 0, "")
			}
		}
	}

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (p *printService) PrintSeminarReport(seminar *model.Seminar) ([]byte, error) {
	pwd, err := os.Getwd()
	if err != nil {
		logoped.ErrorLog.Println("Error getting pwd: ", err)
		return []byte{}, err
	}
	pdf := fpdf.New("P", "mm", "A4", filepath.Join(pwd, "font"))
	pdf.AddFont("Arimo-Regular", "", "Arimo-Regular.json")
	pdf.AddFont("Arimo-Bold", "", "Arimo-Bold.json")
	pdf.AddFont("Helvetica", "", "helvetica_1251.json")
	latTr := pdf.UnicodeTranslatorFromDescriptor("iso-8859-16")
	cirTr := pdf.UnicodeTranslatorFromDescriptor("cp1251")
	trObj := newTranslationDetails(pdf, "Helvetica", "Arimo-Regular", 11, latTr, cirTr)

	pdf.SetMargins(15.0, marginTop, marginRight)
	pdf.AddPage()

	createSimpleHeader(pdf, cirTr)

	pdf.Ln(5)
	pdf.SetTextColor(47, 83, 150)
	pdf.Text(30, pdf.GetY(), trObj.translate("ИЗВЕШТАЈ О РЕАЛИЗОВАНОМ СЕМИНАРУ УНАПРЕЂЕЊА ЗНАЊА", 13))
	pdf.Line(30, pdf.GetY()+1, 180, pdf.GetY()+1)
	pdf.Ln(12)

	var seminarDay *model.SeminarDay
	if len(seminar.Days) > 0 {
		seminarDay = &seminar.Days[0]
	}

	seminarDay, err = SeminarDayService.GetSeminarDayByID(int(seminarDay.ID))
	if err != nil {
		logoped.ErrorLog.Println("Error getting SeminarDay: ", err)
		return []byte{}, err
	}

	presenceTrue := 0
	presenceFalse := 0
	notPresendet := []model.Client{}
	for _, p := range seminarDay.Presence {
		if p.Presence != nil && *p.Presence {
			presenceTrue++
		} else {
			presenceFalse++
			notPresendet = append(notPresendet, p.Client)
		}
	}

	pdf.Text(15, pdf.GetY(), trObj.translDef(fmt.Sprintf("Семинар унапређења знања реализован %s године са почетком у %s и завршетком у", seminar.Start.Format("02.01.2006."), seminarDay.Date.Format("15:04"))))
	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), trObj.translDef(fmt.Sprintf("%s, на адреси Центра за обуку у Србобрану, Туријски пут 17. Семинару је присуствовало %s", seminarDay.Date.Add(375*time.Minute).Format("15:04"), strconv.Itoa(presenceTrue))))

	sentence := fmt.Sprintf("од најављених %s полазника", strconv.Itoa(len(seminarDay.Presence)))
	if len(notPresendet) == 0 {
		sentence = sentence + "."
	}
	if len(notPresendet) == 1 {
		sentence = sentence + "(није присуствовао " + notPresendet[0].Person.FullName() + " " + *notPresendet[0].JMBG + ")."
	}
	if len(notPresendet) > 1 {
		sentence = sentence + "(нису присуствовали "
		for i, c := range notPresendet {
			sentence = sentence + c.Person.FullName() + " " + *c.JMBG
			if i+1 < len(notPresendet) {
				sentence = sentence + ", "
			}
		}

		sentence = sentence + ")."
	}

	sentenceSplited, _ := splitLine(sentence, 130)
	for _, s := range sentenceSplited {
		pdf.Ln(5)
		pdf.Text(15, pdf.GetY(), trObj.translDef(s))
	}

	lines, _ := splitLine(fmt.Sprintf("Тема семинара била је: %s.", seminarDay.Name), 80)
	for _, line := range lines {
		pdf.Ln(5)
		pdf.Text(15, pdf.GetY(), trObj.translDef(line))

	}
	pdf.Ln(10)

	ch := 6.0
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFillColor(146, 208, 80)
	pdf.CellFormat(45, ch, trObj.translDef("Тема семинара:"), "TLR", 0, "C", true, 0, "")
	pdf.CellFormat(29, ch, trObj.translDef("Датум почетка"), "TLR", 0, "C", true, 0, "")
	pdf.CellFormat(28, ch, trObj.translDef("Време"), "TLR", 0, "C", true, 0, "")
	pdf.CellFormat(48, ch, trObj.translDef("Шифра обуке:"), "TLR", 0, "C", true, 0, "")
	pdf.CellFormat(28, ch, trObj.translDef("Број"), "TLR", 0, "C", true, 0, "")
	pdf.Ln(ch)
	pdf.CellFormat(45, ch, "", "BLR", 0, "C", true, 0, "")
	pdf.CellFormat(29, ch, trObj.translDef("обуке:"), "BLR", 0, "C", true, 0, "")
	pdf.CellFormat(28, ch, trObj.translDef("реализације:"), "BLR", 0, "C", true, 0, "")
	pdf.CellFormat(48, ch, "", "BLR", 0, "C", true, 0, "")
	pdf.CellFormat(28, ch, trObj.translDef("полазника:"), "BLR", 0, "C", true, 0, "")

	pdf.Ln(ch)
	ch = 7.0
	pdf.SetTextColor(47, 83, 150)
	lines, _ = splitLine(seminarDay.Name, 22)
	for i, line := range lines {
		if i == 0 {
			pdf.CellFormat(45, ch, trObj.translate(line, 9), "TLR", 0, "C", false, 0, "")
			pdf.CellFormat(29, ch, seminar.Start.Format("02.01.2006."), "TLR", 0, "C", false, 0, "")
			pdf.CellFormat(28, ch, fmt.Sprintf("%s - %s", seminarDay.Date.Format("15:04"), seminarDay.Date.Add(time.Minute*375).Format("15:04")), "TLR", 0, "C", false, 0, "")
			pdf.CellFormat(48, ch, trObj.translDef(seminar.GetCode()), "TLR", 0, "C", false, 0, "")
			pdf.CellFormat(28, ch, strconv.Itoa(presenceTrue), "TLR", 0, "C", false, 0, "")
		} else {
			borderStr := "LR"
			if i+1 == len(lines) {
				borderStr = "BLR"
			}
			pdf.CellFormat(45, ch, trObj.translate(line, 9), borderStr, 0, "C", false, 0, "")
			pdf.CellFormat(29, ch, "", borderStr, 0, "C", false, 0, "")
			pdf.CellFormat(28, ch, "", borderStr, 0, "C", false, 0, "")
			pdf.CellFormat(48, ch, "", borderStr, 0, "C", false, 0, "")
			pdf.CellFormat(28, ch, "", borderStr, 0, "C", false, 0, "")
		}
		if i+1 < len(lines) {
			pdf.Ln(ch)
		}
	}
	pdf.Ln(15)
	pdf.Text(15, pdf.GetY(), trObj.translDef("Наставу су реализовали следећи предавачи:"))
	pdf.Ln(5)
	h := 27.0
	pdf.Rect(15, pdf.GetY()-2, 30, h, "FD")
	pdf.Rect(45, pdf.GetY()-2, 65, h, "FD")
	pdf.Rect(110, pdf.GetY()-2, 15, h, "FD")
	pdf.Rect(125, pdf.GetY()-2, 15, h, "FD")
	pdf.Rect(140, pdf.GetY()-2, 15, h, "FD")
	pdf.Rect(155, pdf.GetY()-2, 15, h, "FD")
	pdf.Rect(170, pdf.GetY()-2, 15, h, "FD")
	pdf.Rect(185, pdf.GetY()-2, 15, h, "FD")

	pdf.SetTextColor(0, 0, 0)
	pdf.Text(20, pdf.GetY()+7, trObj.translDef("Предавачи:"))
	pdf.Text(65, pdf.GetY()+7, trObj.translDef("Наставни час"))

	pdf.TransformBegin()
	pdf.TransformRotate(90, 119, pdf.GetY()+8)
	pdf.Text(106, pdf.GetY()+4, trObj.translDef("Број"))
	pdf.Text(106, pdf.GetY()+8, trObj.translDef("одржаних"))
	pdf.Text(106, pdf.GetY()+12, trObj.translDef("часова"))
	pdf.TransformEnd()

	pdf.TransformBegin()
	pdf.TransformRotate(90, 136, pdf.GetY()+11)
	pdf.Text(123, pdf.GetY()+4, trObj.translDef("Час држи на"))
	pdf.Text(123, pdf.GetY()+8, trObj.translDef("интересантан"))
	pdf.Text(123, pdf.GetY()+12, trObj.translDef("начин"))
	pdf.TransformEnd()

	pdf.TransformBegin()
	pdf.TransformRotate(90, 149, pdf.GetY()+8)
	pdf.Text(136, pdf.GetY()+4, trObj.translDef("Успешно"))
	pdf.Text(136, pdf.GetY()+8, trObj.translDef("објашњава"))
	pdf.Text(136, pdf.GetY()+12, trObj.translDef("градиво"))
	pdf.TransformEnd()

	pdf.TransformBegin()
	pdf.TransformRotate(90, 166, pdf.GetY()+11)
	pdf.Text(153, pdf.GetY()+4, trObj.translDef("Спреман је да"))
	pdf.Text(153, pdf.GetY()+8, trObj.translDef("одговара на"))
	pdf.Text(153, pdf.GetY()+12, trObj.translDef("питања"))
	pdf.TransformEnd()

	pdf.TransformBegin()
	pdf.TransformRotate(90, 179, pdf.GetY()+8)
	pdf.Text(166, pdf.GetY()+4, trObj.translDef("Користио"))
	pdf.Text(166, pdf.GetY()+8, trObj.translDef("је добре"))
	pdf.Text(166, pdf.GetY()+12, trObj.translDef("примере"))
	pdf.TransformEnd()

	pdf.TransformBegin()
	pdf.TransformRotate(90, 194, pdf.GetY()+8)
	pdf.Text(181, pdf.GetY()+4, trObj.translDef("Просечна"))
	pdf.Text(181, pdf.GetY()+8, trObj.translDef("оцена"))
	pdf.TransformEnd()

	teacherClassMap := map[string][]string{}
	for _, c := range seminarDay.Classes {
		teacher := ""
		if c.Teacher != nil {
			teacher = c.Teacher.Person.FullName()
		}
		classes, ok := teacherClassMap[teacher]
		if ok {
			classes = append(classes, c.Name)
			teacherClassMap[teacher] = classes
		} else {
			teacherClassMap[teacher] = []string{c.Name}
		}
	}

	pdf.Ln(25)
	pdf.SetTextColor(47, 83, 150)
	for k, classes := range teacherClassMap {
		splitedTeacherName := strings.Split(k, " ")
		firstName := ""
		lastName := ""
		if len(splitedTeacherName) > 1 {
			firstName = strings.Split(k, " ")[0]
			lastName = strings.Split(k, " ")[1]
		}

		lines, _ := splitLine(classes[0], 30)
		if len(classes) == 1 && len(lines) == 1 {
			pdf.CellFormat(30, ch, trObj.translDef(firstName), "TLR", 0, "C", false, 0, "")
			pdf.CellFormat(65, ch, trObj.translDef(lines[0]), "TLR", 0, "L", false, 0, "")
			pdf.CellFormat(15, ch, strconv.Itoa(len(classes)), "TLR", 0, "C", false, 0, "")
			pdf.CellFormat(15, ch, "", "TLR", 0, "C", false, 0, "")
			pdf.CellFormat(15, ch, "", "TLR", 0, "C", false, 0, "")
			pdf.CellFormat(15, ch, "", "TLR", 0, "C", false, 0, "")
			pdf.CellFormat(15, ch, "", "TLR", 0, "C", false, 0, "")
			pdf.CellFormat(15, ch, "", "TLR", 0, "C", false, 0, "")
			pdf.Ln(ch)

			pdf.CellFormat(30, ch, trObj.translDef(lastName), "BLR", 0, "C", false, 0, "")
			pdf.CellFormat(65, ch, "", "BLR", 0, "C", false, 0, "")
			pdf.CellFormat(15, ch, "", "BLR", 0, "C", false, 0, "")
			pdf.CellFormat(15, ch, "", "BLR", 0, "C", false, 0, "")
			pdf.CellFormat(15, ch, "", "BLR", 0, "C", false, 0, "")
			pdf.CellFormat(15, ch, "", "BLR", 0, "C", false, 0, "")
			pdf.CellFormat(15, ch, "", "BLR", 0, "C", false, 0, "")
			pdf.CellFormat(15, ch, "", "BLR", 0, "C", false, 0, "")
			pdf.Ln(ch)

			continue
		}

		nameDone := false
		lastNameDone := false
		borderStr := "1"
		borderStrName := "L"
		borderStrMarks := "LR"
		classesNum := ""

		for j, class := range classes {
			lines, _ = splitLine(class, 30)
			for i, line := range lines {
				name := ""
				if !nameDone {
					name = firstName
					nameDone = true
					borderStr = "TLR"
				} else if !lastNameDone {
					name = lastName
					lastNameDone = true
					borderStr = "LR"
				} else {
					borderStr = "LR"
				}

				if i == 0 {
					borderStr = "TLR"
				}
				if i == 0 && j == 0 {
					classesNum = strconv.Itoa(len(classes))
				} else {
					classesNum = ""
				}

				if i+1 == len(lines) {
					borderStr = "BLR"
				}
				if j+1 == len(classes) && i+1 == len(lines) {
					borderStrName = "BL"
					borderStrMarks = "BLR"
				}

				pdf.CellFormat(30, ch, trObj.translDef(name), borderStrName, 0, "C", false, 0, "")
				pdf.CellFormat(65, ch, trObj.translDef(line), borderStr, 0, "L", false, 0, "")
				pdf.CellFormat(15, ch, classesNum, borderStrMarks, 0, "C", false, 0, "")
				pdf.CellFormat(15, ch, "", borderStrMarks, 0, "C", false, 0, "")
				pdf.CellFormat(15, ch, "", borderStrMarks, 0, "C", false, 0, "")
				pdf.CellFormat(15, ch, "", borderStrMarks, 0, "C", false, 0, "")
				pdf.CellFormat(15, ch, "", borderStrMarks, 0, "C", false, 0, "")
				pdf.CellFormat(15, ch, "", borderStrMarks, 0, "C", false, 0, "")
				pdf.Ln(ch)
			}
		}

	}

	pdf.AddPage()
	ch = 10.0
	pdf.Text(15, pdf.GetY(), trObj.translDef("Списак полазника:"))
	pdf.Ln(3)

	pdf.Rect(15, pdf.GetY(), 186, 5, "FD")
	pdf.SetTextColor(0, 0, 0)
	pdf.Text(90, pdf.GetY()+4, trObj.translDef("Списак полазника"))
	pdf.SetTextColor(47, 83, 150)
	pdf.Ln(5)

	sort.Slice(seminar.Trainees, func(i, j int) bool {
		return *seminar.Trainees[i].Client.JMBG < *seminar.Trainees[j].Client.JMBG
	})

	for i := 0; i < len(seminar.Trainees); i++ {
		pdf.CellFormat(10, ch, strconv.Itoa(i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(47, ch, trObj.translDef(seminar.Trainees[i].Client.Person.FullName()), "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, ch, *seminar.Trainees[i].Client.JMBG, "1", 0, "C", false, 0, "")
		pdf.CellFormat(10, ch, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(42, ch, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(42, ch, "", "1", 0, "C", false, 0, "")
		pdf.Ln(ch)
	}

	pdf.AddPage()
	pdf.Text(15, pdf.GetY(), trObj.translDef("Евалуација наставе:"))
	pdf.Ln(3)

	pdf.Rect(15, pdf.GetY(), 186, 5, "FD")
	pdf.SetTextColor(0, 0, 0)
	pdf.Text(90, pdf.GetY()+4, trObj.translDef("Евалуација семинара"))
	pdf.SetTextColor(47, 83, 150)
	pdf.Ln(5)

	ch = 5
	pdf.CellFormat(80, ch, trObj.translDef("Простор је био пријатан за рад(осветљење,"), "TLR", 0, "L", false, 0, "")
	pdf.CellFormat(13, ch, "", "TLR", 0, "L", false, 0, "")
	pdf.CellFormat(80, ch, trObj.translDef("Предавачи су се придржавали сатнице"), "TLR", 0, "L", false, 0, "")
	pdf.CellFormat(13, ch, "", "TLR", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(80, ch, trObj.translDef("температура, столица, акустика)"), "BLR", 0, "L", false, 0, "")
	pdf.CellFormat(13, ch, "", "BLR", 0, "L", false, 0, "")
	pdf.CellFormat(80, ch, trObj.translDef("(почетка/завршетка часа)"), "BLR", 0, "L", false, 0, "")
	pdf.CellFormat(13, ch, "", "BLR", 0, "L", false, 0, "")
	pdf.Ln(5)

	pdf.CellFormat(80, ch, trObj.translDef("Атмосфера током семинара је била"), "1", 0, "L", false, 0, "")
	pdf.CellFormat(13, ch, "", "1", 0, "L", false, 0, "")
	pdf.CellFormat(80, ch, trObj.translDef("Трајање пауза је било одговарајуће"), "1", 0, "L", false, 0, "")
	pdf.CellFormat(13, ch, "", "1", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(80, ch, trObj.translDef("Храна током обуке је била одговарајућа"), "1", 0, "L", false, 0, "")
	pdf.CellFormat(13, ch, "", "1", 0, "L", false, 0, "")
	pdf.CellFormat(80, ch, trObj.translDef("Пиће током обуке је било одговарајуће"), "1", 0, "L", false, 0, "")
	pdf.CellFormat(13, ch, "", "1", 0, "L", false, 0, "")
	pdf.Ln(5)

	pdf.CellFormat(80, ch, trObj.translDef("Предавачи су подржавали комуникацију и"), "TLR", 0, "L", false, 0, "")
	pdf.CellFormat(13, ch, "", "TLR", 0, "L", false, 0, "")
	pdf.CellFormat(80, ch, trObj.translDef("У којој мери је обука испунила Ваша"), "TLR", 0, "L", false, 0, "")
	pdf.CellFormat(13, ch, "", "TLR", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(80, ch, trObj.translDef("интеракцију полазника"), "BLR", 0, "L", false, 0, "")
	pdf.CellFormat(13, ch, "", "BLR", 0, "L", false, 0, "")
	pdf.CellFormat(80, ch, trObj.translDef("очекивања"), "BLR", 0, "L", false, 0, "")
	pdf.CellFormat(13, ch, "", "BLR", 0, "L", false, 0, "")

	pdf.Ln(15)
	pdf.Text(15, pdf.GetY(), trObj.translDef("За време семинара унапређења знања спроведено је и улазно и излазно тестирање полазника."))
	pdf.Ln(4)

	clientTests, err := TestService.GetClientTestBySeminarDayID(int(seminarDay.ID))
	if err != nil {
		logoped.ErrorLog.Println("Error getting client test by seminar day id: ", err)
		return []byte{}, err
	}

	testsMap := map[string][]model.ClientTest{}
	for _, ct := range clientTests {
		_, ok := testsMap[ct.Jmbg]
		if ok {
			testsMap[ct.Jmbg] = append(testsMap[ct.Jmbg], ct)
		} else {
			testsMap[ct.Jmbg] = []model.ClientTest{ct}
		}
	}

	sum1 := 0.0
	num1 := 0.0
	sum2 := 0.0
	num2 := 0.0
	for _, t := range testsMap {
		if len(t) != 2 {
			continue
		}

		if t[0].CreatedAt.Before(t[1].CreatedAt) {
			num1++
			sum1 = sum1 + t[0].Result

			num2++
			sum2 = sum2 + t[1].Result
		} else {
			num2++
			sum2 = sum2 + t[0].Result

			num1++
			sum1 = sum1 + t[1].Result
		}
	}

	percentIn := 0.0
	if sum1 > 0 {
		percentIn = (sum1 / num1) * 100
	}
	percentOut := 0.0
	if sum2 > 0 {
		percentIn = (sum2 / num2) * 100
	}

	pdf.Text(15, pdf.GetY(), trObj.translDef(fmt.Sprintf("Просечан резултат на улазном тесту био је %d%% тачних одговора, док је на излазном тесту", int(percentIn))))
	pdf.Ln(4)
	pdf.Text(15, pdf.GetY(), trObj.translDef(fmt.Sprintf("резултат био %d%% тачних одговора. Проценат успешности приказан је у наредној табели:", int(percentOut))))
	pdf.Ln(8)

	pdf.Rect(15, pdf.GetY(), 186, 5, "FD")
	pdf.SetTextColor(0, 0, 0)
	pdf.Text(90, pdf.GetY()+4, trObj.translDef("Евалуација семинара"))
	pdf.SetTextColor(47, 83, 150)
	pdf.Ln(5)

	for i := 0; i < len(seminar.Trainees); i++ {
		p := ""
		p = fmt.Sprintf("Полазник %d", i+1)

		tests, _ := testsMap[*seminar.Trainees[i].Client.JMBG]

		t1 := 0
		t2 := 0
		d := 0
		if len(tests) > 1 {
			if tests[0].CreatedAt.Before(tests[1].CreatedAt) {
				t1 = int(tests[0].Result * 100)
				t2 = int(tests[1].Result * 100)
				d = int(tests[1].Result*100 - tests[0].Result*100)
			} else {
				t1 = int(tests[1].Result * 100)
				t2 = int(tests[0].Result * 100)
				d = int(tests[0].Result*100 - tests[1].Result*100)
			}
		}

		pdf.CellFormat(26, ch, trObj.translDef(p), "1", 0, "L", false, 0, "")
		pdf.CellFormat(12, ch, fmt.Sprintf("%d%%", t1), "1", 0, "L", false, 0, "")
		pdf.CellFormat(12, ch, fmt.Sprintf("%d%%", t2), "1", 0, "L", false, 0, "")
		pdf.CellFormat(12, ch, fmt.Sprintf("%d%%", d), "1", 0, "L", false, 0, "")

		if (i+1)%3 == 0 {
			pdf.Ln(5)
		}
	}

	pdf.Ln(7)
	pdf.Text(15, pdf.GetY()+4, trObj.translDef("Семинар унапређења знања успешно је реализован уз активно учешће полазника."))

	pdf.Ln(15)
	pdf.Text(175, pdf.GetY(), trObj.translDef("Србобран,"))
	pdf.Ln(4)
	pdf.Text(160, pdf.GetY(), trObj.translDef(seminar.Start.Format("02.01.2006.")+" године"))

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (p *printService) PrintSeminarReport2(seminar *model.Seminar) ([]byte, error) {
	pwd, err := os.Getwd()
	if err != nil {
		logoped.ErrorLog.Println("Error getting pwd: ", err)
		return []byte{}, err
	}
	pdf := fpdf.New("P", "mm", "A4", filepath.Join(pwd, "font"))
	pdf.AddFont("Arimo-Regular", "", "Arimo-Regular.json")
	pdf.AddFont("Arimo-Bold", "", "Arimo-Bold.json")
	pdf.AddFont("Helvetica", "", "helvetica_1251.json")
	latTr := pdf.UnicodeTranslatorFromDescriptor("iso-8859-16")
	cirTr := pdf.UnicodeTranslatorFromDescriptor("cp1251")
	trObj := newTranslationDetails(pdf, "Helvetica", "Arimo-Regular", 11, latTr, cirTr)

	pdf.SetMargins(15.0, marginTop, marginRight)
	pdf.AddPage()

	pdf.Ln(20)
	pdf.SetTextColor(47, 83, 150)
	pdf.Text(85, pdf.GetY(), trObj.translate("Срболаб доо", 15))
	pdf.Ln(12)

	pdf.Text(55, pdf.GetY(), trObj.translate("Центар за едукацију и развој Срболаб", 17))
	pdf.Ln(12)

	pdf.Text(40, pdf.GetY(), trObj.translate("Извештај о реализованој периодичној обуци", 19))
	pdf.Ln(12)

	pdf.Text(65, pdf.GetY(), trObj.translate(fmt.Sprintf("Центар за обуку у %s", seminar.ClassRoom.Location.GetLocationForSentence()), 15))
	pdf.Ln(20)

	pdf.Text(20, pdf.GetY(), trObj.translate("1. Основни подаци о Центру и времену одржавања", 13))
	pdf.Ln(5)

	ch := 7.0
	pdf.SetFillColor(232, 238, 248)
	pdf.CellFormat(60, ch, trObj.translDef("Mesto održavanja obuka"), "1", 0, "C", true, 0, "")
	pdf.SetFillColor(255, 255, 255)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(115, ch, trObj.translDef(seminar.ClassRoom.Location.Address.GetStreetWithNumber()+", "+seminar.ClassRoom.Location.Address.Place), "1", 0, "L", true, 0, "")
	pdf.Ln(ch)
	pdf.SetFillColor(232, 238, 248)
	pdf.SetTextColor(47, 83, 150)
	pdf.CellFormat(60, ch, trObj.translDef("Vrste obuka koje su realizovane"), "1", 0, "C", true, 0, "")
	pdf.SetFillColor(255, 255, 255)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(115, ch, trObj.translDef("Периодична ЦПЦ обука - 7 часова"), "1", 0, "L", true, 0, "")
	pdf.Ln(ch)
	pdf.SetFillColor(232, 238, 248)
	pdf.SetTextColor(47, 83, 150)
	pdf.CellFormat(60, ch, trObj.translDef("Podaci o realizaciji obuke"), "TLR", 0, "C", true, 0, "")
	pdf.SetFillColor(255, 255, 255)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(115, ch, trObj.translDef(seminar.Start.Format("02.01.2006.")+" године"), "1", 0, "L", true, 0, "")
	pdf.Ln(ch)
	pdf.SetFillColor(232, 238, 248)
	pdf.SetTextColor(47, 83, 150)
	pdf.CellFormat(60, ch, trObj.translDef(""), "BLR", 0, "C", true, 0, "")
	pdf.SetFillColor(255, 255, 255)
	pdf.SetTextColor(0, 0, 0)
	start := ""
	end := ""
	if len(seminar.Days) > 0 {
		start = seminar.Days[0].Date.Format("15:04")
		end = seminar.Days[0].Date.Add(375 * time.Minute).Format("15:04")
	}

	pdf.CellFormat(115, ch, trObj.translDef(fmt.Sprintf("%s до %s", start, end)), "1", 0, "L", true, 0, "")
	pdf.Ln(20)

	pdf.SetTextColor(47, 83, 150)
	pdf.Text(20, pdf.GetY(), trObj.translate(fmt.Sprintf("2.  Полазници (%s)", seminar.GetCode()), 13))
	pdf.Ln(5)

	pdf.SetFillColor(232, 238, 248)
	pdf.Rect(15, pdf.GetY(), 180, 6, "FD")
	pdf.SetTextColor(0, 0, 0)
	pdf.Text(64, pdf.GetY()+5, trObj.translate("Периодична ЦПЦ обука - 7 часова", 13))
	pdf.Ln(6)
	pdf.SetTextColor(47, 83, 150)
	pdf.Rect(15, pdf.GetY(), 25, 18, "FD")
	pdf.Text(21, pdf.GetY()+10, trObj.translate("Датум", 13))
	pdf.Rect(40, pdf.GetY(), 25, 18, "FD")
	pdf.Text(44, pdf.GetY()+5, trObj.translate("Укупан", 13))
	pdf.Text(47, pdf.GetY()+10, trObj.translate("број", 13))
	pdf.Text(41, pdf.GetY()+15, trObj.translate("полазника", 13))
	pdf.Rect(65, pdf.GetY(), 40, 18, "FD")
	pdf.Text(68, pdf.GetY()+7, trObj.translate("Правно/физичко", 13))
	pdf.Text(80, pdf.GetY()+14, trObj.translate("лице", 13))
	pdf.Rect(105, pdf.GetY(), 30, 18, "FD")
	pdf.Text(115, pdf.GetY()+7, trObj.translate("Број", 13))
	pdf.Text(108, pdf.GetY()+14, trObj.translate("полазника", 13))
	pdf.Rect(135, pdf.GetY(), 30, 18, "FD")
	pdf.Text(141, pdf.GetY()+10, trObj.translate("Попуст", 13))
	pdf.Rect(165, pdf.GetY(), 30, 18, "FD")
	pdf.Text(171.5, pdf.GetY()+7, trObj.translate("Начин", 13))
	pdf.Text(169, pdf.GetY()+14, trObj.translate("плаћања", 13))

	pdf.Ln(18)

	totalNum := len(seminar.Trainees)
	companyMap := map[string]int{}
	for _, t := range seminar.Trainees {
		companyMap[t.Client.Company.Name]++
	}

	pdf.SetFont("Helvetica", "", 10)
	pdf.SetFillColor(255, 255, 255)
	firstLine := true
	borderGlobal := "LR"
	mapRowNum := 0
	for k, v := range companyMap {
		mapRowNum++
		if k == "" {
			k = "Физичко лице"
		}

		lines, _ := splitLine(k, 19)
		for i, line := range lines {
			date := ""
			total := ""
			if firstLine {
				date = trObj.translate(seminar.Start.Format("02.01.2006."), 10)
				total = strconv.Itoa(totalNum)
				firstLine = false
			}

			numByCompany := ""
			border := "LR"
			if i == 0 {
				numByCompany = strconv.Itoa(v)
				border = "TLR"
			} else {
				ch = 5
			}

			if mapRowNum == len(companyMap) && i+1 == len(lines) {
				if len(lines) == 1 {
					border = "1"
				} else {
					border = "BLR"
					borderGlobal = "BLR"
				}
			}

			pdf.CellFormat(25, ch, date, borderGlobal, 0, "C", true, 0, "")
			pdf.CellFormat(25, ch, total, borderGlobal, 0, "C", true, 0, "")
			pdf.CellFormat(40, ch, trObj.translate(line, 10), border, 0, "L", true, 0, "")
			pdf.CellFormat(30, ch, numByCompany, border, 0, "C", true, 0, "")
			pdf.CellFormat(30, ch, "", border, 0, "C", true, 0, "")
			pdf.CellFormat(30, ch, "", border, 0, "C", true, 0, "")
			pdf.Ln(ch)

			ch = 7
		}
	}

	pdf.SetFillColor(232, 238, 248)
	pdf.SetTextColor(47, 83, 150)
	pdf.Rect(15, pdf.GetY(), 25, ch, "FD")
	pdf.Text(21, pdf.GetY()+5, trObj.translate("Укупно", 12))
	pdf.Rect(40, pdf.GetY(), 25, ch, "FD")
	pdf.Text(44, pdf.GetY()+5, strconv.Itoa(totalNum))
	pdf.Rect(65, pdf.GetY(), 130, ch, "FD")
	pdf.Text(74, pdf.GetY()+5, trObj.translate("Напомена:", 12))

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func getClassTime(d time.Time, i int) string {
	switch i {
	case 1:
		return d.Format("15:04") + " - " + d.Add(45*time.Minute).Format("15:04")
	case 2:
		return d.Add(55*time.Minute).Format("15:04") + " - " + d.Add(100*time.Minute).Format("15:04")
	case 3:
		return d.Add(105*time.Minute).Format("15:04") + " - " + d.Add(150*time.Minute).Format("15:04")
	case 4:
		return d.Add(175*time.Minute).Format("15:04") + " - " + d.Add(220*time.Minute).Format("15:04")
	case 5:
		return d.Add(225*time.Minute).Format("15:04") + " - " + d.Add(270*time.Minute).Format("15:04")
	case 6:
		return d.Add(280*time.Minute).Format("15:04") + " - " + d.Add(325*time.Minute).Format("15:04")
	case 7:
		return d.Add(330*time.Minute).Format("15:04") + " - " + d.Add(375*time.Minute).Format("15:04")
	}
	return ""
}

func splitLine(text string, length int) ([]string, float64) {
	words := strings.Split(text, " ")
	lineNum := 1.0

	lines := []string{}
	currentLine := ""
	for _, word := range words {
		wordR := []rune(word)
		currentLineR := []rune(currentLine)
		if len(currentLineR)+len(wordR)+1 > length {
			lines = append(lines, currentLine)
			lineNum++
			currentLine = ""
		}
		currentLine = currentLine + " " + word
	}

	return append(lines, currentLine), lineNum
}
