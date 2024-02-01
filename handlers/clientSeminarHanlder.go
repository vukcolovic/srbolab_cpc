package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"srbolab_cpc/logoped"
	"srbolab_cpc/model"
	"srbolab_cpc/service"
)

func UpdateClientSeminar(w http.ResponseWriter, r *http.Request) {
	var cs model.ClientSeminar
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cs)
	if err != nil {
		logoped.ErrorLog.Println("Error decoding client seminar: ", err)
		SetErrorResponse(w, NewJSONDecodeError("clientSeminar"))
		return
	}

	updatedClientSeminar, err := service.ClientSeminarService.UpdateClientSeminar(cs)
	if err != nil {
		logoped.ErrorLog.Println("Error updating client seminar: ", err)
		SetErrorResponse(w, errors.New("Greška prilikom ažuriranja klijent seminara: "+err.Error()))
		return
	}

	SetSuccessResponse(w, updatedClientSeminar)
}

type CreateBulkRequest struct {
	SeminarID int   `json:"seminar_id"`
	ClientIDs []int `json:"client_ids"`
}

func CreateClientSeminarBulk(w http.ResponseWriter, r *http.Request) {
	var req CreateBulkRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		logoped.ErrorLog.Println("Error decoding client seminar buld: ", err)
		SetErrorResponse(w, NewJSONDecodeError("CreateBulkRequest"))
		return
	}

	err = service.ClientSeminarService.CreateBulk(req.SeminarID, req.ClientIDs)
	if err != nil {
		logoped.ErrorLog.Println("Error updating client seminar: ", err)
		SetErrorResponse(w, errors.New("Greška prilikom ažuriranja klijent seminara: "+err.Error()))
		return
	}

	SetSuccessResponse(w, nil)
}
