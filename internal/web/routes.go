package web

func (a *App) routes() {
	a.router.GET("/", a.home)
}
