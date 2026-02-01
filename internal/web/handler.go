package web

import (
	"log"
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

	stats, err := a.repo.GetGuildStatistics(guildId, time.Now())
	if err != nil {
		c.JSON(500, gin.H{
			"error":   "INTERNAL_ERROR",
			"message": err,
		})
		return
	}

	members := make([]GuildMember, len(stats))
	for i, s := range stats {
		members[i] = NewGuildMember(s)
	}

	log.Printf("res %v", stats)

	c.JSON(200, GuildResponse{
		Guild:   NewGuild(*guild),
		Members: members,
	})
}
