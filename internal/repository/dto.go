package repository

import (
	"study-bot/internal/discord"
	"time"
)

type MemberForm struct {
	MemberName string
	MemberId   string
}

type VoiceStateForm struct {
	GuildId   *string
	ChannelId string
	MemberId  string
	SessionId string
	EnteredAt time.Time
}

type GuildChannelForm struct {
	GuildId     string              `db:"guildId"`
	ChannelId   string              `db:"channelId"`
	ChannelName string              `db:"channelName"`
	ChannelType discord.ChannelType `db:"channelType"`
}
