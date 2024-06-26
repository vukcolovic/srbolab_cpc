package server

import (
	"net/http"
	"srbolab_cpc/handlers"
	"srbolab_cpc/middleware"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func RunServer(host string) {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Use(middleware.AuthToken)

	r.HandleFunc("/api/login", handlers.Login).Methods("POST")
	r.HandleFunc("/api/corporate-ip", handlers.IsCorporateIp).Methods("GET")

	s := r.PathPrefix("/api/users").Subrouter()
	s.HandleFunc("/register", handlers.Register).Methods("POST").Name("users_register")
	s.HandleFunc("/update", handlers.UpdateUser).Methods("POST").Name("users_update")
	s.HandleFunc("/list", handlers.ListUsers).Methods("GET").Name("users_list")
	s.HandleFunc("/id/{id}", handlers.GetUserByID).Methods("GET")
	s.HandleFunc("/delete/{id}", handlers.DeleteUser).Methods("GET").Name("users_delete")
	s.HandleFunc("/count", handlers.CountUsers).Methods("GET")
	s.HandleFunc("/teachers", handlers.ListTeachers).Methods("GET")

	s = r.PathPrefix("/api/roles").Subrouter()
	s.HandleFunc("/list", handlers.ListRoles).Methods("GET")
	s.HandleFunc("/id/{id}", handlers.GetRoleByID).Methods("GET")

	s = r.PathPrefix("/api/clients").Subrouter()
	s.HandleFunc("/create", handlers.CreateClient).Methods("POST")
	s.HandleFunc("/create-not-verified", handlers.CreateClientNotVerified).Methods("POST")
	s.HandleFunc("/update", handlers.UpdateClient).Methods("POST")
	s.HandleFunc("/list", handlers.ListClients).Methods("POST")
	s.HandleFunc("/id/{id}", handlers.GetClientByID).Methods("GET")
	s.HandleFunc("/jmbg/{jmbg}", handlers.GetClientByJMBG).Methods("GET")
	s.HandleFunc("/delete/{id}", handlers.DeleteClient).Methods("GET")
	s.HandleFunc("/count", handlers.CountClients).Methods("GET")
	s.HandleFunc("/download/id/{id}/filename/{filename}", handlers.DownloadClientFile).Methods("GET")

	s = r.PathPrefix("/api/seminars").Subrouter()
	s.HandleFunc("/create", handlers.CreateSeminar).Methods("POST")
	s.HandleFunc("/update", handlers.UpdateSeminar).Methods("POST")
	s.HandleFunc("/list/status/{status}", handlers.ListSeminarsByStatus).Methods("GET")
	s.HandleFunc("/list", handlers.ListSeminars).Methods("POST")
	s.HandleFunc("/id/{id}", handlers.GetSeminarByID).Methods("GET")
	s.HandleFunc("/delete/{id}", handlers.DeleteSeminar).Methods("GET")
	s.HandleFunc("/count", handlers.CountSeminars).Methods("GET")
	s.HandleFunc("/download/id/{id}/filename/{filename}", handlers.DownloadSeminarFile).Methods("GET")

	s = r.PathPrefix("/api/client-seminar").Subrouter()
	s.HandleFunc("/update", handlers.UpdateClientSeminar).Methods("POST")
	s.HandleFunc("/insert-bulk", handlers.CreateClientSeminarBulk).Methods("POST")

	s = r.PathPrefix("/api/seminar-days").Subrouter()
	s.HandleFunc("/create-all/{seminar_id}", handlers.CreateAllSeminarDaysForSeminar).Methods("GET")
	s.HandleFunc("/create", handlers.CreateSeminarDay).Methods("POST")
	s.HandleFunc("/update", handlers.UpdateSeminarDay).Methods("POST")
	s.HandleFunc("/list/{seminar_id}", handlers.ListSeminarDays).Methods("GET")
	s.HandleFunc("/id/{id}", handlers.GetSeminarDayByID).Methods("GET")
	s.HandleFunc("/teachers/id/{id}", handlers.GetTeachersFromSeminarDay).Methods("GET")
	s.HandleFunc("/jmbg/{jmbg}", handlers.GetSeminarDayWithTestByJMBG).Methods("GET")
	s.HandleFunc("/download/id/{id}/filename/{filename}", handlers.DownloadSeminarDayFile).Methods("GET")
	//s.HandleFunc("/delete/{id}", handlers.DeleteSeminarDay).Methods("GET")

	s = r.PathPrefix("/api/class-names").Subrouter()
	s.HandleFunc("/create", handlers.CreateSeminarClassName).Methods("POST")
	s.HandleFunc("/update", handlers.UpdateSeminarClassName).Methods("POST")
	s.HandleFunc("/list", handlers.ListSeminarClassNames).Methods("GET")
	s.HandleFunc("/id/{id}", handlers.GetSeminarClassNameByID).Methods("GET")

	s = r.PathPrefix("/api/theme-names").Subrouter()
	s.HandleFunc("/create", handlers.CreateSeminarDayThemeName).Methods("POST")
	s.HandleFunc("/update", handlers.UpdateSeminarDayThemeName).Methods("POST")
	s.HandleFunc("/list", handlers.ListSeminarDayThemeNames).Methods("GET")
	s.HandleFunc("/id/{id}", handlers.GetSeminarDayThemeNameByID).Methods("GET")

	s = r.PathPrefix("/api/companies").Subrouter()
	s.HandleFunc("/create", handlers.CreateCompany).Methods("POST")
	s.HandleFunc("/update", handlers.UpdateCompany).Methods("POST")
	s.HandleFunc("/list", handlers.ListCompanies).Methods("GET")
	s.HandleFunc("/id/{id}", handlers.GetCompanyByID).Methods("GET")
	s.HandleFunc("/delete/{id}", handlers.DeleteCompany).Methods("GET")
	s.HandleFunc("/count", handlers.CountCompanies).Methods("GET")

	s = r.PathPrefix("/api/locations").Subrouter()
	s.HandleFunc("/list", handlers.ListLocations).Methods("GET")
	s.HandleFunc("/create", handlers.CreateLocation).Methods("POST")
	s.HandleFunc("/update", handlers.UpdateLocation).Methods("POST")
	s.HandleFunc("/id/{id}", handlers.GetLocationByID).Methods("GET")
	s.HandleFunc("/count", handlers.CountLocations).Methods("GET")

	s = r.PathPrefix("/api/class-rooms").Subrouter()
	s.HandleFunc("/location/{locationId}", handlers.ListClassRoomsByLocation).Methods("GET")
	s.HandleFunc("/list", handlers.ListClassRooms).Methods("GET")
	s.HandleFunc("/create", handlers.CreateClassRoom).Methods("POST")
	s.HandleFunc("/update", handlers.UpdateClassRoom).Methods("POST")
	s.HandleFunc("/id/{id}", handlers.GetClassRoomByID).Methods("GET")

	s = r.PathPrefix("/api/seminar-types").Subrouter()
	s.HandleFunc("/list", handlers.ListBaseSeminarTypes).Methods("GET")
	s.HandleFunc("/themes/seminar-type/{seminarTypeId}", handlers.ListSeminarThemesBySeminarType).Methods("GET")
	s.HandleFunc("/themes/list", handlers.ListSeminarThemes).Methods("GET")

	s = r.PathPrefix("/api/seminar-statuses").Subrouter()
	s.HandleFunc("/list", handlers.ListSeminarStatuses).Methods("GET")

	s = r.PathPrefix("/api/print").Subrouter()
	s.HandleFunc("/seminar/student-list/{seminar_id}", handlers.PrintSeminarStudentList).Methods("GET")
	s.HandleFunc("/seminar/confirmation/{seminar_id}", handlers.PrintConfirmations).Methods("GET")
	s.HandleFunc("/seminar/confirmation-receive/{seminar_id}", handlers.PrintConfirmationReceives).Methods("GET")
	s.HandleFunc("/seminar/muster/{seminar_day_id}", handlers.PrintMuster).Methods("GET")
	s.HandleFunc("/seminar/check-in/{seminar_id}", handlers.PrintCheckIn).Methods("GET")
	s.HandleFunc("/seminar/teacher-evidence/{seminar_day_id}", handlers.PrintSeminarEvidence).Methods("GET")
	s.HandleFunc("/seminar-day/test/barcode", handlers.PrintTestBarcode).Methods("GET")
	s.HandleFunc("/seminar-day/training-realization/{seminar_day_id}", handlers.PrintPlanTreningRealization).Methods("GET")
	s.HandleFunc("/seminar/payments/{seminar_id}", handlers.PrintPayments).Methods("GET")
	s.HandleFunc("/seminar/report/{seminar_id}", handlers.PrintReport).Methods("GET")
	s.HandleFunc("/seminar/report2/{seminar_id}", handlers.PrintReport2).Methods("GET")
	s.HandleFunc("/seminar/exam-registration/{seminar_id}", handlers.PrintExamRegistration).Methods("GET")
	s.HandleFunc("/test/{test_id}", handlers.PrintTest).Methods("GET")

	s = r.PathPrefix("/api/questions").Subrouter()
	s.HandleFunc("/list/seminar-theme/{id}", handlers.ListQuestionsBySeminarThemeID).Methods("GET")
	s.HandleFunc("/list", handlers.ListQuestions).Methods("GET")
	s.HandleFunc("/create", handlers.CreateQuestion).Methods("POST")
	s.HandleFunc("/update", handlers.UpdateQuestion).Methods("POST")
	s.HandleFunc("/id/{id}", handlers.GetQuestionByID).Methods("GET")

	s = r.PathPrefix("/api/tests").Subrouter()
	s.HandleFunc("/list", handlers.ListTests).Methods("GET")
	s.HandleFunc("/list/seminar-theme/{id}", handlers.ListTestsBySeminarThemeID).Methods("GET")
	s.HandleFunc("/create", handlers.CreateTest).Methods("POST")
	s.HandleFunc("/update", handlers.UpdateTest).Methods("POST")
	s.HandleFunc("/id/{id}", handlers.GetTestByID).Methods("GET")
	s.HandleFunc("/id/{id}/seminar-day/{seminar-day}/client-jmbg/{jmbg}", handlers.GetClientTestsBySeminarDayAndJMBG).Methods("GET")
	s.HandleFunc("/client-test/create", handlers.SaveClientTest).Methods("POST")
	s.HandleFunc("/client-test/seminar-day/{seminar-day}", handlers.GetClientTestsBySeminarDay).Methods("GET")

	s = r.PathPrefix("/api/excel").Subrouter()
	s.HandleFunc("/client-tests/{seminar-day}", handlers.PrintClientTestsBySeminarDay).Methods("GET")
	s.HandleFunc("/clients", handlers.PrintListOfCients).Methods("GET")
	s.HandleFunc("/list_trainees/{seminar-day}", handlers.PrintListTraineesBySeminarDay).Methods("GET")
	s.HandleFunc("/seminars-report/clients", handlers.PrintSeminarsReportOfClients).Methods("POST")
	s.HandleFunc("/seminars-report/teachers", handlers.PrintSeminarsReportOfTeachers).Methods("POST")

	s = r.PathPrefix("/api/survey-questions").Subrouter()
	s.HandleFunc("/list", handlers.ListSurveyQuestions).Methods("GET")
	s.HandleFunc("/create", handlers.CreateSurveyQuestion).Methods("POST")
	s.HandleFunc("/id/{id}", handlers.GetSurveyQuestionByID).Methods("GET")

	s = r.PathPrefix("/api/surveys").Subrouter()
	s.HandleFunc("/list", handlers.ListSurveys).Methods("GET")
	s.HandleFunc("/create", handlers.CreateSurvey).Methods("POST")
	s.HandleFunc("/id/{id}", handlers.GetSurveyByID).Methods("GET")
	s.HandleFunc("/active", handlers.GetActiveSurveys).Methods("GET")
	s.HandleFunc("/client-survey/create", handlers.SaveClientSurvey).Methods("POST")

	s = r.PathPrefix("/api/partners").Subrouter()
	s.HandleFunc("/list", handlers.ListPartners).Methods("GET")
	s.HandleFunc("/create", handlers.CreatePartner).Methods("POST")
	s.HandleFunc("/update", handlers.UpdatePartner).Methods("POST")
	s.HandleFunc("/id/{id}", handlers.GetPartnerByID).Methods("GET")
	s.HandleFunc("/count", handlers.CountPartners).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // All origins
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
	})

	http.ListenAndServe(host, c.Handler(r))
}
