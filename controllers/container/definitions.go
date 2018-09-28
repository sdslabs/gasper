package container

type dockerConfig struct {
	Image string `form:"image" json:"image" binding:"required"`
}
