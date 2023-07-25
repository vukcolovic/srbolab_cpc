package server

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"srbolab_cpc/handlers"
)

func RunServer(host string) {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.HandleFunc("/login", handlers.Login).Methods("POST")

	s := r.PathPrefix("/api/users").Subrouter()
	s.HandleFunc("/register", handlers.Register).Methods("POST")
	s.HandleFunc("/update", handlers.UpdateUser).Methods("POST")
	s.HandleFunc("/list", handlers.ListUsers).Methods("GET")
	s.HandleFunc("/id/{id}", handlers.GetUserByID).Methods("GET")
	s.HandleFunc("/delete/{id}", handlers.DeleteUser).Methods("GET")
	s.HandleFunc("/count", handlers.CountUsers).Methods("GET")

	s = r.PathPrefix("/api/clients").Subrouter()
	s.HandleFunc("/create", handlers.CreateClient).Methods("POST")
	s.HandleFunc("/update", handlers.UpdateClient).Methods("POST")
	s.HandleFunc("/list", handlers.ListClients).Methods("POST")
	s.HandleFunc("/id/{id}", handlers.GetClientByID).Methods("GET")
	s.HandleFunc("/jmbg/{jmbg}", handlers.GetClientByJMBG).Methods("GET")
	s.HandleFunc("/delete/{id}", handlers.DeleteClient).Methods("GET")
	s.HandleFunc("/count", handlers.CountClients).Methods("GET")

	s = r.PathPrefix("/api/seminars").Subrouter()
	s.HandleFunc("/create", handlers.CreateSeminar).Methods("POST")
	s.HandleFunc("/update", handlers.UpdateSeminar).Methods("POST")
	s.HandleFunc("/list/status/{status}", handlers.ListSeminarsByStatus).Methods("GET")
	s.HandleFunc("/list", handlers.ListSeminars).Methods("GET")
	s.HandleFunc("/id/{id}", handlers.GetSeminarByID).Methods("GET")
	s.HandleFunc("/delete/{id}", handlers.DeleteSeminar).Methods("GET")
	s.HandleFunc("/count", handlers.CountSeminars).Methods("GET")

	s = r.PathPrefix("/api/seminar-days").Subrouter()
	s.HandleFunc("/create-all/{seminar_id}", handlers.CreateAllSeminarDaysForSeminar).Methods("GET")
	s.HandleFunc("/create", handlers.CreateSeminarDay).Methods("POST")
	s.HandleFunc("/update", handlers.UpdateSeminarDay).Methods("POST")
	s.HandleFunc("/list/{seminar_id}", handlers.ListSeminarDays).Methods("GET")
	s.HandleFunc("/id/{id}", handlers.GetSeminarDayByID).Methods("GET")
	//s.HandleFunc("/delete/{id}", handlers.DeleteSeminarDay).Methods("GET")

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

	s = r.PathPrefix("/api/seminar-statuses").Subrouter()
	s.HandleFunc("/list", handlers.ListSeminarStatuses).Methods("GET")

	s = r.PathPrefix("/api/print").Subrouter()
	s.HandleFunc("/seminar/student-list/{seminar_id}", handlers.PrintSeminarStudentList).Methods("GET")
	s.HandleFunc("/seminar/confirmation-statement/{seminar_id}", handlers.PrintConfirmationStatements).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // All origins
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
	})

	http.ListenAndServe(host, c.Handler(r))
}
