package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"srbolab_cpc/logoped"
	"srbolab_cpc/model"
	"srbolab_cpc/service"
	"strconv"
)

func ListQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := service.QuestionService.GetAllQuestions(0, 10000)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste pitanja: "+err.Error()))
		return
	}

	SetSuccessResponse(w, questions)
}

func GetQuestionByID(w http.ResponseWriter, req *http.Request) {
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
		SetErrorResponse(w, NewWrongParamFormatErrorError("locationID", questionIDParam))
		return
	}
	question, err := service.QuestionService.GetQuestionByID(questionID)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja pitanja: "+err.Error()))
		return
	}

	SetSuccessResponse(w, question)
}

func CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var question model.Question
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&question)
	if err != nil {
		logoped.ErrorLog.Println("unable to retrieve the just parsed code")
		SetErrorResponse(w, NewJSONDecodeError("Question"))
		return
	}

	createdQuestion, err := service.QuestionService.CreateQuestion(question)
	if err != nil {
		logoped.ErrorLog.Println("Error creating question " + err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom kreiranja pitanja: "+err.Error()))
		return
	}

	SetSuccessResponse(w, createdQuestion)
}

func UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	var question model.Question
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&question)
	if err != nil {
		logoped.ErrorLog.Println("Error decoding question: ", err)
		SetErrorResponse(w, NewJSONDecodeError("Question"))
		return
	}

	updatedQuestion, err := service.QuestionService.UpdateQuestion(question)
	if err != nil {
		logoped.ErrorLog.Println("Error updating question: ", err)
		SetErrorResponse(w, errors.New("Greška prilikom ažuriranja pitanja: "+err.Error()))
		return
	}

	SetSuccessResponse(w, updatedQuestion)
}
