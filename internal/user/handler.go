package user

import (
	"encoding/json"
	"github.com/bdzhalalov/pr-review-assigner/internal/user/dto"
	"github.com/bdzhalalov/pr-review-assigner/pkg/render"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewUserHandler(service *Service) *Handler {
	return &Handler{
		service: service,
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
