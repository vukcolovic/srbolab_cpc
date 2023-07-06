package handlers

import (
	"errors"
	"net/http"
	"srbolab_cpc/logoped"
	"srbolab_cpc/service"
)

func ListSeminarTypes(w http.ResponseWriter, r *http.Request) {
	seminarTypes, err := service.SeminarTypeService.GetAllSeminarTypes()
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste tipova seminara: "+err.Error()))
		return
	}

	SetSuccessResponse(w, seminarTypes)
}
