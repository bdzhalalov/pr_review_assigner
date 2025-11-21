package team

import (
	"github.com/gorilla/mux"
	"net/http"
)

func TeamRouter(h *Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/add", h.CreateTeam).Methods(http.MethodPost)
	router.HandleFunc("/get", h.GetTeam).Methods(http.MethodGet)

	return router
}
