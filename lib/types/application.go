package types

// ApplicationConfig defines the config for various service apps
type ApplicationConfig struct {
	DockerImage  string
	ConfFunction func(string, string) string
}

// Request configs...

// StaticAppConfig defined the request structure for creating new static app
type StaticAppConfig struct {
	Name      string `json:"name" form:"name" binding:"required"`
	UserID    int    `json:"user_id" form:"user_id" binding:"required"`
	GithubURL string `json:"github_url" form:"github_url" binding:"required"`
}
