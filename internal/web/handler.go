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

func (a *App) guildInfo(c *gin.Context) {
	guildId := c.Param("guildId")
	guild, err := a.repo.GetGuild(guildId)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "INTERNAL_ERROR",
		})
		return
	}

	c.JSON(200, gin.H{
		"name": guild.GuildName,
	})
}
