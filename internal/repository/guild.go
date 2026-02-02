package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"study-bot/internal/discord"
)

func (c Conn) GetGuild(guildId string) (*Guild, error) {
	var guild Guild

	query := "SELECT idx, guildName, guildId FROM Guild WHERE guildId = ?"
	if err := c.db.QueryRow(query, guildId).Scan(&guild.Idx, &guild.GuildName, &guild.GuildId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &guild, nil
}

func (c Conn) InsertGuild(guildName string, guildId string) error {
	query := "INSERT INTO Guild(guildName, guildId) VALUES (?,?)"
	if _, err := c.db.Exec(query, guildName, guildId); err != nil {
		return err
	}

	return nil
}

func (c Conn) InsertGuildChannel(guildId string, channelName string, channelId string, channelType discord.ChannelType) error {
	query := `INSERT INTO GuildChannel(channelId, guildId, channelName, channelType) 
		VALUES (:channelId, :guildId, :channelName, :channelType) 
		ON DUPLICATE KEY UPDATE channelName = :channelName, channelType = :channelType`
	if _, err := c.db.NamedExec(query, GuildChannelForm{
		GuildId:     guildId,
		ChannelId:   channelId,
		ChannelName: channelName,
		ChannelType: channelType,
	}); err != nil {
		return err
	}

	return nil
}

func (c Conn) GetGuildDMChannels(guildId string) ([]GuildChannel, error) {
	query := `SELECT * FROM GuildChannel where guildId = ? and channelType = 0`

	channels := []GuildChannel{}
	if err := c.db.Select(&channels, query, guildId); err != nil {
		return []GuildChannel{}, fmt.Errorf("Fail to get GuildChannels %w", err)
	}

	return channels, nil
}
