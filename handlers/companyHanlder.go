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

func CreateCompany(w http.ResponseWriter, r *http.Request) {
	var company model.Company
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&company)
	if err != nil {
		logoped.ErrorLog.Println("unable to retrieve the just parsed code")
		SetErrorResponse(w, NewJSONDecodeError("Company"))
		return
	}

	createdCompany, err := service.CompanyService.CreateCompany(company)
	if err != nil {
		logoped.ErrorLog.Println("Error creating company " + err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom kreiranja firme: "+err.Error()))
		return
	}

	SetSuccessResponse(w, createdCompany)
}

func ListCompanies(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	skipParam, ok := queryParams["skip"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter skip")
		SetErrorResponse(w, NewMissingRequestParamError("skip"))
		return
	}
	skip, err := strconv.Atoi(skipParam[0])
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("skip", skipParam[0]))
		return
	}

	takeParam, ok := queryParams["take"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter take")
		SetErrorResponse(w, NewMissingRequestParamError("take"))
		return
	}
	take, err := strconv.Atoi(takeParam[0])
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("take", takeParam[0]))
		return
	}

	companies, err := service.CompanyService.GetAllCompanies(skip, take)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste firmi: "+err.Error()))
		return
	}

	SetSuccessResponse(w, companies)
}

func GetCompanyByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	companyIDParam, ok := vars["id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter id")
		SetErrorResponse(w, NewMissingRequestParamError("id"))
		return
	}

	companyID, err := strconv.Atoi(companyIDParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("companyID", companyIDParam))
		return
	}
	company, err := service.CompanyService.GetCompanyByID(companyID)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja firme: "+err.Error()))
		return
	}

	SetSuccessResponse(w, company)
}

func UpdateCompany(w http.ResponseWriter, r *http.Request) {
	var company model.Company
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&company)
	if err != nil {
		logoped.ErrorLog.Println("Error decoding company: ", err)
		SetErrorResponse(w, NewJSONDecodeError("Company"))
		return
	}

	updatedCompany, err := service.CompanyService.UpdateCompany(company)
	if err != nil {
		logoped.ErrorLog.Println("Error updating company: ", err)
		SetErrorResponse(w, errors.New("Greška prilikom ažuriranja firme: "+err.Error()))
		return
	}

	SetSuccessResponse(w, updatedCompany)
}

func DeleteCompany(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	companyID, ok := vars["id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter id")
		SetErrorResponse(w, NewMissingRequestParamError("id"))
		return
	}

	companyIDInt, err := strconv.Atoi(companyID)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("companyID", companyID))
		return
	}

	err = service.CompanyService.DeleteCompany(companyIDInt)
	if err != nil {
		logoped.ErrorLog.Println("error deleting company ", err.Error())
		SetErrorResponse(w, errors.New("Greska prilikom brisanja firme: "+err.Error()))
		return
	}

	SetSuccessResponse(w, nil)
}

func CountCompanies(w http.ResponseWriter, req *http.Request) {
	count, err := service.CompanyService.GetCompaniesCount()
	if err != nil {
		logoped.ErrorLog.Println("error getting company count")
		SetErrorResponse(w, errors.New("Greška prilikom dobijanja broja firmi: "+err.Error()))
		return
	}

	SetSuccessResponse(w, count)
}
