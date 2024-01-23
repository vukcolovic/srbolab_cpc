package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"srbolab_cpc/logoped"
	"srbolab_cpc/model"
	"srbolab_cpc/service"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func ListSurveyQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := service.SurveyService.GetAllSurveyQuestions()
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste pitanja: "+err.Error()))
		return
	}

	SetSuccessResponse(w, questions)
}

func GetSurveyQuestionByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	questionIDParam, ok := vars["id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter id")
		SetErrorResponse(w, NewMissingRequestParamError("id"))
		return
	}

	questionID, err := strconv.Atoi(questionIDParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("QuestionID", questionIDParam))
		return
	}
	question, err := service.SurveyService.GetSurveyQuestionByID(questionID)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja pitanja: "+err.Error()))
		return
	}

	SetSuccessResponse(w, question)
}

func CreateSurveyQuestion(w http.ResponseWriter, r *http.Request) {
	var question model.SurveyQuestion
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&question)
	if err != nil {
		logoped.ErrorLog.Println("unable to retrieve the just parsed code")
		SetErrorResponse(w, NewJSONDecodeError("SurveyQuestion"))
		return
	}

	createdQuestion, err := service.SurveyService.CreateSurveyQuestion(question)
	if err != nil {
		logoped.ErrorLog.Println("Error creating question " + err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom kreiranja pitanja: "+err.Error()))
		return
	}

	SetSuccessResponse(w, createdQuestion)
}

func ListSurveys(w http.ResponseWriter, r *http.Request) {
	surveys, err := service.SurveyService.GetAllSurveys()
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste anketa: "+err.Error()))
		return
	}

	SetSuccessResponse(w, surveys)
}

func GetSurveyByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	surveyIDParam, ok := vars["id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter id")
		SetErrorResponse(w, NewMissingRequestParamError("id"))
		return
	}

	surveyID, err := strconv.Atoi(surveyIDParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("SurveyID", surveyIDParam))
		return
	}
	survey, err := service.SurveyService.GetSurveyByID(surveyID)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja ankete: "+err.Error()))
		return
	}

	SetSuccessResponse(w, survey)
}

func CreateSurvey(w http.ResponseWriter, r *http.Request) {
	var survey model.Survey
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&survey)
	if err != nil {
		logoped.ErrorLog.Println("unable to retrieve the just parsed code")
		SetErrorResponse(w, NewJSONDecodeError("Survey"))
		return
	}

	createdSurvey, err := service.SurveyService.CreateSurvey(survey)
	if err != nil {
		logoped.ErrorLog.Println("Error creating survey " + err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom kreiranja ankete: "+err.Error()))
		return
	}

	SetSuccessResponse(w, createdSurvey)
}

func GetActiveSurveys(w http.ResponseWriter, req *http.Request) {
	surveys, err := service.SurveyService.GetActiveSurveys()
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja aktivnih anketa: "+err.Error()))
		return
	}

	SetSuccessResponse(w, surveys)
}

func SaveClientSurvey(w http.ResponseWriter, r *http.Request) {
	var clientSurvey model.ClientSurvey
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&clientSurvey)
	if err != nil {
		logoped.ErrorLog.Println("unable to retrieve the just parsed code")
		SetErrorResponse(w, NewJSONDecodeError("ClientSurvey"))
		return
	}

	msg, err := isSurveyValid(&clientSurvey)
	if err != nil {
		logoped.ErrorLog.Println("Error saving client survey, survey is not valid, error: ", err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom snimanja ankete: "+err.Error()))
		return
	}
	if len(msg) > 0 {
		logoped.ErrorLog.Println("Error saving client survey, survey is not valid, message: ", msg)
		SetErrorResponse(w, errors.New(msg))
		return
	}

	client, err := service.ClientService.GetClientByJMBG(clientSurvey.JMBG)
	if err != nil {
		logoped.ErrorLog.Println("Error creating client survey " + err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom kreiranja ankete: "+err.Error()))
		return
	}

	clientSurvey.Client = *client
	clientSurvey.ClientID = client.ID

	createdClientSurvey, err := service.SurveyService.CreateClientSurvey(clientSurvey)
	if err != nil {
		logoped.ErrorLog.Println("Error creating client survey " + err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom kreiranja ankete: "+err.Error()))
		return
	}

	SetSuccessResponse(w, createdClientSurvey)
}

func isSurveyValid(clientSurvey *model.ClientSurvey) (string, error) {
	client, err := service.ClientService.GetClientByJMBG(clientSurvey.JMBG)
	if err != nil {
		return "", err
	}
	if client == nil {
		return "Korisnik sa ovim jmbg-om ne postoji u sistemu.", nil
	}

	day, err := service.SeminarDayService.GetSeminarDayByID(int(clientSurvey.SeminarDay.ID))
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

	dayTime := clientSurvey.SeminarDay.Date
	now := time.Now()
	if now.Day() != dayTime.Day() || now.Month() != dayTime.Month() || now.Year() != dayTime.Year() {
		return "Ovu ankety nije dozvoljeno raditi danas.", nil
	}

	if clientSurvey.SeminarDay.Seminar.SeminarStatusID != model.SEMINAR_STATUS_IN_PROGRESS {
		return "Ovu anketu nije dozvoljeno raditi, seminar nije u toku.", nil
	}

	clientSurveys, err := service.SurveyService.GetClientSurveysBySeminarDayIDAndClientID(int(clientSurvey.SeminarDay.ID), int(client.ID))
	if err != nil {
		return "", err
	}

	if len(clientSurveys) > 8 {
		return "Ovu anketu nije dozvoljeno raditi, klijent je već odradio anketu u toku dana.", nil
	}

	return "", nil
}
