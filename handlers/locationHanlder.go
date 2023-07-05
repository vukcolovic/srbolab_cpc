package handlers

import (
	"errors"
	"net/http"
	"srbolab_cpc/logoped"
	"srbolab_cpc/service"
)

func ListLocations(w http.ResponseWriter, r *http.Request) {
	locations, err := service.LocationService.GetAllLocations()
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste lokacija: "+err.Error()))
		return
	}

	SetSuccessResponse(w, locations)
}
