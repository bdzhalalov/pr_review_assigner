package user

import (
	"encoding/json"
	app "github.com/bdzhalalov/pr-review-assigner/internal/app/user"
	"github.com/bdzhalalov/pr-review-assigner/internal/user/dto"
	"github.com/bdzhalalov/pr-review-assigner/pkg/render"
	"net/http"
)

type Handler struct {
	service *Service
	app     *app.UserApp
}

func NewUserHandler(service *Service, app *app.UserApp) *Handler {
	return &Handler{
		service: service,
		app:     app,
	}
}

func (h *Handler) UpdateActivity(w http.ResponseWriter, r *http.Request) {
	var input dto.UpdateUserActivityDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		render.RenderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedUser, err := h.service.UpdateUserActivity(input)
	if err != nil {
		render.RenderJSON(w, err.Message, err.Code)
		return
	}

	render.RenderJSON(w, updatedUser, http.StatusOK)
}

func (h *Handler) GetUserReviews(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")
	if userId == "" {
		render.RenderJSON(w, "User id is required", http.StatusBadRequest)
		return
	}

	input := app.GetPrReviewsRequest{
		ReviewerID: userId,
	}

	output, err := h.app.GetUserPullRequestReviews(input)
	if err != nil {
		render.RenderJSON(w, err.Message, err.Code)
		return
	}

	render.RenderJSON(w, output, http.StatusOK)
}
