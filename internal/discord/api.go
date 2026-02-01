package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type GuildCommandOption struct {
	Type                     int     `json:"type"`
	Name                     string  `json:"name"`
	Description              *string `json:"description"`
	DescriptionLocalizations *string `json:"description_localizations"`
}

type MakeGuildCommandBody struct {
	Name        string               `json:"name"`
	Type        int                  `json:"type"`
	Description string               `json:"description"`
	Options     []GuildCommandOption `json:"options"`
}

type GuildCommand struct {
	Id                       string               `json:"id"`
	Name                     string               `json:"name"`
	NameLocalizations        *string              `json:"name_localizations"`
	Type                     int                  `json:"type"`
	Description              *string              `json:"description"`
	DescriptionLocalizations *string              `json:"description_localizations"`
	ApplicationId            string               `json:"application_id"`
	Version                  string               `json:"version"`
	DefaultMemberPermissions *string              `json:"default_member_permissions"`
	Options                  []GuildCommandOption `json:"options"`
}

func MakeGuildCommand(applicationId string, guildId string, body MakeGuildCommandBody) (GuildCommand, error) {
	url := fmt.Sprintf("https://discord.com/api/v10/applications/%s/guilds/%s/commands", applicationId, guildId)

	payload, err := json.Marshal(body)
	if err != nil {
		return GuildCommand{}, err
	}

	reqs, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return GuildCommand{}, err
	}

	token := fmt.Sprintf("Bot %s", os.Getenv("DISCORD_BOT_TOKEN"))
	reqs.Header.Add("Authorization", token)
	reqs.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(reqs)
	if err != nil {
		return GuildCommand{}, err
	}
	defer resp.Body.Close()

	log.Printf("Command Created : %s %v", url, body)
	resBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return GuildCommand{}, err
	}

	var command GuildCommand
	if err := json.Unmarshal(resBody, &command); err != nil {
		return GuildCommand{}, err
	}

	return command, nil
}
