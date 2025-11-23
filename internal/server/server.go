package server

import (
	"errors"
	"github.com/bdzhalalov/pr-review-assigner/config"
	prApp "github.com/bdzhalalov/pr-review-assigner/internal/app/pullrequest"
	teamApp "github.com/bdzhalalov/pr-review-assigner/internal/app/team"
	userApp "github.com/bdzhalalov/pr-review-assigner/internal/app/user"
	"github.com/bdzhalalov/pr-review-assigner/internal/pullrequest"
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
	pr := pullrequest.NewPrRepository(db)

	ts := team.NewTeamService(tr, logger)
	us := user.NewUserService(ur, logger)
	ps := pullrequest.NewPrService(pr, logger)

	ta := teamApp.NewTeamApp(ts, us)
	pa := prApp.InitPRApps(ts, us, ps)
	ua := userApp.NewUserApp(ps, us)

	th := team.NewTeamHandler(ta, ts)
	teamRouter := team.TeamRouter(th)

	uh := user.NewUserHandler(us, ua)
	userRouter := user.UserRouter(uh)

	ph := pullrequest.NewPrHandler(pa, ps)
	prRouter := pullrequest.PullRequestRouter(ph)

	router := MainRouter(Routers{
		teamRouter: teamRouter,
		userRouter: userRouter,
		prRouter:   prRouter,
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
