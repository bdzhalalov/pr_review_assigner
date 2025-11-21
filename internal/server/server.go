package server

import (
	"errors"
	"github.com/bdzhalalov/pr-review-assigner/config"
	"github.com/bdzhalalov/pr-review-assigner/internal/app/add_team"
	"github.com/bdzhalalov/pr-review-assigner/internal/team"
	"github.com/bdzhalalov/pr-review-assigner/internal/user"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

type APIServer struct {
	config *config.Config
	logger *logrus.Logger
	server *http.Server
}

func Init(config *config.Config, logger *logrus.Logger, db *gorm.DB) *APIServer {
	tr := team.NewTeamRepository(db)
	ur := user.NewUserRepository(db)

	ts := team.NewTeamService(tr, logger)
	us := user.NewUserService(ur, logger)

	teamApp := add_team.NewTeamApp(ts, us)

	th := team.NewTeamHandler(teamApp, ts, logger)
	teamRouter := team.TeamRouter(th)

	router := MainRouter(Routers{
		teamRouter: teamRouter,
	})

	return &APIServer{
		config: config,
		logger: logger,
		server: &http.Server{
			Addr:    config.Addr,
			Handler: router,
		},
	}
}

func (s *APIServer) Run() error {

	s.logger.Info("Running API server on port" + s.config.Addr)

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.WithError(err).Fatal("Failed to start API server")
	}

	return nil
}
