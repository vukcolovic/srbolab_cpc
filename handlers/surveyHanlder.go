package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"srbolab_cpc/logoped"
	"srbolab_cpc/model"
	"srbolab_cpc/service"
	"strconv"

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

func GetActiveSurvey(w http.ResponseWriter, req *http.Request) {
	survey, err := service.SurveyService.GetActiveSurvey()
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja ankete: "+err.Error()))
		return
	}

	SetSuccessResponse(w, survey)
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
