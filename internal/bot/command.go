package bot

import (
	"fmt"
	"log"
	"os"
	"study-bot/internal/discord"
)

func (b *Bot) registryGuildCommand(guildId string) {
	optionDescription := "사용자명이 뭐에요?"
	if _, err := discord.MakeGuildCommand(b.applicationId, guildId, discord.MakeGuildCommandBody{
		Name:        "info",
		Type:        discord.ChatInput,
		Description: "사용자의 현황을 볼 수 있는 명령어에요!",
		Options: []discord.GuildCommandOption{
			{
				Name:        "user",
				Type:        discord.Mentionable,
				Description: &optionDescription,
			},
		},
	}); err != nil {
		log.Printf("Fail to create Command : %v", err)
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
