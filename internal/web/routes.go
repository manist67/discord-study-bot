package web

import (
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (a *App) routes() {
	api := a.router.Group("/api")
	api.GET("/", a.home)
	api.GET("/:guildId", a.guildInfo)
	api.GET("/:guildId/:memberId", a.memberInfo)
}

func (a *App) statics() {
	fsSub, err := fs.Sub(contents, "dist")
	if err != nil {
		log.Panicf("Fail to serve react : %v", err)
	}

	fileServer := http.FileServer(http.FS(fsSub))
	a.router.NoRoute(func(c *gin.Context) {
		path := strings.TrimPrefix(c.Request.URL.Path, "/")
		_, err := fsSub.Open(path)
		if err != nil {
			indexFile, _ := fs.ReadFile(fsSub, "index.html")
			c.Data(http.StatusOK, "text/html", indexFile)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}
