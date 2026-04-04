package boot

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/config"
	"helpdesk/backend/internal/container"
	"helpdesk/backend/internal/database"
	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/routes"
	"helpdesk/backend/seed"
)

type App struct {
	engine *gin.Engine
	cfg    config.Config
}

func NewApp() (*App, error) {
	cfg := config.Load()
	if err := logger.Init(cfg); err != nil {
		return nil, fmt.Errorf("initialize logger: %w", err)
	}
	log := logger.L()
	log.Info().Msg("bootstrapping backend application")

	if err := database.RunMigrations(cfg.MySQLDSN()); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}
	log.Info().Msg("migrations completed")

	db, err := database.Connect(cfg)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}
	if err := seed.SeedAll(db, cfg); err != nil {
		return nil, fmt.Errorf("seed data: %w", err)
	}
	log.Info().Msg("seed process completed")

	tokenMaker, err := auth.NewPasetoMaker(cfg.PasetoSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("initialize paseto maker: %w", err)
	}
	log.Info().Msg("paseto maker initialized")

	c := container.New(db, tokenMaker)

	engine := gin.Default()
	routes.Register(engine, c)

	return &App{
		engine: engine,
		cfg:    cfg,
	}, nil
}

func (a *App) Run() error {
	defer func() {
		_ = logger.Sync()
	}()

	addr := fmt.Sprintf(":%s", a.cfg.Port)
	logger.L().Info().Str("address", addr).Msg("starting HTTP server")
	return a.engine.Run(addr)
}
