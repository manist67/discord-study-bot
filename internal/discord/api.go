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

func InteractionCallback(interactionId string, interactionToken string, body InteractionCallbackForm) error {
	url := fmt.Sprintf("https://discord.com/api/v10/interactions/%s/%s/callback", interactionId, interactionToken)
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	reqs, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	reqs.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(reqs)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("(%d) %v %v", resp.StatusCode, string(resBody), body)
	}
	return nil
}

func SendMessage(channelId string, form MessageForm) error {
	url := fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", channelId)
	payload, err := json.Marshal(form)
	if err != nil {
		return err
	}

	reqs, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	token := fmt.Sprintf("Bot %s", os.Getenv("DISCORD_BOT_TOKEN"))
	reqs.Header.Add("Authorization", token)
	reqs.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(reqs)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Printf("Message Created : %s %v", url, form)
	resBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Fail to response interaction %v", string(resBody))
	}
	return nil

}
