package v1

type GetContext struct {
	Token   string `json:"token"`
	GuildID string `json:"guildID"`
}

type Args struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

type Commands struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Args        []*Args `json:"args"`
}

type BuildFunc struct {
	Name        string      `json:"name"`
	Runtime     string      `json:"runtime"`
	GuildID     string      `json:"guildID"`
	Version     string      `json:"version"`
	Description string      `json:"description"`
	Commands    []*Commands `json:"commands"`
}

type ContextResp struct {
	ClientID       string `json:"clientID"`
	GuildID        string `json:"guildID"`
	GuildName      string `json:"guildName"`
	ServerURL      string `json:"serverURL"`
	CurrentContext bool   `json:"currentContext"`
	JWToken        string `json:"jwToken"`
}
