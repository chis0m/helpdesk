package boot

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"helpdesk/backend/internal/config"
	"helpdesk/backend/internal/container"
	"helpdesk/backend/internal/database"
	"helpdesk/backend/internal/routes"
	"helpdesk/backend/seed"
)

type App struct {
	engine *gin.Engine
	cfg    config.Config
}

func NewApp() (*App, error) {
	cfg := config.Load()

	if err := database.RunMigrations(cfg.MySQLDSN()); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}
	if err := seed.SeedAll(db, cfg); err != nil {
		return nil, fmt.Errorf("seed data: %w", err)
	}

	c := container.New(db)

	engine := gin.Default()
	routes.Register(engine, c)

	return &App{
		engine: engine,
		cfg:    cfg,
	}, nil
}

func (a *App) Run() error {
	addr := fmt.Sprintf(":%s", a.cfg.Port)
	return a.engine.Run(addr)
}
