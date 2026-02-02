package repository

import (
	"time"
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
