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

func PrintClientTestsBySeminarDay(w http.ResponseWriter, r *http.Request) {
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
		SetErrorResponse(w, errors.New("Greška prilikom štampanja Excel izveštaja za testove vozača: "+err.Error()))
		return
	}

	excel, err := service.ExcelService.CreateClientTestsBySeminarDayReport(clientTests)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom štampanja Excel izveštaja za testove vozača: "+err.Error()))
		return
	}

	SetSuccessResponse(w, excel)
}

func PrintListTraineesBySeminarDay(w http.ResponseWriter, r *http.Request) {
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

	seminarDay, err := service.SeminarDayService.GetSeminarDayByID(seminarDayID)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom štampanja Excel izveštaja za spisak polaznika: "+err.Error()))
		return
	}

	excel, err := service.ExcelService.CreateListClientsBySeminarDayReport(seminarDay)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom štampanja Excel izveštaja za spisak polaznika: "+err.Error()))
		return
	}

	SetSuccessResponse(w, excel)
}

func PrintSeminarsReportOfClients(w http.ResponseWriter, r *http.Request) {
	var filter model.SeminarFilter
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&filter)
	if err != nil {
		logoped.ErrorLog.Println("Error decoding seminar filter: ", err)
		SetErrorResponse(w, NewJSONDecodeError("SeminarFilter"))
		return
	}

	seminars, err := service.SeminarService.GetAllSeminarsWithTrainees(0, 10000, filter)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom štampanja Excel izveštaja za statistiku polaznika: "+err.Error()))
		return
	}

	excel, err := service.ExcelService.CreateSeminarsReportOfClients(seminars)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom štampanja Excel izveštaja za statistiku polaznika: "+err.Error()))
		return
	}

	SetSuccessResponse(w, excel)
}

func PrintSeminarsReportOfTeachers(w http.ResponseWriter, r *http.Request) {
	var filter model.SeminarFilter
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&filter)
	if err != nil {
		logoped.ErrorLog.Println("Error decoding seminar filter: ", err)
		SetErrorResponse(w, NewJSONDecodeError("SeminarFilter"))
		return
	}

	seminars, err := service.SeminarService.GetAllSeminarsWithSeminarDays(0, 10000, filter)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom štampanja Excel izveštaja za statistiku predavača: "+err.Error()))
		return
	}

	excel, err := service.ExcelService.CreateSeminarsReportOfTeachers(seminars)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom štampanja Excel izveštaja za statistiku predavača: "+err.Error()))
		return
	}

	SetSuccessResponse(w, excel)
}
