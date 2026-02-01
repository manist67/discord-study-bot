package web

import (
	"time"

	"github.com/gin-gonic/gin"
)

func (a *App) home(c *gin.Context) {
	c.JSON(200, gin.H{
		"time": time.Now(),
	})
}
