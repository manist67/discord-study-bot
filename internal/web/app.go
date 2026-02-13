package web

import (
	"context"
	"embed"
	"study-bot/internal/repository"

	"github.com/gin-gonic/gin"
)

//go:embed dist/*
var contents embed.FS

type App struct {
	repo   *repository.Conn
	router *gin.Engine
}

func NewWeb(conn *repository.Conn) *App {
	app := App{
		repo:   conn,
		router: gin.Default(),
	}

	app.router.Handlers = append(app.router.Handlers, CORSMiddleware())

	app.routes()
	app.statics()

	return &app
}

func (a *App) Run(ctx context.Context) {
	a.router.Run()
}
