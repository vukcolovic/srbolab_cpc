package handlers

import (
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	//var user model.User
	//decoder := json.NewDecoder(r.Body)
	//err := decoder.Decode(&user)
	//if err != nil {
	//	logoped.ErrorLog.Println("unable to retrieve the just parsed code")
	//	util.SetErrorResponse(w, util.NewJSONDecodeError("User"))
	//	return
	//}
	//
	//loginResponse, err := service.UsersService.Login(user)
	//if err != nil {
	//	logoped.ErrorLog.Println("Error loging " + err.Error())
	//	util.SetErrorResponse(w, errors.New("Greska prilikom prijave korisnika: "+err.Error()))
	//	return
	//}

	SetSuccessResponse(w, "")
}
