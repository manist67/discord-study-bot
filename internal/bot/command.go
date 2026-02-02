package bot

import (
	"errors"
	"fmt"
	"log"
	"os"
	"study-bot/internal/discord"
	"study-bot/internal/repository"
	"time"
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

func (b *Bot) enterVoiceChannel(member *repository.Member, payload discord.VoiceStatePayload) error {
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

	if err := discord.SendMessage(channels[0].ChannelId, discord.MessageForm{
		Content: fmt.Sprintf("%s 님 안녕하세요! 입장 시간 : %s", member.MemberName, now.Local().Format("2006-01-02 15:04:05")),
	}); err != nil {
		return fmt.Errorf("Fail to send message %w", err)
	}

	return nil
}

func (b *Bot) leaveVoiceChannel(state *repository.VoiceState, member *repository.Member) error {
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

	if err := discord.SendMessage(channels[0].ChannelId, discord.MessageForm{
		Content: fmt.Sprintf("%s 님 조심히 들어가세요! 활동 시간 : %s ~ %s",
			member.MemberName,
			state.EnteredAt.Local().Format("2006-01-02 15:04:05"),
			leaveDate.Local().Format("2006-01-02 15:04:05")),
	}); err != nil {
		return fmt.Errorf("Fail to send message %w", err)
	}

	return nil
}
