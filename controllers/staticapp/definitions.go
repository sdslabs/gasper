package staticapp

// appConfig is json binding config for creating new static page
type appConfig struct {
	Name      string `json:"name" form:"name" binding:"required"`
	UserID    int    `json:"user_id" form:"user_id" binding:"required"`
	GithubURL string `json:"github_url" form:"github_url" binding:"required"`
}

// app is an interface containing methods for static pages
type app interface {
	ReadAndWriteConfig() error
}
