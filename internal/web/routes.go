package web

func (a *App) routes() {
	api := a.router.Group("/api")
	api.GET("/", a.home)
	api.GET("/:guildId", a.guildInfo)
	api.GET("/:guildId/:memberId", a.home)
}
