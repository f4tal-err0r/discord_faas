package discord

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OAuth2Response struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type GuildMember struct {
	User        OAuth2Response `json:"user"`
	Nickname    string         `json:"nick"`
	Roles       []string       `json:"roles"`
	Permissions int64          `json:"permissions"`
}

func IdentGuildMember(token, gid string) (*GuildMember, error) {
	endpoint := fmt.Sprintf("https://discord.com/api/users/@me/guilds/%s/member", gid)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch guild member: %s", resp.Status)
	}

	var member GuildMember
	if err := json.NewDecoder(resp.Body).Decode(&member); err != nil {
		return nil, err
	}

	return &member, nil
}
