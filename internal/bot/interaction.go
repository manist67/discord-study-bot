package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"study-bot/internal/discord"
)

func (b *Bot) handleInteraction(p json.RawMessage) {
	var payload discord.InteractionPayload
	if err := json.Unmarshal(p, &payload); err != nil {
		log.Printf("Fail to unmarshal payload %v %s", err, string(p))
		return
	}
	switch payload.Data.Name {
	case "info":
		b.handleInfoInteraction(payload)
	case "set_channel":
		b.handleSetChannel(payload)
	}
}

func (b *Bot) handleInfoInteraction(payload discord.InteractionPayload) {
	data := payload.Data
	if data.GuildId == nil {
		log.Printf("Info command used outside of a guild")
		return
	}

	if len(data.Options) > 0 {
		b.responseMemberInfoLink(payload.Id, payload.Token, *data.GuildId, data.Options[0].Value)
	} else {
		b.responseGuildInfoLink(payload.Id, payload.Token, *data.GuildId)
	}
}

func (b *Bot) handleSetChannel(payload discord.InteractionPayload) {
	interactionId := payload.Id
	token := payload.Token

	if payload.Guild == nil || payload.Channel == nil {
		if err := discord.InteractionCallback(interactionId, token, discord.InteractionCallbackForm{
			Type: discord.ChannelMessageWithSource,
			Data: discord.InteractionCallbackData{
				Content: fmt.Sprintf("채널 설정 실패 : %s", "guild, channel 없음"),
			},
		}); err != nil {
			log.Printf("ERROR : %v", err)
		}
		return
	}

	if err := b.repo.UpdateGuildChannel(payload.Guild.Id, payload.Channel.Id); err != nil {
		if err := discord.InteractionCallback(interactionId, token, discord.InteractionCallbackForm{
			Type: discord.ChannelMessageWithSource,
			Data: discord.InteractionCallbackData{
				Content: fmt.Sprintf("채널 설정 실패 : %s", "DB 업데이트 실패"),
			},
		}); err != nil {
			log.Printf("ERROR : %v", err)
		}
		return
	}

	b.responseSetting(interactionId, token, payload.Channel.Name)
}

func (b *Bot) responseSetting(interactionId string, token string, channelName string) {
	if err := discord.InteractionCallback(interactionId, token, discord.InteractionCallbackForm{
		Type: discord.ChannelMessageWithSource,
		Data: discord.InteractionCallbackData{
			Content: fmt.Sprintf("설정 채널 : %s", channelName),
		},
	}); err != nil {
		log.Printf("ERROR : %v", err)
	}
}

func (b *Bot) responseGuildInfoLink(interactionId string, token string, guildId string) {
	guild, err := b.repo.GetGuild(guildId)
	if err != nil {
		log.Printf("Fail to get guild : %v", err)
		if err := discord.InteractionCallback(interactionId, token, discord.InteractionCallbackForm{
			Type: discord.ChannelMessageWithSource,
			Data: discord.InteractionCallbackData{
				Content: "오류가 발생했어요! 잠시 후 다시 시도해주세요.",
			},
		}); err != nil {
			log.Printf("Fail to create Command : %s %s %v", interactionId, guildId, err)
			return
		}
	}

	if err := discord.InteractionCallback(interactionId, token, discord.InteractionCallbackForm{
		Type: discord.ChannelMessageWithSource,
		Data: discord.InteractionCallbackData{
			Content: fmt.Sprintf("%s 의 활동 내역이에요!\n%s/%s", guild.GuildName, os.Getenv("HOST_URL"), guildId),
		},
	}); err != nil {
		log.Printf("Fail to create Command : %s %s %v", interactionId, guildId, err)
	}
}

func (b *Bot) responseMemberInfoLink(interactionId string, token string, guildId string, memberId string) {
	member, err := b.repo.GetMemberById(memberId)
	if err != nil {
		log.Printf("Fail to get member : %v", err)
		if err := discord.InteractionCallback(interactionId, token, discord.InteractionCallbackForm{
			Type: discord.ChannelMessageWithSource,
			Data: discord.InteractionCallbackData{
				Content: "오류가 발생했어요! 잠시 후 다시 시도해주세요.",
			},
		}); err != nil {
			log.Printf("Fail to create Command : %s %s %s %v", interactionId, guildId, memberId, err)
			return
		}
	}

	if err := discord.InteractionCallback(interactionId, token, discord.InteractionCallbackForm{
		Type: discord.ChannelMessageWithSource,
		Data: discord.InteractionCallbackData{
			Content: fmt.Sprintf("%s 님의 활동 내역이에요!\n%s/%s/%s", member.MemberName, os.Getenv("HOST_URL"), guildId, memberId),
		},
	}); err != nil {
		log.Printf("Fail to create Command : %s %s %s %v", interactionId, guildId, memberId, err)
	}
}
