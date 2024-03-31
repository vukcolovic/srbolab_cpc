package handlers

import (
	"errors"
	"net/http"
	"srbolab_cpc/logoped"
	"srbolab_cpc/service"
	"strconv"

	"github.com/gorilla/mux"
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
	report, err := service.PrintService.PrintTestBarcode()
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

func PrintPayments(w http.ResponseWriter, req *http.Request) {
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
	seminar, err := service.SeminarService.GetSeminarByIDWithClientFiles(seminarId)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška štampanja, greška prilikom povlačenja seminara: "+err.Error()))
		return
	}

	report, err := service.PrintService.PrintPayments(seminar)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška štampanja: "+err.Error()))
		return
	}

	SetSuccessResponse(w, report)
}

func PrintReport(w http.ResponseWriter, req *http.Request) {
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

	report, err := service.PrintService.PrintSeminarReport(seminar)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška štampanja: "+err.Error()))
		return
	}

	SetSuccessResponse(w, report)
}

func PrintReport2(w http.ResponseWriter, req *http.Request) {
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

	report, err := service.PrintService.PrintSeminarReport2(seminar)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška štampanja: "+err.Error()))
		return
	}

	SetSuccessResponse(w, report)
}

func PrintTest(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	testIdParam, ok := vars["test_id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter test_id")
		SetErrorResponse(w, NewMissingRequestParamError("test_id"))
		return
	}

	test_id, err := strconv.Atoi(testIdParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("testId", testIdParam))
		return
	}
	test, err := service.TestService.GetTestWithAnswersByID(test_id)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška štampanja, greška prilikom povlačenja testa: "+err.Error()))
		return
	}

	report, err := service.PrintService.PrintTest(test)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška štampanja: "+err.Error()))
		return
	}

	SetSuccessResponse(w, report)
}
