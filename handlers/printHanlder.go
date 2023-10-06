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

func PrintConfirmationReceives(w http.ResponseWriter, req *http.Request) {
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

	report, err := service.PrintService.PrintConfirmationReceives(seminar)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška štampanja: "+err.Error()))
		return
	}

	SetSuccessResponse(w, report)
}

func PrintMuster(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	seminarDayIdParam, ok := vars["seminar_day_id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter seminar_day_id")
		SetErrorResponse(w, NewMissingRequestParamError("seminar_day_id"))
		return
	}

	seminarDayId, err := strconv.Atoi(seminarDayIdParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("seminarDayId", seminarDayIdParam))
		return
	}
	seminarDay, err := service.SeminarDayService.GetSeminarDayByID(seminarDayId)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška štampanja, greška prilikom povlačenja seminar dana: "+err.Error()))
		return
	}

	report, err := service.PrintService.PrintMuster(seminarDay)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška štampanja: "+err.Error()))
		return
	}

	SetSuccessResponse(w, report)
}

func PrintSeminarEvidence(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	seminarDayIdParam, ok := vars["seminar_day_id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter seminar_day_id")
		SetErrorResponse(w, NewMissingRequestParamError("seminar_day_id"))
		return
	}

	seminarDayId, err := strconv.Atoi(seminarDayIdParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("seminarDayId", seminarDayIdParam))
		return
	}
	seminarDay, err := service.SeminarDayService.GetSeminarDayByID(seminarDayId)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška štampanja, greška prilikom povlačenja seminar dana: "+err.Error()))
		return
	}

	report, err := service.PrintService.PrintSeminarEvidence(seminarDay)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška štampanja: "+err.Error()))
		return
	}

	SetSuccessResponse(w, report)
}

func PrintTestBarcode(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	seminarDayIdParam, ok := vars["seminar_day_id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter seminar_day_id")
		SetErrorResponse(w, NewMissingRequestParamError("seminar_day_id"))
		return
	}

	seminarDayId, err := strconv.Atoi(seminarDayIdParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("seminarDayId", seminarDayIdParam))
		return
	}
	seminarDay, err := service.SeminarDayService.GetSeminarDayByID(seminarDayId)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška štampanja, greška prilikom povlačenja seminar dana: "+err.Error()))
		return
	}

	report, err := service.PrintService.PrintTestBarcode(seminarDay)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška štampanja: "+err.Error()))
		return
	}

	SetSuccessResponse(w, report)
}

func PrintCheckIn(w http.ResponseWriter, req *http.Request) {
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

	report, err := service.PrintService.PrintCheckIn(seminar)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška štampanja: "+err.Error()))
		return
	}

	SetSuccessResponse(w, report)
}

func PrintPlanTreningRealization(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	seminarDayIdParam, ok := vars["seminar_day_id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter seminar_day_id")
		SetErrorResponse(w, NewMissingRequestParamError("seminar_day_id"))
		return
	}

	seminarDayId, err := strconv.Atoi(seminarDayIdParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("seminarDayId", seminarDayIdParam))
		return
	}
	seminarDay, err := service.SeminarDayService.GetSeminarDayByID(seminarDayId)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška štampanja, greška prilikom povlačenja seminar dana: "+err.Error()))
		return
	}

	report, err := service.PrintService.PrintPlanTreningRealization(seminarDay)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška štampanja: "+err.Error()))
		return
	}

	SetSuccessResponse(w, report)
}
