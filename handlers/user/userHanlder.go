package user

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user model.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		loger.ErrorLog.Println("unable to retrieve the just parsed code")
		SetErrorResponse(w, util.NewJSONDecodeError("User"))
		return
	}

	createdUser, err := service.UsersService.CreateUser(user)
	if err != nil {
		loger.ErrorLog.Println("Error creating user " + err.Error())
		SetErrorResponse(w, errors.New("Greska prilikom registracije korisnika: "+err.Error()))
		return
	}

	SetSuccessResponse(w, createdUser)
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	skipParam, ok := queryParams["skip"]
	if !ok {
		loger.ErrorLog.Println("missing parameter skip")
		SetErrorResponse(w, util.NewMissingRequestParamError("skip"))
		return
	}
	skip, err := strconv.Atoi(skipParam[0])
	if err != nil {
		loger.ErrorLog.Println(err.Error())
		SetErrorResponse(w, util.NewWrongParamFormatErrorError("skip", skipParam[0]))
		return
	}

	takeParam, ok := queryParams["take"]
	if !ok {
		loger.ErrorLog.Println("missing parameter take")
		SetErrorResponse(w, util.NewMissingRequestParamError("take"))
		return
	}
	take, err := strconv.Atoi(takeParam[0])
	if err != nil {
		loger.ErrorLog.Println(err.Error())
		SetErrorResponse(w, util.NewWrongParamFormatErrorError("take", takeParam[0]))
		return
	}

	users, err := service.UsersService.GetAllUsers(skip, take)
	if err != nil {
		loger.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greska prilikom povlacenja liste korisnika: "+err.Error()))
		return
	}

	SetSuccessResponse(w, users)
}

func GetUserByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userIdParam, ok := vars["id"]
	if !ok {
		loger.ErrorLog.Println("missing parameter id")
		SetErrorResponse(w, util.NewMissingRequestParamError("id"))
		return
	}

	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		loger.ErrorLog.Println(err.Error())
		SetErrorResponse(w, util.NewWrongParamFormatErrorError("userId", userIdParam))
		return
	}
	user, err := service.UsersService.GetUserByID(userId)
	if err != nil {
		loger.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greska prilikom povlacenja korisnika: "+err.Error()))
		return
	}

	SetSuccessResponse(w, user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		loger.ErrorLog.Println("Error decoding user: ", err)
		SetErrorResponse(w, util.NewJSONDecodeError("User"))
		return
	}

	updatedUser, err := service.UsersService.UpdateUser(user)
	if err != nil {
		loger.ErrorLog.Println("Error updating user: ", err)
		SetErrorResponse(w, errors.New("Greska prilikom azuriranja korisnika: "+err.Error()))
		return
	}

	SetSuccessResponse(w, updatedUser)
}

func DeleteUser(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userId, ok := vars["id"]
	if !ok {
		loger.ErrorLog.Println("missing parameter id")
		SetErrorResponse(w, util.NewMissingRequestParamError("id"))
		return
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		loger.ErrorLog.Println(err.Error())
		SetErrorResponse(w, util.NewWrongParamFormatErrorError("userId", userId))
		return
	}

	err = service.UsersService.DeleteUser(userIdInt)
	if err != nil {
		loger.ErrorLog.Println("error deleting user ", err.Error())
		SetErrorResponse(w, errors.New("Greska prilikom brisanja korisnika: "+err.Error()))
		return
	}

	SetSuccessResponse(w, nil)
}

func CountUsers(w http.ResponseWriter, req *http.Request) {
	count, err := service.UsersService.GetUsersCount()
	if err != nil {
		loger.ErrorLog.Println("error getting user count")
		SetErrorResponse(w, errors.New("Greska prilikom dobijanja broja korisnika: "+err.Error()))
		return
	}

	SetSuccessResponse(w, count)
}
