package bot

import (
	"log"
	"study-bot/internal/discord"
)

func (b *Bot) registryGuildCommand(guildId string) {
	optionDescription := "사용자명이 뭐에요?"
	if _, err := discord.MakeGuildCommand(b.applicationId, guildId, discord.MakeGuildCommandBody{
		Name:        "info",
		Type:        1,
		Description: "사용자의 현황을 볼 수 있는 명령어에요!",
		Options: []discord.GuildCommandOption{
			{
				Name:        "user",
				Type:        9,
				Description: &optionDescription,
			},
		},
	}); err != nil {
		log.Printf("Fail to create Command : %v", err)
	}
}
