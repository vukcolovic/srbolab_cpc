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

func CreateClient(w http.ResponseWriter, r *http.Request) {
	var client model.Client
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&client)
	if err != nil {
		logoped.ErrorLog.Println("unable to retrieve the just parsed code")
		SetErrorResponse(w, NewJSONDecodeError("Client"))
		return
	}

	createdClient, err := service.ClientService.CreateClient(client)
	if err != nil {
		logoped.ErrorLog.Println("Error creating client " + err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom kreiranja klijenta: "+err.Error()))
		return
	}

	SetSuccessResponse(w, createdClient)
}

func ListClients(w http.ResponseWriter, r *http.Request) {
	var filter model.ClientFilter
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&filter)
	if err != nil {
		logoped.ErrorLog.Println("Error decoding client filter: ", err)
		SetErrorResponse(w, NewJSONDecodeError("ClientFilter"))
		return
	}

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

	users, err := service.ClientService.GetAllClients(skip, take, filter)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste klijenata: "+err.Error()))
		return
	}

	SetSuccessResponse(w, users)
}

func GetClientByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	clientIdParam, ok := vars["id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter id")
		SetErrorResponse(w, NewMissingRequestParamError("id"))
		return
	}

	clientId, err := strconv.Atoi(clientIdParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("clientId", clientIdParam))
		return
	}
	client, err := service.ClientService.GetClientByID(clientId)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja klijenta: "+err.Error()))
		return
	}

	SetSuccessResponse(w, client)
}

func GetClientByJMBG(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	jmbg, ok := vars["jmbg"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter jmbg")
		SetErrorResponse(w, NewMissingRequestParamError("jmbg"))
		return
	}

	client, err := service.ClientService.GetClientByJMBG(jmbg)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja klijenta: "+err.Error()))
		return
	}

	SetSuccessResponse(w, client)
}

func UpdateClient(w http.ResponseWriter, r *http.Request) {
	var client model.Client
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&client)
	if err != nil {
		logoped.ErrorLog.Println("Error decoding client: ", err)
		SetErrorResponse(w, NewJSONDecodeError("Client"))
		return
	}

	updatedClient, err := service.ClientService.UpdateClient(client)
	if err != nil {
		logoped.ErrorLog.Println("Error updating client: ", err)
		SetErrorResponse(w, errors.New("Greška prilikom ažuriranja klijenta: "+err.Error()))
		return
	}

	SetSuccessResponse(w, updatedClient)
}

func DeleteClient(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	clientId, ok := vars["id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter id")
		SetErrorResponse(w, NewMissingRequestParamError("id"))
		return
	}

	clientIdInt, err := strconv.Atoi(clientId)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("clientId", clientId))
		return
	}

	err = service.ClientService.DeleteClient(clientIdInt)
	if err != nil {
		logoped.ErrorLog.Println("error deleting client ", err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom brisanja klijenta: "+err.Error()))
		return
	}

	SetSuccessResponse(w, nil)
}

func CountClients(w http.ResponseWriter, req *http.Request) {
	count, err := service.ClientService.GetClientsCount()
	if err != nil {
		logoped.ErrorLog.Println("error getting client count")
		SetErrorResponse(w, errors.New("Greška prilikom dobijanja broja klijenata: "+err.Error()))
		return
	}

	SetSuccessResponse(w, count)
}
