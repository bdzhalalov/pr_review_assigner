package team

import (
	"encoding/json"
	app "github.com/bdzhalalov/pr-review-assigner/internal/app/team"
	"github.com/bdzhalalov/pr-review-assigner/internal/team/dto"
	"github.com/bdzhalalov/pr-review-assigner/pkg/render"
	"net/http"
)

type Handler struct {
	app     *app.TeamApp
	service *Service
}

func NewTeamHandler(app *app.TeamApp, service *Service) *Handler {
	return &Handler{
		app:     app,
		service: service,
	}
}

func (h *Handler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var input app.AddTeamRequest
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
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		render.RenderJSON(w, "Team name is required", http.StatusBadRequest)
		return
	}
	input := dto.TeamRequestDTO{
		TeamName: teamName,
	}

	teamObj, err := h.service.GetByName(input)
	if err != nil {
		render.RenderJSON(w, err.Message, err.Code)
		return
	}

	render.RenderJSON(w, teamObj, http.StatusOK)
}
