package container

// Config is a struct definition for container controller
type Config struct {
	Image string `form:"image" json:"image" binding:"required"`
}
