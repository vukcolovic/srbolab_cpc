package service

import (
	"bytes"
	"fmt"
	"github.com/go-pdf/fpdf"
	"github.com/skip2/go-qrcode"
	"os"
	"path/filepath"
	"srbolab_cpc/logoped"
	"srbolab_cpc/model"
	"srbolab_cpc/util"
	"strconv"
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
	PrintConfirmationStatements(seminar *model.Seminar) ([]byte, error)
	PrintConfirmations(seminar *model.Seminar) ([]byte, error)
	PrintConfirmationReceives(seminar *model.Seminar) ([]byte, error)
	PrintMuster(day *model.SeminarDay) ([]byte, error)
	PrintCheckIn(seminar *model.Seminar) ([]byte, error)
	PrintSeminarEvidence(day *model.SeminarDay) ([]byte, error)
	PrintTestBarcode() ([]byte, error)
	PrintPlanTreningRealization(day *model.SeminarDay) ([]byte, error)
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
	trObj := newTranslationDetails(pdf, "Helvetica", "Arimo-Regular", 8, latTr, cirTr)

	pdf.SetMargins(marginLeft, marginTop, marginRight)
	pdf.AddPage()

	pdf.SetFont("Helvetica", "", 8)
	createSimpleHeader(pdf, cirTr)

	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), trObj.translDef("Број документа: "))
	pdf.Text(35, pdf.GetY(), "")
	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), trObj.translDef(("Шифра обуке: ")))
	pdf.Text(35, pdf.GetY(), trObj.translDef(seminar.GetCode()))
	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), trObj.translDef("Место: "))
	pdf.Text(27, pdf.GetY(), trObj.translDef(seminar.ClassRoom.Location.Address.Place))
	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), trObj.translDef("Регистрациони лист - списак полазника за "+seminar.SeminarTheme.Name))
	pdf.Text(140, pdf.GetY(), trObj.translDef("Datum: "+time.Now().Format("01.02.2006")))

	ch := 8.0
	pdf.Ln(ch)
	pdf.CellFormat(20, ch, trObj.translDef("Редни број"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(75, ch, trObj.translDef("Име и презиме"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(35, ch, trObj.translDef("ЈМБГ"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(55, ch, trObj.translDef("Фирма у којој сте запослени"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(35, ch, trObj.translDef("Телефон"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(60, ch, trObj.translDef("Потпис"), "1", 0, "C", false, 0, "")

	for i, cs := range seminar.Trainees {
		pdf.Ln(ch)
		pdf.CellFormat(20, ch, strconv.Itoa(i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(75, ch, trObj.translDef(cs.Client.Person.FullName()), "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, ch, *cs.Client.JMBG, "1", 0, "C", false, 0, "")
		pdf.CellFormat(55, ch, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, ch, "", "1", 0, "C", false, 0, "")
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
	pdf.AddFont("Helvetica", "", "helvetica_1251.json")
	latTr := pdf.UnicodeTranslatorFromDescriptor("iso-8859-16")
	cirTr := pdf.UnicodeTranslatorFromDescriptor("cp1251")
	trObj := newTranslationDetails(pdf, "Helvetica", "Arimo-Regular", 11, latTr, cirTr)

	pdf.SetMargins(marginLeft, 40, marginRight)

	notPassedClientIds := make(map[uint]string)
	for _, day := range seminar.Days {
		for _, p := range day.Presence {
			if !*p.Presence && !day.Date.After(time.Now()) {
				notPassedClientIds[p.ClientID] = ""
			}
		}
	}

	for _, client := range seminar.Trainees {
		if _, exist := notPassedClientIds[client.ClientID]; exist {
			continue
		}
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
		pdf.Text(10, pdf.GetY(), trObj.translate("законским и подзаконсим актима и да исто неће бити од утицаја на прижање услуга од стране руковаоца.", 10))
		pdf.Ln(25)

		pdf.Text(15, pdf.GetY(), trObj.translDef("Датум: "))
		pdf.Line(30, pdf.GetY(), 70, pdf.GetY())
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(35, pdf.GetY()-1, time.Now().Format("02.01.2006"))
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

	for _, client := range seminar.Trainees {
		if _, exist := notPassedClientIds[client.ClientID]; exist {
			continue
		}
		pdf.AddPage()

		pdf.SetFont("Arimo-Regular", "", 11)

		pdf.Ln(5)
		pdf.Text(15, pdf.GetY(), trObj.translDef("Број:"))
		pdf.Text(30, pdf.GetY(), trObj.translDef(seminar.GetCode()))
		pdf.Ln(5)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(15, pdf.GetY(), trObj.translDef("Дана:"))
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(30, pdf.GetY(), time.Now().Format("02.01.2006."))
		pdf.Ln(5)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(15, pdf.GetY(), trObj.translDef("Место:"))
		pdf.SetFont("Arimo-Bold", "", 11)
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
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.CellFormat(wr, ch-1, trObj.translDef(client.Client.Person.FirstName+" ("+client.Client.Person.MiddleName+") "+client.Client.Person.LastName), "LRT", 0, "L", false, 0, "")
		pdf.Ln(ch - 1)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.CellFormat(wl, ch-1, trObj.translDef("родитеља, презиме"), "LRB", 0, "L", false, 0, "")
		pdf.CellFormat(wr, ch-1, "", "LRB", 0, "L", false, 0, "")
		pdf.Ln(ch - 1)
		pdf.CellFormat(wl, ch, trObj.translDef("ЈМБГ"), "1", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", 11)
		for _, l := range trObj.translDef(*client.Client.JMBG) {
			pdf.CellFormat(wr/13, ch, string(l), "1", 0, "C", false, 0, "")
		}
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.CellFormat(wl, ch, trObj.translDef("Место прбивалишта"), "1", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.CellFormat(wr, ch, trObj.translDef(client.Client.Address.Place), "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.CellFormat(wl, ch, trObj.translDef("Адреса прбивалишта"), "1", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.CellFormat(wr, ch, trObj.translDef(client.Client.Address.Street+" "+client.Client.Address.HouseNumber), "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.CellFormat(wl, ch, trObj.translDef("Редни број семинара"), "1", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", 11)
		//fixme just closed seminars and passed
		completedSeminars := 0
		if client.Client.InitialCompletedSeminars != nil {
			completedSeminars = *client.Client.InitialCompletedSeminars
		}

		//changed after consulatition with Marko Jovanovic
		//completedInSrbolab := 0
		//for _, s := range client.Client.Seminars {
		//	if s.Pass != nil && *s.Pass {
		//		completedInSrbolab++
		//	}
		//}
		seminarNumber := completedSeminars
		if seminarNumber > 5 {
			seminarNumber = 5
		}
		cx := 89.5 + float64(seminarNumber)*7
		if seminarNumber == 1 {
			cx = cx + 1
		}
		if seminarNumber == 5 {
			cx = cx + 1
		}
		if seminarNumber > 0 {
			pdf.Circle(cx, 137, 3, "")
		}
		pdf.CellFormat(wr, ch, " I    II    III    IV    V", "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.CellFormat(wl, ch-1, trObj.translDef("Датум похађања"), "LRT", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.CellFormat(wr, ch-1, "od "+startSeminar.Format("02.01.2006"), "LRT", 0, "L", false, 0, "")
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
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.CellFormat(wr, ch-1, "do "+endSeminar.Format("02.01.2006"), "LRB", 0, "L", false, 0, "")
		pdf.Ln(ch - 1)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.CellFormat(wl, ch-1, trObj.translDef("Место похађања"), "LRT", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.CellFormat(wr, ch-1, trObj.translDef(seminar.ClassRoom.Location.Address.Place+", "+seminar.ClassRoom.Location.Address.Street+" "+seminar.ClassRoom.Location.Address.HouseNumber), "LRT", 0, "L", false, 0, "")
		pdf.Ln(ch - 1)
		pdf.SetFont("Arimo-Regular", "", 11)

		pdf.CellFormat(wl, ch-1, trObj.translDef(fmt.Sprintf("%s обуке", seminarType)), "LRB", 0, "L", false, 0, "")
		pdf.CellFormat(wr, ch-1, "", "LRB", 0, "L", false, 0, "")
		pdf.Ln(ch - 1)
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.CellFormat(wl, ch, trObj.translDef("Врста ЦПЦ"), "1", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", 11)
		if client.Client.CLicence != nil && *client.Client.CLicence {
			pdf.Circle(97.5, 178.5, 2.5, "")
		}
		if client.Client.DLicence != nil && *client.Client.DLicence {
			pdf.Circle(128.0, 178.5, 2.5, "")
		}
		pdf.CellFormat(wr, ch, trObj.translDef("1. превоз терета  2. превоз путника"), "1", 0, "L", false, 0, "")
		pdf.Ln(20)

		pdf.SetFont("Arimo-Bold", "", 11)
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

		pdf.SetFont("Arimo-Bold", "", 12)
		pdf.Text(40, pdf.GetY(), trObj.translate(fmt.Sprintf("Изјава о преузимању потврде и завршеној %s обуци", seminar.SeminarTheme.BaseSeminarType.GetSeminarTypeForSentence()), 12))
		pdf.Ln(5)
		pdf.Text(60, pdf.GetY(), trObj.translate("на обавезним семинарима унапређења знања", 12))
		pdf.Ln(20)

		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(20, pdf.GetY(), trObj.translDef("Дана"))
		pdf.Line(33, pdf.GetY(), 60, pdf.GetY())
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(33, pdf.GetY()-1, latTr(time.Now().Format("02.01.2006.")))
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(62, pdf.GetY(), trObj.translDef("године, "))
		pdf.Line(77, pdf.GetY(), 135, pdf.GetY())
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(80, pdf.GetY()-1, trObj.translDef(client.Client.Person.FullName()))
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(137, pdf.GetY(), trObj.translDef("(име и презиме), ЈМБГ"))
		pdf.Ln(6)
		pdf.Line(20, pdf.GetY(), 55, pdf.GetY())
		pdf.SetFont("Arimo-Bold", "", 11)
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
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Line(31, pdf.GetY(), 58, pdf.GetY())
		pdf.Text(32, pdf.GetY()-1, time.Now().Format("02.01.2006"))
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(60, pdf.GetY(), trObj.translDef("године."))
	}

	for company, clients := range companyClientsMap {
		pdf.AddPage()

		createSimpleHeader(pdf, cirTr)

		pdf.SetFont("Arimo-Bold", "", 12)
		pdf.Text(40, pdf.GetY(), trObj.translate(fmt.Sprintf("Изјава о преузимању потврде и завршеној %s обуци", seminar.SeminarTheme.BaseSeminarType.GetSeminarTypeForSentence()), 12))
		pdf.Ln(5)
		pdf.Text(60, pdf.GetY(), trObj.translate("на обавезним семинарима унапређења знања", 12))
		pdf.Ln(20)

		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(15, pdf.GetY(), trObj.translDef("Дана"))
		pdf.Line(27, pdf.GetY(), 57, pdf.GetY())
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(30, pdf.GetY()-1, trObj.translDef(time.Now().Format("02.01.2006.")))
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(60, pdf.GetY(), trObj.translDef("године, "))
		pdf.Line(75, pdf.GetY(), 135, pdf.GetY())
		pdf.Text(135, pdf.GetY(), trObj.translDef("(име и презиме), запослен у"))
		pdf.Ln(6)
		pdf.Text(15, pdf.GetY(), trObj.translDef("у фирми"))
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Text(30, pdf.GetY()-1, trObj.translDef(company))
		pdf.Line(28, pdf.GetY(), 178, pdf.GetY())
		pdf.SetFont("Arimo-Regular", "", 11)
		pdf.Text(180, pdf.GetY(), ",")
		pdf.Ln(6)
		pdf.Text(15, pdf.GetY(), trObj.translDef(("ЈМБГ")))
		pdf.Line(28, pdf.GetY(), 70, pdf.GetY())
		pdf.Text(70, pdf.GetY(), trObj.translDef(fmt.Sprintf(", је преузеп потврде о завршеној %s обуци на", seminar.SeminarTheme.BaseSeminarType.GetSeminarTypeForSentence())))
		pdf.Ln(6)
		pdf.Text(15, pdf.GetY(), trObj.translDef("обавезним семинарима унапређења знања за следећа лица:"))
		pdf.Ln(10)

		ch := 5.0
		pdf.CellFormat(10, ch, trObj.translDef("РБ"), "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, ch, trObj.translDef("Име"), "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, ch, trObj.translDef("Презиме"), "1", 0, "L", false, 0, "")
		pdf.CellFormat(35, ch, trObj.translDef("ЈМБГ"), "1", 0, "L", false, 0, "")
		pdf.CellFormat(50, ch, trObj.translDef("Број потврде"), "1", 0, "L", false, 0, "")

		for i, client := range clients {
			pdf.Ln(ch)
			pdf.CellFormat(10, ch, strconv.Itoa(i+1), "1", 0, "L", false, 0, "")
			pdf.CellFormat(40, ch, trObj.translDef(client.Client.Person.FirstName), "1", 0, "L", false, 0, "")
			pdf.CellFormat(40, ch, trObj.translDef(client.Client.Person.LastName), "1", 0, "L", false, 0, "")
			pdf.CellFormat(35, ch, *client.Client.JMBG, "1", 0, "L", false, 0, "")
			pdf.CellFormat(50, ch, seminar.GetCode()+"/"+strconv.Itoa(int(client.ClientID)), "1", 0, "L", false, 0, "")
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
		pdf.SetFont("Arimo-Bold", "", 11)
		pdf.Line(26, pdf.GetY(), 48, pdf.GetY())
		pdf.Text(27, pdf.GetY()-1, time.Now().Format("02.01.2006"))
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
	trObj := newTranslationDetails(pdf, "Helvetica", "Arimo-Regular", 8, latTr, cirTr)

	pdf.SetMargins(marginLeft, marginTop, marginRight)
	pdf.AddPage()

	pdf.SetFont("Arimo-Regular", "", 8)
	createSimpleHeader(pdf, cirTr)

	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), trObj.translDef("Mesto: "))
	pdf.Text(27, pdf.GetY(), trObj.translDef(day.Seminar.ClassRoom.Location.Address.Place))
	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), trObj.translDef("Šifra obuke: "))
	pdf.Text(35, pdf.GetY(), day.Seminar.GetCode())
	pdf.Ln(5)
	pdf.Text(15, pdf.GetY(), trObj.translDef("Datum: "))
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

	for i, cs := range day.Presence {
		pdf.Ln(ch)
		pdf.CellFormat(10, ch, strconv.Itoa(i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(45, ch, "", "1", 0, "C", false, 0, "")
		pdf.Text(25, 71+float64(i*14), trObj.translDef(cs.Client.Person.FullName()))
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

		pdf.SetFont("Arimo-Bold", "", 15)
		pdf.Text(85, pdf.GetY(), trObj.translate("П Р И Ј А В А*", 15))

		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.Ln(10)
		pdf.Text(70, pdf.GetY(), trObj.translDef("за похађање обавезног семинара"))
		pdf.Ln(5)
		pdf.Text(87, pdf.GetY(), trObj.translDef("унапређења знања"))
		pdf.Ln(10)

		pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.Text(15, pdf.GetY(), trObj.translDef("ЛИЧНИ ПОДАЦИ"))
		pdf.Ln(1)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("Име (име једног родитеља) презиме:"), "1", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.CellFormat(110, ch, trObj.translDef(client.Client.Person.FullNameWithMiddleName()), "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("ЈМБГ:"), "1", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.CellFormat(110, ch, *client.Client.JMBG, "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("Датум рођења:"), "1", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.CellFormat(110, ch, client.Client.GetBirthDate(), "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("Место рођења, држава:"), "1", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.CellFormat(110, ch, trObj.translDef(*client.Client.PlaceBirth+", "+*client.Client.CountryBirth), "1", 0, "L", false, 0, "")

		pdf.Ln(17)
		pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.Text(15, pdf.GetY(), trObj.translDef("ПОДАЦИ О ПРЕБИВАЛИШТУ/БОРАВИШТУ"))
		pdf.Ln(1)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("Место:"), "1", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.CellFormat(110, ch, trObj.translDef(client.Client.Address.Place), "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("Поштански број:"), "1", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.CellFormat(110, ch, trObj.translDef(client.Client.Address.PostCode), "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("Улица и кућни број:"), "1", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.CellFormat(110, ch, trObj.translDef(client.Client.Address.GetStreetWithNumber()), "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("Телефон/Мобилни:"), "1", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.CellFormat(110, ch, trObj.translDef(client.Client.Person.PhoneNumber), "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("е-маил:"), "1", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.CellFormat(110, ch, trObj.translDef(client.Client.Person.Email), "1", 0, "L", false, 0, "")

		pdf.Ln(17)
		pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.Text(15, pdf.GetY(), trObj.translDef("ПОДАЦИ О КВАЛИФИКАЦИОНОЈ КАРТИЦИ ВОЗАЧА"))
		pdf.Ln(1)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("Број картице(СРБ број)*:"), "1", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.CellFormat(110, ch, *client.Client.CPCNumber, "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("Серијски број картице(СРБ број):"), "1", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.CellFormat(110, ch, "", "1", 0, "L", false, 0, "")
		pdf.Ln(ch)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.CellFormat(70, ch, trObj.translDef("Рок важења картице:"), "1", 0, "L", false, 0, "")
		pdf.SetFont("Arimo-Bold", "", fontSize)
		cpcDate := ""
		if client.Client.CPCDate != nil {
			cpcDate = client.Client.CPCDate.Format("02.01.2006.")
		}
		pdf.CellFormat(110, ch, cpcDate, "1", 0, "L", false, 0, "")

		pdf.Ln(17)
		pdf.SetFont("Arimo-Bold", "", fontSize)
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
		pdf.Text(10, pdf.GetY(), trObj.translDef("Уз попуњен образац Пријаве за похађање семинара, прилаже се:"))
		pdf.Ln(5)
		pdf.Text(20, pdf.GetY(), trObj.translDef("- докаж о уплати трошкова за похађање семинара, по важећој тарифи"))

		pdf.Ln(20)
		pdf.Text(18, pdf.GetY(), "U ")
		pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.Text(23, pdf.GetY(), trObj.translDef(seminar.ClassRoom.Location.Address.Place))
		pdf.Line(23, pdf.GetY()+1, 65, pdf.GetY()+1)
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.Text(65.5, pdf.GetY(), trObj.translDef(", дана"))
		pdf.SetFont("Arimo-Bold", "", fontSize)
		pdf.Line(78, pdf.GetY()+1, 110, pdf.GetY()+1)
		pdf.Text(80, pdf.GetY(), time.Now().Format("02.01.2006."))
		pdf.SetFont("Arimo-Regular", "", fontSize)
		pdf.Text(110.5, pdf.GetY(), trObj.translDef(", године"))

		pdf.Ln(30)
		pdf.Text(135, pdf.GetY(), trObj.translDef("Потпис подносиоца пријаве: "))
		pdf.Ln(15)
		pdf.Line(135, pdf.GetY(), 190, pdf.GetY())
		pdf.Ln(10)

		pdf.SetFont("Arimo-Regular", "", 9)
		pdf.Text(15, pdf.GetY(), trObj.translDef("*Образац пријаве попунити читко штампаним словима"))
		pdf.Ln(5)
		pdf.Text(15, pdf.GetY(), trObj.translDef("*Уписати број картице (СРБ број) или број „Потврде о пријему захтева за издавање сертификата о стручној"))
		pdf.Ln(4)
		pdf.Text(17, pdf.GetY(), trObj.translDef("компетенцији и квалификационе картице возача, уколико сте покренули поступак издавања"))
		pdf.Ln(4)
		pdf.Text(17, pdf.GetY(), trObj.translDef("квалификаионе картице и сертификата"))
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
	createSimpleHeader(pdf, cirTr)

	pdf.Ln(5)
	pdf.SetFont("Arimo-Bold", "", 9)
	pdf.Text(15, pdf.GetY(), trObj.translDef(fmt.Sprintf("Дневник предавача семинара унапређења знања на %s обуци професионалних возача", day.Seminar.SeminarTheme.BaseSeminarType.GetSeminarTypeForSentence())))
	pdf.Ln(5)
	pdf.SetFont("Arimo-Regular", "", 9)
	pdf.Text(15, pdf.GetY(), trObj.translDef("Датум одржавања семинара"))
	pdf.Text(80, pdf.GetY(), day.Date.Format("02.01.2006."))
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

	getClassTime := func(d time.Time, i int) string {
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

	for i := 0; i < len(day.Classes); i++ {
		pdf.Ln(ch)
		pdf.CellFormat(20, ch, strconv.Itoa(i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, ch, getClassTime(day.Date, i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(90, ch, trObj.translDef(day.Classes[i].Name), "1", 0, "C", false, 0, "")
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
	pdf.Text(120, pdf.GetY(), trObj.translDef(day.Seminar.GetCode()))

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (p *printService) PrintTestBarcode() ([]byte, error) {
	url := fmt.Sprintf("%s/do-test")
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
	pdf.Text(38, pdf.GetY(), trObj.translate(fmt.Sprintf("План реализације наставе за %s обуку - 7 часова", day.Seminar.SeminarTheme.BaseSeminarType.GetSeminarTypeForSentence()), 14))
	pdf.Ln(5)

	pdf.SetFont("Arimo-Bold", "", 11)
	ch := 8.0
	chs := 5.0
	cw1 := 12.0
	cw2 := 69.0
	cw3 := 53.0
	cw4 := 23.0
	cw5 := 23.0
	pdf.CellFormat(180, ch, trObj.translDef(util.GetDaySerbian(day.Date))+" "+day.Date.Format("01.02.2006"), "1", 0, "C", false, 0, "")
	pdf.Ln(ch)
	pdf.SetFillColor(180, 197, 231)
	pdf.CellFormat(cw1, ch, trObj.translDef("Р.бр"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw2, ch, trObj.translDef("Назив часа"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw3, ch, trObj.translDef("Предавач"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw4, ch, trObj.translDef("Почетак"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw5, ch, trObj.translDef("Крај"), "1", 0, "C", true, 0, "")
	pdf.SetFont("Arimo-Regular", "", 10)
	pdf.Ln(ch)
	pdf.CellFormat(cw1, chs*2.0, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw2, chs*2.0, trObj.translate("Prijava i evidentiranje polaznika", 10), "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw3, chs*2.0, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw4, chs*2.0, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw5, chs*2.0, "", "1", 0, "C", false, 0, "")

	pdf.Ln(chs * 2.0)
	pdf.CellFormat(cw1, chs*3.0, "1", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw2, chs*3.0, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw3, chs*3.0, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw4, chs*3.0, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw5, chs*3.0, "", "1", 0, "C", false, 0, "")

	pdf.SetFont("Arimo-Bold", "", 11)
	pdf.Ln(chs * 3.0)
	pdf.SetFillColor(180, 197, 231)
	pdf.CellFormat(cw1, chs*2, latTr(""), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw2, chs*2, trObj.translDef("Pauza za kafu"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw3, chs*2, trObj.translDef("Predavač"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw4, chs*2, latTr("10"), "LTB", 0, "R", true, 0, "")
	pdf.CellFormat(cw5, chs*2, trObj.translDef("minuta"), "TBR", 0, "L", true, 0, "")

	pdf.SetFont("Arimo-Regular", "", 10)
	pdf.Ln(chs * 2.0)
	pdf.CellFormat(cw1, chs*2.0, "2", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw2, chs*2.0, latTr("Prijava i evidentiranje polaznika"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw3, chs*2.0, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw4, chs*2.0, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw5, chs*2.0, "", "1", 0, "C", false, 0, "")

	pdf.SetFont("Arimo-Bold", "", 11)
	pdf.Ln(chs * 2.0)
	pdf.SetFillColor(180, 197, 231)
	pdf.CellFormat(cw1, chs*2, latTr(""), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw2, chs*2, latTr("Pauza"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw3, chs*2, latTr(""), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw4, chs*2, latTr("5"), "LTB", 0, "R", true, 0, "")
	pdf.CellFormat(cw5, chs*2, latTr("minuta"), "TBR", 0, "L", true, 0, "")

	pdf.SetFont("Arimo-Regular", "", 10)
	pdf.Ln(chs * 2.0)
	pdf.CellFormat(cw1, chs*2.0, "3", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw2, chs*2.0, latTr("Prijava i evidentiranje polaznika"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw3, chs*2.0, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw4, chs*2.0, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw5, chs*2.0, "", "1", 0, "C", false, 0, "")

	pdf.SetFont("Arimo-Bold", "", 11)
	pdf.Ln(chs * 2.0)
	pdf.SetFillColor(180, 197, 231)
	pdf.CellFormat(cw1, chs*2, latTr(""), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw2, chs*2, latTr("Pauza za doručak"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw3, chs*2, latTr(""), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw4, chs*2, latTr("25"), "LTB", 0, "R", true, 0, "")
	pdf.CellFormat(cw5, chs*2, latTr("minuta"), "TBR", 0, "L", true, 0, "")

	pdf.SetFont("Arimo-Regular", "", 10)
	pdf.Ln(chs * 2.0)
	pdf.CellFormat(cw1, chs*2.0, "4", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw2, chs*2.0, latTr("Prijava i evidentiranje polaznika"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw3, chs*2.0, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw4, chs*2.0, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw5, chs*2.0, "", "1", 0, "C", false, 0, "")

	pdf.SetFont("Arimo-Bold", "", 11)
	pdf.Ln(chs * 2.0)
	pdf.SetFillColor(180, 197, 231)
	pdf.CellFormat(cw1, chs*2, latTr(""), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw2, chs*2, latTr("Pauza"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw3, chs*2, latTr(""), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw4, chs*2, latTr("5"), "LTB", 0, "R", true, 0, "")
	pdf.CellFormat(cw5, chs*2, latTr("minuta"), "TBR", 0, "L", true, 0, "")

	pdf.SetFont("Arimo-Regular", "", 10)
	pdf.Ln(chs * 2.0)
	pdf.CellFormat(cw1, chs*2.0, "5", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw2, chs*2.0, latTr("Prijava i evidentiranje polaznika"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw3, chs*2.0, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw4, chs*2.0, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw5, chs*2.0, "", "1", 0, "C", false, 0, "")

	pdf.SetFont("Arimo-Bold", "", 11)
	pdf.Ln(chs * 2.0)
	pdf.SetFillColor(180, 197, 231)
	pdf.CellFormat(cw1, chs*2, latTr(""), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw2, chs*2, latTr("Pauza za kafu"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw3, chs*2, latTr(""), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw4, chs*2, latTr("15"), "LTB", 0, "R", true, 0, "")
	pdf.CellFormat(cw5, chs*2, latTr("minuta"), "TBR", 0, "L", true, 0, "")

	pdf.SetFont("Arimo-Regular", "", 10)
	pdf.Ln(chs * 2.0)
	pdf.CellFormat(cw1, chs*2.0, "6", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw2, chs*2.0, latTr("Prijava i evidentiranje polaznika"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw3, chs*2.0, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw4, chs*2.0, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw5, chs*2.0, "", "1", 0, "C", false, 0, "")

	pdf.SetFont("Arimo-Bold", "", 11)
	pdf.Ln(chs * 2.0)
	pdf.SetFillColor(180, 197, 231)
	pdf.CellFormat(cw1, chs*2, latTr(""), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw2, chs*2, latTr("Pauza"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw3, chs*2, latTr(""), "1", 0, "C", true, 0, "")
	pdf.CellFormat(cw4, chs*2, latTr("5"), "LTB", 0, "R", true, 0, "")
	pdf.CellFormat(cw5, chs*2, latTr("minuta"), "TBR", 0, "L", true, 0, "")

	pdf.SetFont("Arimo-Regular", "", 10)
	pdf.Ln(chs * 2.0)
	pdf.CellFormat(cw1, chs*2.0, "7", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw2, chs*2.0, latTr("Prijava i evidentiranje polaznika"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw3, chs*2.0, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw4, chs*2.0, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw5, chs*2.0, "", "1", 0, "C", false, 0, "")

	pdf.Ln(chs * 2.0)
	pdf.CellFormat(cw1, chs*2.0, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw2, chs*2.0, latTr("Prijava i evidentiranje polaznika"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw3, chs*2.0, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw4, chs*2.0, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(cw5, chs*2.0, "", "1", 0, "C", false, 0, "")

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}
