package types

// ApplicationContainer is the configuration for creating a container
// for running an application
type ApplicationContainer struct {
	// Name of the container
	Name string
	// Docker image used for creating the container
	Image string
	// Port on which an application is running inside the container
	ApplicationPort int
	// Port of the docker container in the host system
	ContainerPort int
	// Directory inside the docker container for volume mounting purposes
	WorkDir string
	// Directory on the host system for volume mounting purposes
	StoreDir string
	// Environment variables
	Env M
	// Resource limits
	Memory int64
	CPU    int64
	// DNS nameservers to be used for domain name resolutions within the container
	NameServers []string
}

// DatabaseContainer is the configuration for creating a container
// for running an database service
type DatabaseContainer struct {
	// Name of the container
	Name string
	// Docker image used for creating the container
	Image string
	// Port on which a database service is running inside the container
	DatabasePort int
	// Port of the docker container in the host system
	ContainerPort int
	// Directory inside the docker container for volume mounting purposes
	WorkDir string
	// Directory on the host system for volume mounting purposes
	StoreDir string
	// Custom commands to be executed on a container's startup
	Cmd []string
	// Environment variables
	Env M
}

// LizardfsContainer is the configuration for creating a container
// for running the filesystem
type LizardfsContainer struct {
	// Name of the container
	Name string
	// Docker image used for creating the container
	Image string
	// Port on which a database service is running inside the container
	HostPort1 int
	// Port on which a database service is running inside the container
	HostPort2 int
	// Port of the docker container in the host system
	ContainerPort1 int
	// Port of the docker container in the host system
	ContainerPort2 int
	// Directory inside the docker container for volume mounting purposes
	WorkDir string
	// Directory on the host system for volume mounting purposes
	StoreDir string
	// Custom commands to be executed on a container's startup
	Cmd []string
	// Environment variables
	Env M
}

// HasCustomCMD checks whether a database container needs custom CMD commands on boot
func (containerCfg *DatabaseContainer) HasCustomCMD() bool {
	return len(containerCfg.Cmd) > 0
}
