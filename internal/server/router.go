package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Routers struct {
	teamRouter *mux.Router
	userRouter *mux.Router
	prRouter   *mux.Router
}

func MainRouter(r Routers) *mux.Router {
	router := mux.NewRouter()

	group := router.PathPrefix("/api").Subrouter()
	group.PathPrefix("/teams").Handler(http.StripPrefix("/api/teams", r.teamRouter))
	group.PathPrefix("/users").Handler(http.StripPrefix("/api/users", r.userRouter))
	group.PathPrefix("/pullRequest").Handler(http.StripPrefix("/api/pullRequest", r.prRouter))

	return router
}
