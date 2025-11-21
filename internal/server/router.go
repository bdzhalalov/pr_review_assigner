package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Routers struct {
	teamRouter *mux.Router
}

func MainRouter(r Routers) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	group := router.PathPrefix("/api").Subrouter()
	group.PathPrefix("/teams").Handler(http.StripPrefix("/api/teams", r.teamRouter))

	return router
}
