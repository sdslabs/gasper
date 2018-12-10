package static

// staticAppConfig is json binding config for creating new static page
type staticAppConfig struct {
	Name      string `json:"name" form:"name" binding:"required"`
	UserID    int    `json:"user_id" form:"user_id" binding:"required"`
	GithubURL string `json:"github_url" form:"github_url" binding:"required"`
}
