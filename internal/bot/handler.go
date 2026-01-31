package bot

import (
	"encoding/json"
	"log"
	"study-bot/internal/discord"
	"study-bot/internal/repository"
	"time"
)

type DiscordHandler interface {
	OnEvent(event discord.Event)
}

type Bot struct {
	session *discord.Session
	repo    *repository.Conn
}

func NewBot(s *discord.Session, r *repository.Conn) *Bot {
	return &Bot{session: s, repo: r}
}

func (b *Bot) OnEvent(event discord.Event) {
	if event.T == nil || event.D == nil {
		return
	}

	switch *event.T {
	case "GUILD_CREATE":
		b.init(*event.D)
	case "VOICE_STATE_UPDATE":
		b.watchVoiceState(*event.D)
	default:
		log.Printf("Unhandled event type: %s", *event.T)
	}
}

func (b *Bot) init(p json.RawMessage) {
	var payload discord.GuildCreatePayload
	if err := json.Unmarshal(p, &payload); err != nil {
		log.Printf("Fail to unmarshal payload")
		return
	}

	log.Printf("init %d %v", payload.MemberCount, payload)

	guild, err := b.repo.GetGuild(payload.Id)
	if err != nil {
		log.Printf("Err b.repo.GetGuild : %v", err)
		return
	}
	if guild == nil {
		if err := b.repo.InsertGuild(payload.Name, payload.Id); err != nil {
			log.Printf("Err b.repo.InsertGuild : %v", err)
			return
		}
	}
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

	if member == nil {
		m, err := b.repo.InsertMember(user.Username, user.Id)
		if err != nil || m == nil {
			log.Printf("Err b.repo.InsertMember %v %v", user, m)
			return
		}
		member = m
	}

	state, err := b.repo.GetCurrentVoiceStatus(member.MemberId)
	log.Printf("state %s %s", payload.SessionId, member.MemberId)
	if err != nil {
		log.Printf("Err b.repo.GetCurrentVoiceStatus %s %s", payload.SessionId, member.MemberId)
		log.Printf("%v", err)
		return
	}

	if state == nil && payload.ChannelId == nil {
		log.Printf("Unregistered session. aborted")
		return
	}

	if state == nil && payload.ChannelId != nil {
		b.repo.CreateVoiceState(repository.VoiceStateForm{
			GuildId:   payload.GuildId,
			ChannelId: *payload.ChannelId,
			MemberId:  member.MemberId,
			SessionId: payload.SessionId,
			EnteredAt: time.Now(),
		})
		log.Printf("Unregistered session. aborted")
		return
	}

	if payload.ChannelId == nil {
		n := time.Now()
		if b.repo.UpdateVoiceState(state.Idx, n) != nil {
			log.Printf("Fail to update session. sessionId: %s memberId: %s time: %v", payload.SessionId, member.MemberId, n)
			log.Printf("%v", err)
			return
		}
	}
}
