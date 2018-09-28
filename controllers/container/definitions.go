package container

type createConfig struct {
	Image string `form:"image" json:"image" binding:"required"`
}
