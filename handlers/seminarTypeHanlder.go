package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"srbolab_cpc/logoped"
	"srbolab_cpc/service"
	"strconv"
)

func ListBaseSeminarTypes(w http.ResponseWriter, r *http.Request) {
	seminarTypes, err := service.SeminarTypeService.GetAllBaseSeminarTypes()
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste tipova seminara: "+err.Error()))
		return
	}

	SetSuccessResponse(w, seminarTypes)
}

func ListSeminarThemesBySeminarType(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	seminarTypeIdParam, ok := params["seminarTypeId"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter skip")
		SetErrorResponse(w, NewMissingRequestParamError("seminarTypeId"))
		return
	}
	seminarTypeId, err := strconv.Atoi(seminarTypeIdParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("seminarTypeId", seminarTypeIdParam))
		return
	}

	seminarThemes, err := service.SeminarTypeService.GetSeminarThemesBypeID(seminarTypeId)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste tema seminara: "+err.Error()))
		return
	}

	SetSuccessResponse(w, seminarThemes)
}
