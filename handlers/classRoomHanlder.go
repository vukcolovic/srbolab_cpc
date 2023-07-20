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

	classRooms, err := service.ClassRoomService.GetClassRoomsByLocationID(locationId)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste učionica: "+err.Error()))
		return
	}

	SetSuccessResponse(w, classRooms)
}

func ListClassRooms(w http.ResponseWriter, r *http.Request) {
	classRooms, err := service.ClassRoomService.GetAllClassRooms()
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste učionica: "+err.Error()))
		return
	}

	SetSuccessResponse(w, classRooms)
}

func GetClassRoomByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	classRoomIDParam, ok := vars["id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter id")
		SetErrorResponse(w, NewMissingRequestParamError("id"))
		return
	}

	classRoomID, err := strconv.Atoi(classRoomIDParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("classRoomID", classRoomIDParam))
		return
	}
	classRoom, err := service.ClassRoomService.GetClassRoomByID(classRoomID)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja učionice: "+err.Error()))
		return
	}

	SetSuccessResponse(w, classRoom)
}

func CreateClassRoom(w http.ResponseWriter, r *http.Request) {
	var classRoom model.ClassRoom
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&classRoom)
	if err != nil {
		logoped.ErrorLog.Println("unable to retrieve the just parsed code")
		SetErrorResponse(w, NewJSONDecodeError("ClassRoom"))
		return
	}

	createdClassRoom, err := service.ClassRoomService.CreateClassRoom(classRoom)
	if err != nil {
		logoped.ErrorLog.Println("Error creating class room " + err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom kreiranja učionice: "+err.Error()))
		return
	}

	SetSuccessResponse(w, createdClassRoom)
}

func UpdateClassRoom(w http.ResponseWriter, r *http.Request) {
	var classRoom model.ClassRoom
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&classRoom)
	if err != nil {
		logoped.ErrorLog.Println("Error decoding class room: ", err)
		SetErrorResponse(w, NewJSONDecodeError("ClassRoom"))
		return
	}

	updatedClassRoom, err := service.ClassRoomService.UpdateClassRoom(classRoom)
	if err != nil {
		logoped.ErrorLog.Println("Error updating class room: ", err)
		SetErrorResponse(w, errors.New("Greška prilikom ažuriranja učionice: "+err.Error()))
		return
	}

	SetSuccessResponse(w, updatedClassRoom)
}
