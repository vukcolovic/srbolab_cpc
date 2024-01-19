package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"srbolab_cpc/logoped"
	"srbolab_cpc/model"
	"srbolab_cpc/service"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func GetSeminarDayByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	seminarDayIdParam, ok := vars["id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter id")
		SetErrorResponse(w, NewMissingRequestParamError("id"))
		return
	}

	seminarDayId, err := strconv.Atoi(seminarDayIdParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("seminarDayId", seminarDayIdParam))
		return
	}
	seminarDay, err := service.SeminarDayService.GetSeminarDayByID(seminarDayId)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja seminar dana: "+err.Error()))
		return
	}

	SetSuccessResponse(w, seminarDay)
}

func CreateSeminarDay(w http.ResponseWriter, r *http.Request) {
	var day model.SeminarDay
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&day)
	if err != nil {
		logoped.ErrorLog.Println("unable to retrieve the just parsed code")
		SetErrorResponse(w, NewJSONDecodeError("SeminarDay"))
		return
	}

	createdDay, err := service.SeminarDayService.CreateSeminarDay(day)
	if err != nil {
		logoped.ErrorLog.Println("Error creating seminar day" + err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom kreiranja seminar dana: "+err.Error()))
		return
	}

	SetSuccessResponse(w, createdDay)
}

func UpdateSeminarDay(w http.ResponseWriter, r *http.Request) {
	var day model.SeminarDay
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&day)
	if err != nil {
		logoped.ErrorLog.Println("unable to retrieve the just parsed code")
		SetErrorResponse(w, NewJSONDecodeError("SeminarDay"))
		return
	}

	updatedDay, err := service.SeminarDayService.UpdateSeminarDay(day)
	if err != nil {
		logoped.ErrorLog.Println("Error updating seminar day" + err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom ažuriranja seminar dana: "+err.Error()))
		return
	}

	SetSuccessResponse(w, updatedDay)
}

func ListSeminarDays(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seminarIdParam, ok := vars["seminar_id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter seminar_id")
		SetErrorResponse(w, NewMissingRequestParamError("seminar_id"))
		return
	}
	seminarID, err := strconv.Atoi(seminarIdParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("seminar_id", seminarIdParam))
		return
	}

	seminarDays, err := service.SeminarDayService.GetSeminarDaysBySeminarID(seminarID)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste seminar dana: "+err.Error()))
		return
	}

	SetSuccessResponse(w, seminarDays)
}

func CreateAllSeminarDaysForSeminar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seminarIdParam, ok := vars["seminar_id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter seminar_id")
		SetErrorResponse(w, NewMissingRequestParamError("seminar_id"))
		return
	}
	seminarID, err := strconv.Atoi(seminarIdParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("seminar_id", seminarIdParam))
		return
	}

	seminarDays, err := service.SeminarDayService.CreateAllSeminarDaysForSeminar(seminarID)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom kreiranja svih seminar dana za seminar: "+err.Error()))
		return
	}

	SetSuccessResponse(w, seminarDays)
}

func GetSeminarDayWithTestByJMBG(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	jmbg, ok := vars["jmbg"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter jmbg")
		SetErrorResponse(w, NewMissingRequestParamError("jmbg"))
		return
	}

	client, err := service.ClientService.GetClientByJMBGWithSeminars(jmbg)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja klijenta: "+err.Error()))
		return
	}

	seminarDays := []model.SeminarDay{}
	for _, cs := range client.Seminars {
		if cs.Seminar.SeminarStatusID == model.SEMINAR_STATUS_IN_PROGRESS {
			seminarDays = append(seminarDays, cs.Seminar.Days...)
		}
	}

	if len(seminarDays) == 0 {
		SetErrorResponse(w, errors.New("Vozač sa upisanim jmbg-om nije na nekom od aktulenih seminara."))
		return
	}

	for _, day := range seminarDays {
		if day.Date.Day() != time.Now().Day() || day.Date.Month() != time.Now().Month() || day.Date.Year() != time.Now().Year() {
			continue
		}
		if day.TestID == nil || *day.TestID == 0 {
			SetErrorResponse(w, errors.New("Seminar dan nema izabran test."))
			return
		}

		tests, err := service.TestService.GetClientTestBySeminarDayIDAndJMBG(int(day.ID), jmbg)
		if err != nil {
			SetErrorResponse(w, err)
			return
		}

		if len(tests) > 1 {
			SetErrorResponse(w, errors.New("Ovaj test nije dozvoljen, klijent je već odradio dva testa u toku dana."))
			return
		}

		if len(tests) == 1 {
			if !tests[0].CreatedAt.Add(2 * time.Hour).Before(time.Now()) {
				SetErrorResponse(w, errors.New("Nije dozvoljeno snimiti test, rađen je skoro."))
				return
			}
		}

		fullSeminarDay, err := service.SeminarDayService.GetSeminarDayWithTestByID(int(day.ID))
		if err != nil {
			SetErrorResponse(w, err)
			return
		}

		SetSuccessResponse(w, fullSeminarDay)
		return
	}

	SetErrorResponse(w, errors.New("Danas nije predviđen dan seminara."))
}

func DownloadSeminarDayFile(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	seminarDayId, ok := vars["id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter id")
		SetErrorResponse(w, NewMissingRequestParamError("id"))
		return
	}

	filename, ok := vars["filename"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter filename")
		SetErrorResponse(w, NewMissingRequestParamError("filename"))
		return
	}

	content, err := service.FileService.GetFile(service.SeminarDayFolder, seminarDayId, filename)
	if err != nil {
		logoped.ErrorLog.Println("error downloading file")
		SetErrorResponse(w, errors.New("Greška prilikom skidanja fajla: "+err.Error()))
		return
	}

	SetSuccessResponse(w, content)
}

func GetTeachersFromSeminarDay(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	seminarDayIdParam, ok := vars["id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter id")
		SetErrorResponse(w, NewMissingRequestParamError("id"))
		return
	}

	seminarDayId, err := strconv.Atoi(seminarDayIdParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("seminarDayId", seminarDayIdParam))
		return
	}

	teachers, err := service.SeminarDayService.GetTeachersFromSeminarDay(seminarDayId)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja predavača: "+err.Error()))
		return
	}

	SetSuccessResponse(w, teachers)
}
