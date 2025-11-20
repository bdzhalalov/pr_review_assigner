package team

import (
	"encoding/json"
	"github.com/bdzhalalov/pr-review-assigner/internal/app"
	"github.com/bdzhalalov/pr-review-assigner/pkg/render"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	app     *app.TeamApp
	service *Service
	logger  *logrus.Logger
}

func (h *Handler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var input AddTeamDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		render.RenderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	team, err := h.app.AddTeam(input)
	if err != nil {
		render.RenderJSON(w, err.Message, err.Code)
		return
	}

	render.RenderJSON(w, team, http.StatusCreated)
}

func (h *Handler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("user_id")
	if teamName == "" {
		render.RenderJSON(w, "Team name is required", http.StatusBadRequest)
		return
	}

	teamObj, err := h.service.GetByName(teamName)
	if err != nil {
		render.RenderJSON(w, err.Message, err.Code)
		return
	}

	render.RenderJSON(w, teamObj, http.StatusOK)
}
