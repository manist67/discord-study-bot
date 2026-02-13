package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (c Conn) GetCurrentVoiceStatus(memberId string) (*VoiceState, error) {
	query := "SELECT idx, guildId, sessionId, channelId, memberId, enteredAt, leavedAt FROM VoiceState WHERE memberId = ? and leavedAt is null"

	var s VoiceState
	row := c.db.QueryRow(query, memberId)
	if err := row.Scan(&s.Idx, &s.GuildId, &s.SessionId, &s.ChannelId, &s.MemberId, &s.EnteredAt, &s.LeavedAt); err != nil {
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

func (c *Conn) GetGuildStatistics(guildId string, now time.Time) ([]GuildStatistics, error) {
	query := `
		select Member.memberId, Member.memberName, time from (
		select memberId, sum(leavedAt - enteredAt) as time
		from VoiceState as vs
		where 
			guildId = ?
			AND (DATE_FORMAT(enteredAt,  "%y-%m") = ?
				OR DATE_FORMAT(leavedAt,  "%y-%m") = ?)
			group by memberId 
		) as total_table
		left join Member on total_table.memberId  = Member.memberId 
		order by time desc
	`
	nowDate := now.Format("06-01")

	rows, err := c.db.Query(query, guildId, nowDate, nowDate)
	if err != nil {
		return []GuildStatistics{}, fmt.Errorf("GetGuildStatistics %v %v: %w", guildId, now, err)
	}

	var res []GuildStatistics
	for rows.Next() {
		var row GuildStatistics
		if err := rows.Scan(&row.MemberId, &row.MemberName, &row.Time); err != nil {
			return []GuildStatistics{}, fmt.Errorf("GetGuildStatistics scan %v %v: %w", guildId, now, err)
		}
		res = append(res, row)
	}

	return res, nil
}

func (c *Conn) GetIsOnSession(guildId string, memberId string) (bool, error) {
	query := `
		SELECT 
			count(*) > 0 
		FROM 
			VoiceState 
		WHERE 
			leavedAt is null and guildId = ? and memberId =?
	`

	var isOnSession bool
	err := c.db.Get(&isOnSession, query, guildId, memberId)
	if err != nil {
		return false, fmt.Errorf("GetIsOnSession : %v", err)
	}

	return isOnSession, nil
}
