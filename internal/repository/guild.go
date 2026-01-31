package repository

import (
	"database/sql"
	"errors"
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
