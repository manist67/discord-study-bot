package discord

import "encoding/json"

type Event struct {
	Op int              `json:"op"`
	S  *int             `json:"s"`
	T  *string          `json:"t"`
	D  *json.RawMessage `json:"d"`
}

type HandshakePayload struct {
	HeartbeatInterval int `json:"heartbeat_interval"`
}

type IdentifyPayload struct {
	Token      string `json:"token"`
	Intents    int    `json:"intents"`
	Properties struct {
		Os      string `json:"os"`
		Browser string `json:"browser"`
		Device  string `json:"device"`
	} `json:"properties"`
}

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type GuildMember struct {
	User User `json:"user"`
}

type Channel struct {
	Id   string `json:"id"`
	Type int    `json:"type"`
	Name string `json:"name"`
}

type ReadyPayload struct {
	V    int  `json:"v"`
	User User `json:"user"`
}

type GuildCreatePayload struct {
	Id          string        `json:"id"`
	Name        string        `json:"name"`
	Members     []GuildMember `json:"members"`
	Channels    []Channel     `json:"channels"`
	MemberCount int           `json:"member_count"`
}

type VoiceStatePayload struct {
	GuildId   *string     `json:"guild_id"`
	ChannelId *string     `json:"channel_id"`
	UserId    string      `json:"user_id"`
	Member    GuildMember `json:"member"`
	SessionId string      `json:"session_id"`
}

type Application struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type InteractionDataOption struct {
	Value string `json:"value"`
	Name  string `json:"name"`
	Type  int    `json:"type"`
}
type InteractionData struct {
	Id      string                  `json:"id"`
	Name    string                  `json:"name"`
	GuildId *string                 `json:"guild_id"`
	Options []InteractionDataOption `json:"options"`
}

type InteractionPayload struct {
	Id    string          `json:"id"`
	Token string          `json:"token"`
	Data  InteractionData `json:"data"`
}
