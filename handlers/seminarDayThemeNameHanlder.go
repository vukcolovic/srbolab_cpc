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

func GetSeminarDayThemeNameByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	themeNameIdParam, ok := vars["id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter id")
		SetErrorResponse(w, NewMissingRequestParamError("id"))
		return
	}

	themeNameId, err := strconv.Atoi(themeNameIdParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("themeNameId", themeNameIdParam))
		return
	}
	themeName, err := service.SeminarThemeService.GetSeminarDayThemeNameByID(themeNameId)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja naziva teme: "+err.Error()))
		return
	}

	SetSuccessResponse(w, themeName)
}

func CreateSeminarDayThemeName(w http.ResponseWriter, r *http.Request) {
	var themeName model.SeminarDayThemeName
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&themeName)
	if err != nil {
		logoped.ErrorLog.Println("unable to retrieve the just parsed code")
		SetErrorResponse(w, NewJSONDecodeError("SeminarDayThemeName"))
		return
	}

	createdThemeName, err := service.SeminarThemeService.CreateSeminarDayThemeName(themeName)
	if err != nil {
		logoped.ErrorLog.Println("Error creating theme name" + err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom kreiranja naziva teme: "+err.Error()))
		return
	}

	SetSuccessResponse(w, createdThemeName)
}

func UpdateSeminarDayThemeName(w http.ResponseWriter, r *http.Request) {
	var themeName model.SeminarDayThemeName
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&themeName)
	if err != nil {
		logoped.ErrorLog.Println("unable to retrieve the just parsed code")
		SetErrorResponse(w, NewJSONDecodeError("SeminarDayThemeName"))
		return
	}

	updatedThemeName, err := service.SeminarThemeService.UpdateSeminarDayThemeName(themeName)
	if err != nil {
		logoped.ErrorLog.Println("Error updating theme name" + err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom ažuriranja naziva teme: "+err.Error()))
		return
	}

	SetSuccessResponse(w, updatedThemeName)
}

func ListSeminarDayThemeNames(w http.ResponseWriter, r *http.Request) {
	themeNames, err := service.SeminarThemeService.GetAllSeminarDayThemeNames()
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste naziva tema: "+err.Error()))
		return
	}

	SetSuccessResponse(w, themeNames)
}
