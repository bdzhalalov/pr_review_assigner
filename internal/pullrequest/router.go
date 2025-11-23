package pullrequest

import (
	"github.com/gorilla/mux"
	"net/http"
)

func PullRequestRouter(h *Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/create", h.CreatePullRequest).Methods(http.MethodPost)
	router.HandleFunc("/merge", h.UpdatePullRequestStatus).Methods(http.MethodPost)
	router.HandleFunc("/reassign", h.ReassignPullRequestReviewer).Methods(http.MethodPost)

	return router
}
