package bot

import (
	"errors"
	"fmt"
	"log"
	"study-bot/internal/discord"
	"study-bot/internal/repository"
	"time"
)

func (b *Bot) registryGuildCommand(guildId string) {
	{
		// 정보 조회 커맨드 생성
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

	{
		// 정보 세팅 커맨드 생성
		if _, err := discord.MakeGuildCommand(b.applicationId, guildId, discord.MakeGuildCommandBody{
			Name:        "set_channel",
			Type:        discord.ChatInput,
			Description: "입장 메세지가 표시되고 싶은 채널에서 호출해주세요.",
			Options:     []discord.GuildCommandOption{},
		}); err != nil {
			log.Printf("Fail to create Command : %v", err)
		}
	}
}

func (b *Bot) enterVoiceChannel(member *repository.Member, payload discord.VoiceStatePayload, displayName string) error {
	if member == nil {
		return fmt.Errorf("no member")
	}

	now := time.Now()
	if err := b.repo.CreateVoiceState(repository.VoiceStateForm{
		GuildId:   payload.GuildId,
		ChannelId: *payload.ChannelId,
		MemberId:  member.MemberId,
		SessionId: payload.SessionId,
		EnteredAt: now,
	}); err != nil {
		return fmt.Errorf("Fail to create voice state %w", err)
	}

	channels, err := b.repo.GetGuildDMChannels(*payload.GuildId)
	if err != nil {
		return fmt.Errorf("Fail b.repo.GetGuildDMChannels %w", err)
	}

	if len(channels) <= 0 {
		return errors.New("No channel")
	}

	mainChannel := channels[0]
	for _, channel := range channels {
		if channel.IsMain {
			mainChannel = channel
			break
		}
	}

	if err := discord.SendMessage(mainChannel.ChannelId, discord.MessageForm{
		Content: fmt.Sprintf("%s 입장 시간 : %s",
			displayName,
			now.Local().Format("2006-01-02 15:04:05")),
	}); err != nil {
		return fmt.Errorf("Fail to send message %w", err)
	}

	return nil
}

func (b *Bot) leaveVoiceChannel(state *repository.VoiceState, member *repository.Member, displayName string) error {
	if state == nil {
		return errors.New("no state")
	}
	if member == nil {
		return errors.New("no member")
	}

	leaveDate := time.Now()
	if err := b.repo.UpdateVoiceState(state.Idx, leaveDate); err != nil {
		return fmt.Errorf("Fail to update voice state %w", err)
	}

	startAt := state.EnteredAt
	for {
		nextDay := time.Date(startAt.Year(), startAt.Month(), startAt.Day()+1, 0, 0, 0, 0, startAt.Location())
		if nextDay.After(leaveDate) {
			if err := b.repo.UpsertParticipating(state.GuildId, member.MemberId, startAt, leaveDate); err != nil {
				return err
			}
			break
		}
		if err := b.repo.UpsertParticipating(state.GuildId, member.MemberId, startAt, nextDay); err != nil {
			return err
		}
		startAt = nextDay
	}

	channels, err := b.repo.GetGuildDMChannels(state.GuildId)
	if err != nil {
		return fmt.Errorf("Fail b.repo.GetGuildDMChannels %w", err)
	}

	if len(channels) <= 0 {
		return errors.New("No channel")
	}

	mainChannel := channels[0]
	for _, channel := range channels {
		if channel.IsMain {
			mainChannel = channel
			break
		}
	}

	if err := discord.SendMessage(mainChannel.ChannelId, discord.MessageForm{
		Content: fmt.Sprintf("%s 퇴장 시간 : %s",
			displayName,
			leaveDate.Local().Format("2006-01-02 15:04:05")),
	}); err != nil {
		return fmt.Errorf("Fail to send message %w", err)
	}

	return nil
}
