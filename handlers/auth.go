package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"srbolab_cpc/logoped"
	"srbolab_cpc/model"
	"srbolab_cpc/service"
	"srbolab_cpc/util"
)

var CorporateIps []string

func Login(w http.ResponseWriter, r *http.Request) {
	var request model.LoginRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		logoped.ErrorLog.Println("unable to retrieve the just parsed code")
		SetErrorResponse(w, NewJSONDecodeError("LoginRequest"))
		return
	}

	loginResponse, err := service.AuthService.Login(request)
	if err != nil {
		logoped.ErrorLog.Println("Error loging " + err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom prijave korisnika: "+err.Error()))
		return
	}

	SetSuccessResponse(w, loginResponse)
}

type Corporate struct {
	IsCorporate bool `json:"is_corporate"`
}

func IsCorporateIp(w http.ResponseWriter, r *http.Request) {
	ip, _ := GetIP(r)

	isCorporate := util.Contains(CorporateIps, ip)

	SetSuccessResponse(w, Corporate{IsCorporate: isCorporate})
}
