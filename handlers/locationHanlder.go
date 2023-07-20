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

func ListLocations(w http.ResponseWriter, r *http.Request) {
	locations, err := service.LocationService.GetAllLocations()
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste lokacija: "+err.Error()))
		return
	}

	SetSuccessResponse(w, locations)
}

func GetLocationByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	locationIDParam, ok := vars["id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter id")
		SetErrorResponse(w, NewMissingRequestParamError("id"))
		return
	}

	locationID, err := strconv.Atoi(locationIDParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("locationID", locationIDParam))
		return
	}
	location, err := service.LocationService.GetLocationByID(locationID)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja lokacije: "+err.Error()))
		return
	}

	SetSuccessResponse(w, location)
}

func CreateLocation(w http.ResponseWriter, r *http.Request) {
	var location model.Location
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&location)
	if err != nil {
		logoped.ErrorLog.Println("unable to retrieve the just parsed code")
		SetErrorResponse(w, NewJSONDecodeError("Location"))
		return
	}

	createdLocation, err := service.LocationService.CreateLocation(location)
	if err != nil {
		logoped.ErrorLog.Println("Error creating location " + err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom kreiranja lokacije: "+err.Error()))
		return
	}

	SetSuccessResponse(w, createdLocation)
}

func UpdateLocation(w http.ResponseWriter, r *http.Request) {
	var location model.Location
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&location)
	if err != nil {
		logoped.ErrorLog.Println("Error decoding location: ", err)
		SetErrorResponse(w, NewJSONDecodeError("Location"))
		return
	}

	updatedLocation, err := service.LocationService.UpdateLocation(location)
	if err != nil {
		logoped.ErrorLog.Println("Error updating location: ", err)
		SetErrorResponse(w, errors.New("Greška prilikom ažuriranja lokacije: "+err.Error()))
		return
	}

	SetSuccessResponse(w, updatedLocation)
}

func CountLocations(w http.ResponseWriter, req *http.Request) {
	count, err := service.LocationService.GetLocationsCount()
	if err != nil {
		logoped.ErrorLog.Println("error getting locations count")
		SetErrorResponse(w, errors.New("Greška prilikom dobijanja broja lokacija: "+err.Error()))
		return
	}

	SetSuccessResponse(w, count)
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
