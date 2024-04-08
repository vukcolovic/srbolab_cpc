package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"srbolab_cpc/logoped"
	"srbolab_cpc/model"
	"srbolab_cpc/service"
	"strconv"

	"github.com/gorilla/mux"
)

func CreatePartner(w http.ResponseWriter, r *http.Request) {
	var partner model.Partner
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&partner)
	if err != nil {
		logoped.ErrorLog.Println("unable to retrieve the just parsed code")
		SetErrorResponse(w, NewJSONDecodeError("Partner"))
		return
	}

	createdPartner, err := service.PartnerService.CreatePartner(partner)
	if err != nil {
		logoped.ErrorLog.Println("Error creating partner " + err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom kreiranja partnera: "+err.Error()))
		return
	}

	SetSuccessResponse(w, createdPartner)
}

func ListPartners(w http.ResponseWriter, r *http.Request) {
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

	partners, err := service.PartnerService.GetAllPartners(skip, take)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste partnera: "+err.Error()))
		return
	}

	SetSuccessResponse(w, partners)
}

func GetPartnerByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	partnerIdParam, ok := vars["id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter id")
		SetErrorResponse(w, NewMissingRequestParamError("id"))
		return
	}

	partnerId, err := strconv.Atoi(partnerIdParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("partnerId", partnerIdParam))
		return
	}
	partner, err := service.PartnerService.GetPartnerByID(partnerId)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja partnera: "+err.Error()))
		return
	}

	SetSuccessResponse(w, partner)
}

func UpdatePartner(w http.ResponseWriter, r *http.Request) {
	var partner model.Partner
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&partner)
	if err != nil {
		logoped.ErrorLog.Println("Error decoding partner: ", err)
		SetErrorResponse(w, NewJSONDecodeError("Partner"))
		return
	}

	updatedPartner, err := service.PartnerService.UpdatePartner(partner)
	if err != nil {
		logoped.ErrorLog.Println("Error updating partner: ", err)
		SetErrorResponse(w, errors.New("Greška prilikom ažuriranja partnera: "+err.Error()))
		return
	}

	SetSuccessResponse(w, updatedPartner)
}

func DeletePartner(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	partnerId, ok := vars["id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter id")
		SetErrorResponse(w, NewMissingRequestParamError("id"))
		return
	}

	partnerIdInt, err := strconv.Atoi(partnerId)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("partnerId", partnerId))
		return
	}

	err = service.PartnerService.DeletePartner(partnerIdInt)
	if err != nil {
		logoped.ErrorLog.Println("error deleting partner ", err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom brisanja partnera: "+err.Error()))
		return
	}

	SetSuccessResponse(w, nil)
}

func CountPartners(w http.ResponseWriter, req *http.Request) {
	count, err := service.PartnerService.GetPartnersCount()
	if err != nil {
		logoped.ErrorLog.Println("error getting partner count")
		SetErrorResponse(w, errors.New("Greška prilikom dobijanja broja partnera: "+err.Error()))
		return
	}

	SetSuccessResponse(w, count)
}
