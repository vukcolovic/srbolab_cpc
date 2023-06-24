package server

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"srbolab_cpc/handlers/authhandler"
	"srbolab_cpc/handlers/user"
)

func RunServer(host string) {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.HandleFunc("/login", authhandler.Login).Methods("POST")

	s := r.PathPrefix("/api/users").Subrouter()
	s.HandleFunc("/register", user.Register).Methods("POST")
	s.HandleFunc("/update", user.UpdateUser).Methods("POST")
	s.HandleFunc("/list", user.ListUsers).Methods("GET")
	s.HandleFunc("/id/{id}", user.GetUserByID).Methods("GET")
	s.HandleFunc("/delete/{id}", user.DeleteUser).Methods("GET")
	s.HandleFunc("/count", user.CountUsers).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // All origins
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
	})

	http.ListenAndServe(host, c.Handler(r))
}
