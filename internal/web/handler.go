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

func (a *App) memberInfo(c *gin.Context) {
	guildId := c.Param("guildId")
	memberId := c.Param("memberId")

	member, err := a.repo.GetGuildMember(guildId, memberId)
	if err != nil {
		log.Printf("Error %v", err)
		c.JSON(500, gin.H{
			"error":   "SERVER_ERROR",
			"message": err,
		})
		return
	}

	if member == nil {
		c.JSON(404, gin.H{
			"error":   "INVALID_MEMBER",
			"message": err,
		})
		return
	}
	participatingList, err := a.repo.GetParticipating(guildId, memberId)

	if err != nil {
		log.Printf("Error %v", err)
		c.JSON(500, gin.H{
			"error":   "SERVER_ERROR",
			"message": err,
		})
		return
	}

	list := []Participating{}
	for _, particiapting := range participatingList {
		list = append(list, Participating{
			Date:     particiapting.Date,
			Duration: particiapting.Duration,
		})
	}

	totalDuration, err := a.repo.GetTotalDuration(guildId, memberId)
	if err != nil {
		log.Printf("Error %v", err)
		c.JSON(500, gin.H{
			"error":   "SERVER_ERROR",
			"message": err,
		})
		return
	}

	weekDuration, err := a.repo.GetWeekDuration(guildId, memberId, time.Now())
	if err != nil {
		log.Printf("Error %v", err)
		c.JSON(500, gin.H{
			"error":   "SERVER_ERROR",
			"message": err,
		})
		return
	}

	c.JSON(200, MemberActivity{
		Member: Member{
			Nickname: member.Nickname,
		},
		Total:             totalDuration,
		WeekTotal:         weekDuration,
		ParticipatingList: list,
	})
}
