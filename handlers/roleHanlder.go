package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"srbolab_cpc/logoped"
	"srbolab_cpc/service"
	"strconv"
)

func ListRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := service.RoleService.GetAllRoles()
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja liste rola: "+err.Error()))
		return
	}

	SetSuccessResponse(w, roles)
}

func GetRoleByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	roleIdParam, ok := vars["id"]
	if !ok {
		logoped.ErrorLog.Println("missing parameter id")
		SetErrorResponse(w, NewMissingRequestParamError("id"))
		return
	}

	roleId, err := strconv.Atoi(roleIdParam)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, NewWrongParamFormatErrorError("roleId", roleIdParam))
		return
	}
	role, err := service.RoleService.GetRoleByID(roleId)
	if err != nil {
		logoped.ErrorLog.Println(err.Error())
		SetErrorResponse(w, errors.New("Greška prilikom povlačenja role: "+err.Error()))
		return
	}

	SetSuccessResponse(w, role)
}
