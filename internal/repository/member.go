package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func (c Conn) GetMemberById(memberId string) (*Member, error) {
	var m Member
	if err := c.db.QueryRow("SELECT idx, memberName, memberId FROM Member WHERE memberId = ?", memberId).Scan(&m.Idx, &m.MemberName, &m.MemberId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &m, nil
}

func (c Conn) GetMembersByIds(memberIds []string) ([]Member, error) {
	if len(memberIds) == 0 {
		return []Member{}, nil
	}

	placeholders := make([]string, len(memberIds))
	var args []any

	for idx, v := range memberIds {
		placeholders[idx] = "?"
		args = append(args, v)
	}

	query := fmt.Sprintf("SELECT idx, memberName, memberId FROM Member WHERE memberId in (%s)",
		strings.Join(placeholders, ","))

	rows, err := c.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var members []Member
	for rows.Next() {
		var m Member
		if err := rows.Scan(&m.Idx, &m.MemberName, &m.MemberId); err != nil {
			return nil, err
		}

		members = append(members, m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

func (c Conn) GetGuildMember(guildId string, memberId string) (*GuildMember, error) {
	query := "SELECT guildId, memberId, nickname FROM GuildMember where guildId = ? and memberId = ?"
	var guildMember GuildMember
	err := c.db.Get(&guildMember, query, guildId, memberId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetGuildMember %w", err)
	}

	return &guildMember, nil
}

func (c Conn) InsertGuildMember(guildId string, memberId string, nickname string) error {
	query := `INSERT INTO GuildMember(guildId, memberId, nickname) 
		VALUES (:guildId, :memberId, :nickname)
		ON DUPLICATE KEY UPDATE nickname = :nickname`

	if _, err := c.db.NamedExec(query, GuildMemberForm{
		Nickname: nickname,
		GuildId:  guildId,
		MemberId: memberId,
	}); err != nil {
		return err
	}

	return nil
}

func (c Conn) InsertMember(memberName string, memberId string) (*Member, error) {
	query := "INSERT INTO Member(memberName, memberId) VALUES (?, ?)"
	res, err := c.db.Exec(query, memberName, memberId)
	if err != nil {
		return nil, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Member{
		Idx:        int(lastId),
		MemberName: memberName,
		MemberId:   memberId,
	}, nil
}

func (c Conn) InsertMembers(members []MemberForm) error {
	if len(members) == 0 {
		return nil
	}
	placeholders := make([]string, len(members))
	var args []any
	for idx, m := range members {
		placeholders[idx] = "(?, ?)"
		args = append(args, m.MemberId, m.MemberName)
	}

	query := fmt.Sprintf("INSERT INTO Member(memberId, memberName) VALUES %s", strings.Join(placeholders, ","))

	if _, err := c.db.Exec(query, args...); err != nil {
		return err
	}

	return nil
}
