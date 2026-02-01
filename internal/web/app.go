package web

import (
	"study-bot/internal/repository"

	"github.com/gin-gonic/gin"
)

type App struct {
	repo   *repository.Conn
	router *gin.Engine
}

func NewWeb(conn *repository.Conn) *App {
	app := App{
		repo:   conn,
		router: gin.Default(),
	}

	app.routes()
	return &app
}

func (a *App) Run() {
	a.router.Run()
}
