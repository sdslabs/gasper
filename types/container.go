package types

// ApplicationContainer is the configuration for creating a container
// for running an application
type ApplicationContainer struct {
	Name            string
	Image           string
	ApplicationPort int
	ContainerPort   int
	WorkDir         string
	StoreDir        string
	Env             M
	Memory          int64
	CPU             int64
	NameServers     []string
}
