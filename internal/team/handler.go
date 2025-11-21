package team

import (
	"encoding/json"
	"github.com/bdzhalalov/pr-review-assigner/internal/app/add_team"
	"github.com/bdzhalalov/pr-review-assigner/internal/team/dto"
	"github.com/bdzhalalov/pr-review-assigner/pkg/render"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	app     *add_team.TeamApp
	service *Service
	logger  *logrus.Logger
}

func NewTeamHandler(app *add_team.TeamApp, service *Service, logger *logrus.Logger) *Handler {
	return &Handler{
		app:     app,
		service: service,
		logger:  logger,
	}
}

func (h *Handler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var input add_team.AddTeamRequest
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
	input := dto.GetTeamByNameDTO{
		TeamName: teamName,
	}

	teamObj, err := h.service.GetByName(input)
	if err != nil {
		render.RenderJSON(w, err.Message, err.Code)
		return
	}

	render.RenderJSON(w, teamObj, http.StatusOK)
}
