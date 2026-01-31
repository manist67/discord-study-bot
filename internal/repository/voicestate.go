package repository

import (
	"database/sql"
	"errors"
	"time"
)

func (c Conn) GetCurrentVoiceStatus(memberId string) (*VoiceState, error) {
	query := "SELECT idx, guildId, sessionId, channelId, memberId, enteredAt, leavedAt FROM VoiceState WHERE memberId = ? and leavedAt is null"

	var s VoiceState
	row := c.db.QueryRow(query, memberId)
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

func (c Conn) UpdateVoiceState(idx int, currentTime time.Time) error {
	query := "UPDATE VoiceState SET leavedAt = ? WHERE idx = ?"

	res, err := c.db.Exec(query, currentTime, idx)
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
