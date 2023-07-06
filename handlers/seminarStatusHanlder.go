package handlers

import (
	"errors"
	"net/http"
	"srbolab_cpc/logoped"
	"srbolab_cpc/service"
)

func ListSeminarStatuses(w http.ResponseWriter, r *http.Request) {
	seminarStatuses, err := service.SeminarStatusService.GetAllSeminarStatuses()
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste statusa seminara: "+err.Error()))
		return
	}

	SetSuccessResponse(w, seminarStatuses)
}
