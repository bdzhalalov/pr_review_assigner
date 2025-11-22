package user

import (
	"github.com/gorilla/mux"
	"net/http"
)

func UserRouter(h *Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/setIsActive", h.UpdateActivity).Methods(http.MethodPost)

	return router
}
