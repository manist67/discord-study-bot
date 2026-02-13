package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

func (c Conn) UpsertParticipating(guildId string, memberId string, startTime time.Time, endTime time.Time) error {
	query := `INSERT INTO DailyParticipating(guildId, memberId, date, duration) 
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE DURATION = duration + ?`

	newDuration := endTime.Sub(startTime).Seconds()
	_, err := c.db.Exec(query, guildId, memberId, startTime.Format("2006-01-02"), newDuration, newDuration)
	if err != nil {
		return err
	}

	return nil
}

func (c Conn) GetParticipating(guildId string, memberId string) ([]DailyParticipating, error) {
	query := `SELECT 
			idx, memberId, guildId, date, duration, createdAt, updatedAt 
		FROM DailyParticipating
		WHERE guildId = ? and memberId = ?
		ORDER BY date desc
	`

	var list []DailyParticipating
	err := c.db.Select(&list, query, guildId, memberId)
	if err != nil {
		return []DailyParticipating{}, err
	}

	return list, nil
}

func (c Conn) GetTotalDuration(guildId string, memberId string) (int, error) {
	query := `SELECT sum(duration) from DailyParticipating WHERE guildId = ? and memberId = ?`
	var total int
	err := c.db.Get(&total, query, guildId, memberId)
	if err != nil {
		return total, fmt.Errorf("GetTotalDuration : %w", err)
	}

	return total, nil
}

func (c Conn) GetWeekDuration(guildId string, memberId string, now time.Time) (int, error) {
	query := `SELECT 
			IFNULL(SUM(duration), 0) 
		FROM 
			DailyParticipating 
		WHERE 
			guildId = :guildId 
			AND memberId = :memberId
			AND date >= DATE_SUB(DATE(:now), INTERVAL WEEKDAY(:now) DAY)
			AND date <= DATE_ADD(DATE(:now), INTERVAL (6 - WEEKDAY(:now)) DAY)`
	nQuery, args, err := sqlx.Named(query, GetWeekDurationForm{
		GuildId:  guildId,
		MemberId: memberId,
		Now:      now,
	})
	if err != nil {
		return 0, fmt.Errorf("GetWeekDuration Named : %w", err)
	}

	nQuery = c.db.Rebind(nQuery)
	var total int
	err = c.db.Get(&total, nQuery, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("GetWeekDuration : %w", err)
	}

	return total, nil
}
