package repository

import (
	"database/sql"
	"errors"
	"time"
)

func (c Conn) GetCurrentVoiceStatus(sessionId string, memberId string) (*VoiceState, error) {
	query := "SELECT idx, guildId, sessionId, channelId, memberId, enteredAt, leavedAt FROM VoiceState WHERE sessionId = ? and memberId = ?"

	var s VoiceState
	row := c.db.QueryRow(query, sessionId, memberId)
	if err := row.Scan(&s.Idx, &s.GuildId, &s.SessionId, &s.ChannelId, &s.MemberId, &s.enteredAt, &s.leavedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &s, nil
}

func (c Conn) CreateVoiceState(form VoiceStateForm) error {
	query := "INSERT INTO VoiceState(guildId, channelId, memberId, sessionId, enteredAt) values (?, ?,?,?,?)"

	res, err := c.db.Exec(query, form.GuildId, form.ChannelId, form.MemberId, form.SessionId, form.EnteredAt)
	if err != nil {
		return err
	}
	rowAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return errors.New("Fail to insert")
	}

	return nil
}

func (c Conn) UpdateVoiceState(sessionId string, memberId string, currentTime time.Time) error {
	query := "UPDATE VoiceState SET leavedAt = ? WHERE sessionId = ? and memberId = ?"

	res, err := c.db.Exec(query, currentTime, sessionId, memberId)
	if err != nil {
		return err
	}
	rowAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return errors.New("Fail to insert")
	}

	return nil
}
