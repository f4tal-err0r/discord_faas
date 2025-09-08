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

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type BuildFunc struct {
	Name        string     `json:"name"`
	Runtime     string     `json:"runtime"`
	GuildID     string     `json:"guildID"`
	Version     string     `json:"version"`
	Description string     `json:"description"`
	Commands    []Commands `json:"commands"`
	Envvars     []EnvVar   `json:"envvars"`
}

type ContextResp struct {
	ClientID       string `json:"clientID"`
	GuildID        string `json:"guildID"`
	GuildName      string `json:"guildName"`
	ServerURL      string `json:"serverURL"`
	CurrentContext bool   `json:"currentContext"`
	JWToken        string `json:"jwToken"`
}

type Runtimes struct {
	Name          string   `json:"name" yaml:"name"`
	Lang          string   `json:"lang" yaml:"lang"`
	Version       string   `json:"version" yaml:"version"`
	ReadOnlyFiles []string `json:"readOnlyFiles" yaml:"readOnlyFiles"`
	Build         struct {
		Image    string   `json:"image" yaml:"image"`
		Commands []string `json:"commands" yaml:"commands"`
	} `json:"build" yaml:"build"`
	Run struct {
		Image    string   `json:"image" yaml:"image"`
		Commands []string `json:"commands" yaml:"commands"`
	} `json:"run" yaml:"run"`
}
