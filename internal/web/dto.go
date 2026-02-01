package web

import "study-bot/internal/repository"

type GuildMember struct {
	MemberId   string `json:"memberId"`
	MemberName string `json:"memberName"`
	Time       int    `json:"time"`
}

func NewGuildMember(s repository.GuildStatistics) GuildMember {
	return GuildMember{
		MemberId:   s.MemberId,
		MemberName: s.MemberName,
		Time:       s.Time,
	}
}

type Guild struct {
	Idx       int    `json:"idx"`
	GuildName string `json:"guildName"`
	GuildId   string `json:"guildId"`
}

func NewGuild(g repository.Guild) Guild {
	return Guild(g)
}

type GuildResponse struct {
	Guild   Guild         `json:"guild"`
	Members []GuildMember `json:"members"`
}
