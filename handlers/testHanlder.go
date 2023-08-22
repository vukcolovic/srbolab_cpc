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

func ListTests(w http.ResponseWriter, r *http.Request) {
	tests, err := service.TestService.GetAllTests()
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste testova: "+err.Error()))
		return
	}

	SetSuccessResponse(w, tests)
}

func GetTestByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	testIDParam, ok := vars["id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter id")
		SetErrorResponse(w, NewMissingRequestParamError("id"))
		return
	}

	testID, err := strconv.Atoi(testIDParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("TestID", testIDParam))
		return
	}
	test, err := service.TestService.GetTestByID(testID)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja testa: "+err.Error()))
		return
	}

	SetSuccessResponse(w, test)
}

func CreateTest(w http.ResponseWriter, r *http.Request) {
	var test model.Test
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&test)
	if err != nil {
		logoped.ErrorLog.Println("unable to retrieve the just parsed code")
		SetErrorResponse(w, NewJSONDecodeError("Test"))
		return
	}

	createdTest, err := service.TestService.CreateTest(test)
	if err != nil {
		logoped.ErrorLog.Println("Error creating test " + err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom kreiranja testa: "+err.Error()))
		return
	}

	SetSuccessResponse(w, createdTest)
}

func UpdateTest(w http.ResponseWriter, r *http.Request) {
	var test model.Test
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&test)
	if err != nil {
		logoped.ErrorLog.Println("Error decoding test: ", err)
		SetErrorResponse(w, NewJSONDecodeError("Test"))
		return
	}

	updatedTest, err := service.TestService.UpdateTest(test)
	if err != nil {
		logoped.ErrorLog.Println("Error updating test: ", err)
		SetErrorResponse(w, errors.New("Greška prilikom ažuriranja testa: "+err.Error()))
		return
	}

	SetSuccessResponse(w, updatedTest)
}