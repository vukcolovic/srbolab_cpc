package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"srbolab_cpc/logoped"
	"srbolab_cpc/service"
	"strconv"
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

func ListClassRoomsByLocation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	locationIdParam, ok := params["locationId"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter skip")
		SetErrorResponse(w, NewMissingRequestParamError("locationId"))
		return
	}
	locationId, err := strconv.Atoi(locationIdParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("locationId", locationIdParam))
		return
	}

	classRooms, err := service.LocationService.GetClassRoomsByLocationID(locationId)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste učionica: "+err.Error()))
		return
	}

	SetSuccessResponse(w, classRooms)
}
