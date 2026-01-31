package repository

import "time"

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
