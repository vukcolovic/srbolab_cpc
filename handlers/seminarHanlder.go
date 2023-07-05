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

func CreateSeminar(w http.ResponseWriter, r *http.Request) {
	var seminar model.Seminar
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&seminar)
	if err != nil {
		logoped.ErrorLog.Println("unable to retrieve the just parsed code")
		SetErrorResponse(w, NewJSONDecodeError("Seminar"))
		return
	}

	createdSeminar, err := service.SeminarService.CreateSeminar(seminar)
	if err != nil {
		logoped.ErrorLog.Println("Error creating seminar " + err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom kreiranja seminara: "+err.Error()))
		return
	}

	SetSuccessResponse(w, createdSeminar)
}

func ListSeminars(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	skipParam, ok := queryParams["skip"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter skip")
		SetErrorResponse(w, NewMissingRequestParamError("skip"))
		return
	}
	skip, err := strconv.Atoi(skipParam[0])
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("skip", skipParam[0]))
		return
	}

	takeParam, ok := queryParams["take"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter take")
		SetErrorResponse(w, NewMissingRequestParamError("take"))
		return
	}
	take, err := strconv.Atoi(takeParam[0])
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("take", takeParam[0]))
		return
	}

	seminars, err := service.SeminarService.GetAllSeminars(skip, take)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste seminara: "+err.Error()))
		return
	}

	SetSuccessResponse(w, seminars)
}

func GetSeminarByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	seminarIdParam, ok := vars["id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter id")
		SetErrorResponse(w, NewMissingRequestParamError("id"))
		return
	}

	seminarId, err := strconv.Atoi(seminarIdParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("seminarId", seminarIdParam))
		return
	}
	seminar, err := service.SeminarService.GetSeminarByID(seminarId)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja seminara: "+err.Error()))
		return
	}

	SetSuccessResponse(w, seminar)
}

func UpdateSeminar(w http.ResponseWriter, r *http.Request) {
	var seminar model.Seminar
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&seminar)
	if err != nil {
		logoped.ErrorLog.Println("Error decoding seminar: ", err)
		SetErrorResponse(w, NewJSONDecodeError("seminar"))
		return
	}

	updatedSeminar, err := service.SeminarService.UpdateSeminar(seminar)
	if err != nil {
		logoped.ErrorLog.Println("Error updating client: ", err)
		SetErrorResponse(w, errors.New("Greška prilikom ažuriranja seminara: "+err.Error()))
		return
	}

	SetSuccessResponse(w, updatedSeminar)
}

func DeleteSeminar(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	seminarId, ok := vars["id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter id")
		SetErrorResponse(w, NewMissingRequestParamError("id"))
		return
	}

	seminarIdInt, err := strconv.Atoi(seminarId)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("seminarId", seminarId))
		return
	}

	err = service.SeminarService.DeleteSeminar(seminarIdInt)
	if err != nil {
		logoped.ErrorLog.Println("error deleting seminar ", err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom brisanja seminara: "+err.Error()))
		return
	}

	SetSuccessResponse(w, nil)
}

func CountSeminars(w http.ResponseWriter, req *http.Request) {
	count, err := service.SeminarService.GetSeminarsCount()
	if err != nil {
		logoped.ErrorLog.Println("error getting seminar count")
		SetErrorResponse(w, errors.New("Greška prilikom dobijanja broja seminara: "+err.Error()))
		return
	}

	SetSuccessResponse(w, count)
}
