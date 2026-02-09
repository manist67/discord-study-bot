package bot

import (
	"encoding/json"
	"log"
	"study-bot/internal/discord"
)

type DiscordHandler interface {
	OnEvent(event discord.Event)
}

func (b *Bot) OnEvent(event discord.Event) {
	if event.T == nil || event.D == nil {
		return
	}

	switch *event.T {
	case "READY":
		b.ready(*event.D)
	case "GUILD_CREATE":
		b.createGuild(*event.D)
	case "VOICE_STATE_UPDATE":
		b.watchVoiceState(*event.D)
	case "INTERACTION_CREATE":
		b.handleInteraction(*event.D)
	case "CHANNEL_CREATE":
		b.handleCreateChannel(*event.D)
	case "CHANNEL_DELETE":
		b.handleCreateChannel(*event.D)
	default:
		log.Printf("Unhandled event type: %s", *event.T)
	}
}

func (b *Bot) ready(p json.RawMessage) {
	var payload discord.ReadyPayload
	if err := json.Unmarshal(p, &payload); err != nil {
		log.Printf("Fail to unmarshal payload %s", string(p))
		return
	}

	b.applicationId = payload.User.Id
}

func (b *Bot) createGuild(p json.RawMessage) {
	var payload discord.GuildCreatePayload
	if err := json.Unmarshal(p, &payload); err != nil {
		log.Printf("Fail to unmarshal payload")
		return
	}

	log.Printf("init %s", payload.Name)

	guild, err := b.repo.GetGuild(payload.Id)
	if err != nil {
		log.Printf("Err b.repo.GetGuild : %v", err)
		return
	}
	if guild == nil {
		g, err := b.repo.InsertGuild(payload.Name, payload.Id)
		if err != nil {
			log.Printf("Err b.repo.InsertGuild : %v", err)
			return
		}

		guild = g
	}

	for _, c := range payload.Channels {
		if err := b.repo.InsertGuildChannel(guild.GuildId, c.Name, c.Id, c.Type); err != nil {
			log.Printf("Err b.repo.InsertGuildChannel : %v", err)
			return
		}
	}

	b.registryGuildCommand(guild.GuildId)
}

func (b *Bot) watchVoiceState(p json.RawMessage) {
	var payload discord.VoiceStatePayload
	if err := json.Unmarshal(p, &payload); err != nil {
		log.Printf("Fail to unmarshal payload %s", string(p))
		return
	}

	log.Printf("watchVoiceState %v", payload.Member.User)
	log.Printf("%v", string(p))
	user := payload.Member.User
	member, err := b.repo.GetMemberById(user.Id)
	if err != nil {
		log.Printf("Err b.repo.GetMemberById %s", user.Id)
		return
	}

	// 맴버가 없을 경우 맴버 삽입
	if member == nil {
		m, err := b.repo.InsertMember(user.Username, user.Id)
		if err != nil || m == nil {
			log.Printf("Err b.repo.InsertMember %v %v", user, m)
			return
		}
		member = m
	}

	var displayName string
	if payload.Member.Nick != nil {
		displayName = *payload.Member.Nick
	} else if payload.Member.User.DisplayName != nil {
		displayName = *payload.Member.User.DisplayName
	} else {
		displayName = payload.Member.User.Username
	}

	if payload.GuildId != nil {
		if err := b.repo.InsertGuildMember(*payload.GuildId, member.MemberId, displayName); err != nil {
			log.Printf("Err b.repo.InsertGuildMember: %v", err)
			return
		}
	}

	state, err := b.repo.GetCurrentVoiceStatus(member.MemberId)
	log.Printf("state %s %s", payload.SessionId, member.MemberId)
	if err != nil {
		log.Printf("Err b.repo.GetCurrentVoiceStatus %s %s", payload.SessionId, member.MemberId)
		log.Printf("%v", err)
		return
	}

	// 서버가 켜지기 전 체널에 이미 들어간 경우
	// 일반적으로는 오류임
	if state == nil && payload.ChannelId == nil {
		log.Printf("Unregistered session. aborted")
		return
	}

	if state == nil && payload.ChannelId != nil { // 서버에 입장 경우
		if err := b.enterVoiceChannel(member, payload, displayName); err != nil {
			log.Printf("Fail to create session. sessionId: %s memberId: %s", payload.SessionId, member.MemberId)
			log.Printf("%v", err)
		}
	} else if state != nil && payload.ChannelId == nil { // 서버에 퇴장한 경우
		if err := b.leaveVoiceChannel(state, member, displayName); err != nil {
			log.Printf("Fail to update session. sessionId: %s memberId: %s", payload.SessionId, member.MemberId)
			log.Printf("%v", err)
		}
	} else { // 오류 상황
		log.Printf("Unregistered session. aborted")
		return

	}
}

func (b *Bot) handleCreateChannel(p json.RawMessage) {

}

func (b *Bot) handleDeleteChannel(p json.RawMessage) {

}
