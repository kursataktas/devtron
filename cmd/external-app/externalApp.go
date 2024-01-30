package main

import (
	"fmt"
	"net/http"
	"os"

	authMiddleware "github.com/devtron-labs/authenticator/middleware"
	"github.com/devtron-labs/common-lib/middlewares"
	"github.com/devtron-labs/devtron/client/telemetry"
	"github.com/devtron-labs/devtron/internal/middleware"
	"github.com/devtron-labs/devtron/pkg/auth/user"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

type App struct {
	db             *pg.DB
	sessionManager *authMiddleware.SessionManager
	MuxRouter      *MuxRouter
	Logger         *zap.SugaredLogger
	server         *http.Server
	telemetry      telemetry.TelemetryEventClient
	posthogClient  *telemetry.PosthogClient
	userService    user.UserService
}

func NewApp(db *pg.DB,
	sessionManager *authMiddleware.SessionManager,
	MuxRouter *MuxRouter,
	telemetry telemetry.TelemetryEventClient,
	posthogClient *telemetry.PosthogClient,
	Logger *zap.SugaredLogger,
	userService user.UserService) *App {
	return &App{
		db:             db,
		sessionManager: sessionManager,
		MuxRouter:      MuxRouter,
		Logger:         Logger,
		telemetry:      telemetry,
		posthogClient:  posthogClient,
		userService:    userService,
	}
}
func (app *App) Start() {
	fmt.Println("starting ea module")

	port := 8080 //TODO: extract from environment variable
	app.Logger.Debugw("starting server")
	app.Logger.Infow("starting server on ", "port", port)
	app.MuxRouter.Init()
	//authEnforcer := casbin2.Create()

	_, err := app.telemetry.SendTelemetryInstallEventEA()

	if err != nil {
		app.Logger.Warnw("telemetry installation success event failed", "err", err)
	}
	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: authMiddleware.Authorizer(app.sessionManager, user.WhitelistChecker, app.userService.UserStatusCheckInDb)(app.MuxRouter.Router)}
	app.MuxRouter.Router.Use(middleware.PrometheusMiddleware)
	app.MuxRouter.Router.Use(middlewares.Recovery)
	app.server = server

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Errorw("error in startup", "err", err)
		os.Exit(2)
	}
}

func (app *App) Stop() {
	app.Logger.Info("orchestrator shutdown initiating")
	posthogCl := app.posthogClient.Client
	if posthogCl != nil {
		app.Logger.Info("flushing messages of posthog")
		posthogCl.Close()
	}
}
