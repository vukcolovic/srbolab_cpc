package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"srbolab_cpc/logoped"
	"srbolab_cpc/model"
	"srbolab_cpc/service"
	"strconv"
	"time"
)

func ListTests(w http.ResponseWriter, r *http.Request) {
	tests, err := service.TestService.GetAllTests()
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste testova: "+err.Error()))
		return
	}

	SetSuccessResponse(w, tests)
}

func ListTestsBySeminarThemeID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	seminarThemeIdParam, ok := params["id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter seminarThemeId")
		SetErrorResponse(w, NewMissingRequestParamError("seminarThemeId"))
		return
	}
	seminarThemeId, err := strconv.Atoi(seminarThemeIdParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("seminarThemeId", seminarThemeIdParam))
		return
	}

	tests, err := service.TestService.GetAllTestsBySeminarTheme(seminarThemeId)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste testova: "+err.Error()))
		return
	}

	SetSuccessResponse(w, tests)
}

func GetTestByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	testIDParam, ok := vars["id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter id")
		SetErrorResponse(w, NewMissingRequestParamError("id"))
		return
	}

	testID, err := strconv.Atoi(testIDParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("TestID", testIDParam))
		return
	}
	test, err := service.TestService.GetTestByID(testID)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja testa: "+err.Error()))
		return
	}

	SetSuccessResponse(w, test)
}

func CreateTest(w http.ResponseWriter, r *http.Request) {
	var test model.Test
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&test)
	if err != nil {
		logoped.ErrorLog.Println("unable to retrieve the just parsed code")
		SetErrorResponse(w, NewJSONDecodeError("Test"))
		return
	}

	createdTest, err := service.TestService.CreateTest(test)
	if err != nil {
		logoped.ErrorLog.Println("Error creating test " + err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom kreiranja testa: "+err.Error()))
		return
	}

	SetSuccessResponse(w, createdTest)
}

func UpdateTest(w http.ResponseWriter, r *http.Request) {
	var test model.Test
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&test)
	if err != nil {
		logoped.ErrorLog.Println("Error decoding test: ", err)
		SetErrorResponse(w, NewJSONDecodeError("Test"))
		return
	}

	updatedTest, err := service.TestService.UpdateTest(test)
	if err != nil {
		logoped.ErrorLog.Println("Error updating test: ", err)
		SetErrorResponse(w, errors.New("Greška prilikom ažuriranja testa: "+err.Error()))
		return
	}

	SetSuccessResponse(w, updatedTest)
}

func GetClientTestsBySeminarDayAndJMBG(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	seminarDayIDParam, ok := vars["seminar-day"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter seminar-day")
		SetErrorResponse(w, NewMissingRequestParamError("SeminarDayID"))
		return
	}

	seminarDayID, err := strconv.Atoi(seminarDayIDParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("SeminarDayID", seminarDayIDParam))
		return
	}

	jmbg, ok := vars["jmbg"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter jmbg")
		SetErrorResponse(w, NewMissingRequestParamError("JMBG"))
		return
	}

	clientTests, err := service.TestService.GetClientTestBySeminarDayIDAndJMBG(seminarDayID, jmbg)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja testova klijenta: "+err.Error()))
		return
	}

	SetSuccessResponse(w, clientTests)
}

func GetClientTestsBySeminarDay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seminarDayIDParam, ok := vars["seminar-day"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter seminar-day")
		SetErrorResponse(w, NewMissingRequestParamError("SeminarDayID"))
		return
	}

	seminarDayID, err := strconv.Atoi(seminarDayIDParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("SeminarDayID", seminarDayIDParam))
		return
	}

	clientTests, err := service.TestService.GetClientTestBySeminarDayID(seminarDayID)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja testova klijenta: "+err.Error()))
		return
	}

	SetSuccessResponse(w, clientTests)
}

func SaveClientTest(w http.ResponseWriter, r *http.Request) {
	var clientTest model.ClientTest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&clientTest)
	if err != nil {
		logoped.ErrorLog.Println("unable to retrieve the just parsed code")
		SetErrorResponse(w, NewJSONDecodeError("ClientTest"))
		return
	}
	logoped.InfoLog.Println(fmt.Sprintf("Saving client test, client %s, seminar day %d, test %d",
		clientTest.Jmbg, clientTest.SeminarDay.ID, clientTest.Test.ID))

	msg, err := isTestValid(&clientTest)
	if err != nil {
		logoped.ErrorLog.Println("Error saving client test, test is not valid, error: ", err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom snimanja testa: "+err.Error()))
		return
	}
	if len(msg) > 0 {
		logoped.ErrorLog.Println("Error saving client test, test is not valid, message: ", msg)
		SetErrorResponse(w, errors.New(msg))
		return
	}

	createdClientTest, err := service.TestService.CreateClientTest(clientTest)
	if err != nil {
		logoped.ErrorLog.Println("Error creating client test " + err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom kreiranja testa: "+err.Error()))
		return
	}

	SetSuccessResponse(w, createdClientTest)
}

func isTestValid(clientTest *model.ClientTest) (string, error) {
	client, err := service.ClientService.GetClientByJMBG(clientTest.Jmbg)
	if err != nil {
		return "", err
	}
	if client == nil {
		return "Korisnik sa ovim jmbg-om ne postoji u sistemu.", nil
	}

	day, err := service.SeminarDayService.GetSeminarDayByID(int(clientTest.SeminarDay.ID))
	if err != nil {
		return "", err
	}

	found := false
	for _, tr := range day.Seminar.Trainees {
		if tr.ClientID == client.ID {
			found = true
			break
		}
	}

	if !found {
		return "Korisnik sa ovim jmbg-om nije učesnik seminara.", nil
	}

	dayTime := clientTest.SeminarDay.Date
	now := time.Now()
	if now.Day() != dayTime.Day() || now.Month() != dayTime.Month() || now.Year() != dayTime.Year() {
		return "Ovaj test nije dozvoljeno raditi danas.", nil
	}

	if clientTest.SeminarDay.Seminar.SeminarStatusID != model.SEMINAR_STATUS_IN_PROGRESS {
		return "Ovaj test nije dozvoljeno raditi, seminar nije u toku.", nil
	}

	if clientTest.SeminarDay.TestID == nil {
		return "Ovaj test nije dozvoljeno raditi, test nije odabran.", nil
	}

	tests, err := service.TestService.GetClientTestBySeminarDayIDAndJMBG(int(clientTest.SeminarDay.ID), *client.JMBG)
	if err != nil {
		return "", err
	}

	if len(tests) > 1 {
		return "Ovaj test nije dozvoljeno, klijent je već odradio dva testa u toku dana.", nil
	}

	if len(tests) == 1 {
		if !tests[0].CreatedAt.Add(2 * time.Hour).Before(time.Now()) {
			return "Nije dozvoljeno snimiti test, rađen je skoro.", nil
		}
	}

	clientTest.Client = *client
	return "", nil
}
