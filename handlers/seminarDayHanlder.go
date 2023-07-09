package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"srbolab_cpc/logoped"
	"srbolab_cpc/model"
	"srbolab_cpc/service"
	"strconv"
)

func CreateSeminarDay(w http.ResponseWriter, r *http.Request) {
	var day model.SeminarDay
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&day)
	if err != nil {
		logoped.ErrorLog.Println("unable to retrieve the just parsed code")
		SetErrorResponse(w, NewJSONDecodeError("SeminarDay"))
		return
	}

	createdDay, err := service.SeminarDayService.CreateSeminarDay(day)
	if err != nil {
		logoped.ErrorLog.Println("Error creating seminar day" + err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom kreiranja seminar dana: "+err.Error()))
		return
	}

	SetSuccessResponse(w, createdDay)
}

func UpdateSeminarDay(w http.ResponseWriter, r *http.Request) {
	var day model.SeminarDay
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&day)
	if err != nil {
		logoped.ErrorLog.Println("unable to retrieve the just parsed code")
		SetErrorResponse(w, NewJSONDecodeError("SeminarDay"))
		return
	}

	updatedDay, err := service.SeminarDayService.UpdateSeminarDay(day)
	if err != nil {
		logoped.ErrorLog.Println("Error updating seminar day" + err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom ažuriranja seminar dana: "+err.Error()))
		return
	}

	SetSuccessResponse(w, updatedDay)
}

func ListSeminarDays(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	seminarIDParam, ok := queryParams["seminar_id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter skip")
		SetErrorResponse(w, NewMissingRequestParamError("seminar_id"))
		return
	}
	seminarID, err := strconv.Atoi(seminarIDParam[0])
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("seminar_id", seminarIDParam[0]))
		return
	}

	seminarDays, err := service.SeminarDayService.GetSeminarDaysBySeminarID(seminarID)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste seminar dana: "+err.Error()))
		return
	}

	SetSuccessResponse(w, seminarDays)
}
