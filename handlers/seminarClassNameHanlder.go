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

func GetSeminarClassNameByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	classNameIdParam, ok := vars["id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter id")
		SetErrorResponse(w, NewMissingRequestParamError("id"))
		return
	}

	classNameId, err := strconv.Atoi(classNameIdParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("classNameId", classNameIdParam))
		return
	}
	className, err := service.SeminarClassService.GetSeminarClassNameByID(classNameId)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja naziva časa: "+err.Error()))
		return
	}

	SetSuccessResponse(w, className)
}

func CreateSeminarClassName(w http.ResponseWriter, r *http.Request) {
	var className model.SeminarClassName
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&className)
	if err != nil {
		logoped.ErrorLog.Println("unable to retrieve the just parsed code")
		SetErrorResponse(w, NewJSONDecodeError("SeminarClassName"))
		return
	}

	createdClassName, err := service.SeminarClassService.CreateSeminarClassName(className)
	if err != nil {
		logoped.ErrorLog.Println("Error creating class name" + err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom kreiranja naziva časa: "+err.Error()))
		return
	}

	SetSuccessResponse(w, createdClassName)
}

func UpdateSeminarClassName(w http.ResponseWriter, r *http.Request) {
	var className model.SeminarClassName
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&className)
	if err != nil {
		logoped.ErrorLog.Println("unable to retrieve the just parsed code")
		SetErrorResponse(w, NewJSONDecodeError("SeminarClassName"))
		return
	}

	updatedClassName, err := service.SeminarClassService.UpdateSeminarClassName(className)
	if err != nil {
		logoped.ErrorLog.Println("Error updating class name" + err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom ažuriranja naziva časa: "+err.Error()))
		return
	}

	SetSuccessResponse(w, updatedClassName)
}

func ListSeminarClassNames(w http.ResponseWriter, r *http.Request) {
	classNames, err := service.SeminarClassService.GetAllSeminarClassNames()
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste naziva časova: "+err.Error()))
		return
	}

	SetSuccessResponse(w, classNames)
}
