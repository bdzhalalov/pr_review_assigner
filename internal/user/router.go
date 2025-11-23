package user

import (
	"github.com/gorilla/mux"
	"net/http"
)

func UserRouter(h *Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/setIsActive", h.UpdateActivity).Methods(http.MethodPost)
	router.HandleFunc("/getReview", h.GetUserReviews).Methods(http.MethodGet)

	return router
}
