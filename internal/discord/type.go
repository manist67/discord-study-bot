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

type ResumePayload struct {
	Token     string `json:"token"`
	SessionId string `json:"session_id"`
	Seq       int    `json:"seq"`
}

type User struct {
	Id          string  `json:"id"`
	Username    string  `json:"username"`
	DisplayName *string `json:"display_name"`
}

type GuildMember struct {
	User User    `json:"user"`
	Nick *string `json:"nick"`
}

type Guild struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ChannelType int

const (
	GUILD_TEXT          = 0
	DM                  = 1
	GUILD_VOICE         = 2
	GROUP_DM            = 3
	GUILD_CATEGORY      = 4
	GUILD_ANNOUNCEMENT  = 5
	ANNOUNCEMENT_THREAD = 10
	PUBLIC_THREAD       = 11
	PRIVATE_THREAD      = 12
	GUILD_STAGE_VOICE   = 13
	GUILD_DIRECTORY     = 14
	GUILD_FORUM         = 15
	GUILD_MEDIA         = 16
)

type Channel struct {
	Id      string      `json:"id"`
	Type    ChannelType `json:"type"`
	Name    string      `json:"name"`
	GuildId *string     `json:"guild_id"`
}

type ReadyPayload struct {
	V                int    `json:"v"`
	User             User   `json:"user"`
	SessionId        string `json:"session_id"`
	ResumeGatewayUrl string `json:"resume_gateway_url"`
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
	Id      string          `json:"id"`
	Token   string          `json:"token"`
	Data    InteractionData `json:"data"`
	Guild   *Guild          `json:"guild"`
	Channel *Channel        `json:"channel"`
}

type GuildCommandOptionType int

const (
	ChannelOption GuildCommandOptionType = 7
	Mentionable   GuildCommandOptionType = 9
)

type GuildCommandOption struct {
	Type                     GuildCommandOptionType `json:"type"`
	Name                     string                 `json:"name"`
	Required                 bool                   `json:"required"`
	Description              *string                `json:"description"`
	DescriptionLocalizations *string                `json:"description_localizations"`
}

type CommandType int

const (
	ChatInput CommandType = 1
)

type MakeGuildCommandBody struct {
	Name        string               `json:"name"`
	Type        CommandType          `json:"type"`
	Description string               `json:"description"`
	Options     []GuildCommandOption `json:"options"`
}

type GuildCommand struct {
	Id                       string               `json:"id"`
	Name                     string               `json:"name"`
	NameLocalizations        *string              `json:"name_localizations"`
	Type                     int                  `json:"type"`
	Description              *string              `json:"description"`
	DescriptionLocalizations *string              `json:"description_localizations"`
	ApplicationId            string               `json:"application_id"`
	Version                  string               `json:"version"`
	DefaultMemberPermissions *string              `json:"default_member_permissions"`
	Options                  []GuildCommandOption `json:"options"`
}

type InteractionCallbackData struct {
	Content string `json:"content"`
}

type InteractionCallbackType int

const (
	ChannelMessageWithSource InteractionCallbackType = 4
)

type InteractionCallbackForm struct {
	Type InteractionCallbackType `json:"type"`
	Data InteractionCallbackData `json:"data"`
}

type MessageForm struct {
	Content string `json:"content"`
}
