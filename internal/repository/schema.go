package repository

import (
	"study-bot/internal/discord"
	"time"
)

type Guild struct {
	Idx       int
	GuildName string
	GuildId   string
}
type GuildChannel struct {
	Idx         int                 `db:"idx"`
	ChannelName string              `db:"channelName"`
	GuildId     string              `db:"guildId"`
	ChannelId   string              `db:"channelId"`
	ChannelType discord.ChannelType `db:"channelType"`
	IsMain      bool                `db:"isMain"`
}

type GuildMember struct {
	GuildId    string `db:"guildId"`
	MemberId   string `db:"memberId"`
	Nickname   string `db:"nickname"`
	MemberName string `db:"memberName"`
}

type Member struct {
	Idx        int
	MemberName string
	MemberId   string
}

type VoiceState struct {
	Idx       int
	GuildId   string
	ChannelId string
	MemberId  string
	SessionId string
	EnteredAt time.Time
	LeavedAt  *time.Time
}

type GuildStatistics struct {
	MemberId   string
	MemberName string
	Time       int
}

type DailyParticipating struct {
	Idx       int       `db:"idx"`
	MemberId  string    `db:"memberId"`
	GuildId   string    `db:"guildId"`
	Date      time.Time `db:"date"`
	Duration  int       `db:"duration"`
	CreatedAt time.Time `db:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt"`
}
