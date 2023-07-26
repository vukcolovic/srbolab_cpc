package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"srbolab_cpc/logoped"
	"srbolab_cpc/service"
	"strconv"
)

func PrintSeminarStudentList(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	seminarIdParam, ok := vars["seminar_id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter seminar_id")
		SetErrorResponse(w, NewMissingRequestParamError("seminar_id"))
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
		SetErrorResponse(w, errors.New("Greška štampanja, greška prilikom povlačenja seminara: "+err.Error()))
		return
	}

	report, err := service.PrintService.PrintSeminarStudentList(seminar)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška štampanja: "+err.Error()))
		return
	}

	SetSuccessResponse(w, report)
}

func PrintConfirmationStatements(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	seminarIdParam, ok := vars["seminar_id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter seminar_id")
		SetErrorResponse(w, NewMissingRequestParamError("seminar_id"))
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
		SetErrorResponse(w, errors.New("Greška štampanja, greška prilikom povlačenja seminara: "+err.Error()))
		return
	}

	report, err := service.PrintService.PrintConfirmationStatements(seminar)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška štampanja: "+err.Error()))
		return
	}

	SetSuccessResponse(w, report)
}

func PrintConfirmations(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	seminarIdParam, ok := vars["seminar_id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter seminar_id")
		SetErrorResponse(w, NewMissingRequestParamError("seminar_id"))
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
		SetErrorResponse(w, errors.New("Greška štampanja, greška prilikom povlačenja seminara: "+err.Error()))
		return
	}

	report, err := service.PrintService.PrintConfirmations(seminar)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška štampanja: "+err.Error()))
		return
	}

	SetSuccessResponse(w, report)
}
