package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"srbolab_cpc/logoped"
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
