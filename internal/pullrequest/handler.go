package pullrequest

import (
	"encoding/json"
	app "github.com/bdzhalalov/pr-review-assigner/internal/app/pullrequest"
	"github.com/bdzhalalov/pr-review-assigner/internal/pullrequest/dto"
	"github.com/bdzhalalov/pr-review-assigner/pkg/render"
	"net/http"
)

type Handler struct {
	app     *app.PRApp
	service *Service
}

func NewPrHandler(app *app.PRApp, service *Service) *Handler {
	return &Handler{
		app:     app,
		service: service,
	}
}

func (h *Handler) CreatePullRequest(w http.ResponseWriter, r *http.Request) {
	var input app.AddPrRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		render.RenderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	pr, err := h.app.CreateApp.CreatePullRequest(input)
	if err != nil {
		render.RenderJSON(w, err.Message, err.Code)
		return
	}

	render.RenderJSON(w, pr, http.StatusCreated)
}

func (h *Handler) UpdatePullRequestStatus(w http.ResponseWriter, r *http.Request) {
	var input dto.PrUpdateStatusDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		render.RenderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	pr, err := h.service.UpdatePullRequestStatus(input)
	if err != nil {
		render.RenderJSON(w, err.Message, err.Code)
		return
	}

	render.RenderJSON(w, pr, http.StatusOK)
}

func (h *Handler) ReassignPullRequestReviewer(w http.ResponseWriter, r *http.Request) {
	var input app.ReassignPrRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		render.RenderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.app.ReassignApp.ReassignPullRequest(input)
	if err != nil {
		render.RenderJSON(w, err.Message, err.Code)
		return
	}

	render.RenderJSON(w, response, http.StatusOK)
}
