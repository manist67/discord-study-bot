package repository

import "time"

type Guild struct {
	Idx       int
	GuildName string
	GuildId   string
}

type Member struct {
	Idx        int
	MemberName string
	MemberId   string
}

type VoiceState struct {
	Idx       int
	GuildId   *string
	ChannelId *string
	MemberId  string
	SessionId string
	enteredAt time.Time
	leavedAt  *time.Time
}

type GuildStatistics struct {
	MemberId   string
	MemberName string
	Time       int
}
