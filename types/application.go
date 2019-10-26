package types

import (
	"math"
)

// Application is the interface for creating an application
type Application interface {
	GetName() string
	GetGitRepositoryURL() string
	HasGitAccessToken() bool
	GetGitAccessToken() string
	GetIndex() string
	GetApplicationPort() int
	HasRcFile() bool
	GetBuildCommands() []string
	GetRunCommands() []string
	GetCPULimit() int64
	GetMemoryLimit() int64
	GetEnvVars() map[string]interface{}
	SetDockerImage(image string)
	GetDockerImage() string
	SetContainerID(id string)
	GetContainerID() string
	SetContainerPort(port int)
	GetContainerPort() int
	SetConfGenerator(gen func(string, string) string)
	HasConfGenerator() bool
	InvokeConfGenerator(name, index string) string
}

// Context stores the information related to building and running an application
type Context struct {
	Index  string   `json:"index" bson:"index" valid:"required~Field 'index' inside field 'context' was required but was not provided"`
	Port   int      `json:"port" bson:"port,omitempty" valid:"port~Field 'port' inside field 'context' is not a valid port"`
	RcFile bool     `json:"rc_file" bson:"rc_file,omitempty"`
	Build  []string `json:"build" bson:"build,omitempty"`
	Run    []string `json:"run" bson:"run,omitempty"`
}

// Resources defines the resources requested by an application
type Resources struct {
	// Memory limits in GB
	Memory float64 `json:"memory" bson:"memory,omitempty" valid:"float~Field 'memory' inside field 'resources' should be of type float"`

	// CPU quota in units of CPUs
	CPU float64 `json:"cpu" bson:"cpu,omitempty" valid:"float~Field 'cpu' inside field 'resources' should be of type float"`
}

// ApplicationConfig is the configuration required for creating an application
type ApplicationConfig struct {
	Name           string                      `json:"name" bson:"name" valid:"required~Field 'name' is required but was not provided,alphanum~Field 'name' should only have alphanumeric characters,stringlength(3|40)~Field 'name' should have length between 3 to 40 characters,lowercase~Field 'name' should have only lowercase characters"`
	Password       string                      `json:"password" bson:"password" valid:"required~Field 'password' is required but was not provided"`
	GitURL         string                      `json:"git_url" bson:"git_url" valid:"required~Field 'git_url' is required but was not provided,url~Field 'git_url' is not a valid URL"`
	GitAccessToken string                      `json:"git_access_token,omitempty" bson:"git_access_token,omitempty"`
	Context        Context                     `json:"context" bson:"context"`
	Resources      Resources                   `json:"resources,omitempty" bson:"resources,omitempty"`
	Env            M                           `json:"env,omitempty" bson:"env,omitempty"`
	DockerImage    string                      `json:"docker_image" bson:"docker_image"`
	ContainerID    string                      `json:"container_id" bson:"container_id"`
	ContainerPort  int                         `json:"container_port" bson:"container_port"`
	ConfGenerator  func(string, string) string `json:",omitempty" bson:",omitempty"`
	Language       string                      `json:"language" bson:"language"`
	InstanceType   string                      `json:"instance_type" bson:"instance_type"`
	Rebuild        bool                        `json:"rebuild,omitempty" bson:",omitempty"`
	CloudflareID   string                      `json:"cloudflare_id,omitempty" bson:"cloudflare_id,omitempty"`
	AppURL         string                      `json:"app_url,omitempty" bson:"app_url,omitempty"`
	HostIP         string                      `json:"host_ip,omitempty" bson:"host_ip,omitempty"`
	Owner          string                      `json:"owner,omitempty" bson:"owner,omitempty"`
	Success        bool                        `json:"success,omitempty"`
}

func (app *ApplicationConfig) GetName() string {
	return app.Name
}

func (app *ApplicationConfig) GetGitRepositoryURL() string {
	return app.GitURL
}

func (app *ApplicationConfig) HasGitAccessToken() bool {
	if app.GitAccessToken == "" {
		return false
	}
	return true
}

func (app *ApplicationConfig) GetGitAccessToken() string {
	return app.GitAccessToken
}

func (app *ApplicationConfig) GetIndex() string {
	return app.Context.Index
}

func (app *ApplicationConfig) GetApplicationPort() int {
	if app.Context.Port == 0 {
		app.Context.Port = 80
	}
	return app.Context.Port
}

func (app *ApplicationConfig) HasRcFile() bool {
	return app.Context.RcFile
}

func (app *ApplicationConfig) GetBuildCommands() []string {
	return app.Context.Build
}

func (app *ApplicationConfig) GetRunCommands() []string {
	return app.Context.Run
}

// GetCPULimit returns Application CPU Limit in units of NanoCPUs
func (app *ApplicationConfig) GetCPULimit() int64 {
	if app.Resources.CPU == 0 {
		app.Resources.CPU = DefaultCPUs
	}
	return int64(app.Resources.CPU * math.Pow(10, 9))
}

// GetMemoryLimit returns Application Memory Limit in units of bytes
func (app *ApplicationConfig) GetMemoryLimit() int64 {
	if app.Resources.Memory == 0 {
		app.Resources.Memory = DefaultMemory
	}
	return int64(app.Resources.Memory * math.Pow(1024, 3))
}

func (app *ApplicationConfig) GetEnvVars() map[string]interface{} {
	return app.Env
}

func (app *ApplicationConfig) SetDockerImage(image string) {
	app.DockerImage = image
}

func (app *ApplicationConfig) GetDockerImage() string {
	return app.DockerImage
}

func (app *ApplicationConfig) SetContainerID(id string) {
	app.ContainerID = id
}

func (app *ApplicationConfig) GetContainerID() string {
	return app.ContainerID
}

func (app *ApplicationConfig) SetContainerPort(port int) {
	app.ContainerPort = port
}

func (app *ApplicationConfig) GetContainerPort() int {
	return app.ContainerPort
}

func (app *ApplicationConfig) SetConfGenerator(gen func(string, string) string) {
	app.ConfGenerator = gen
}

func (app *ApplicationConfig) HasConfGenerator() bool {
	if app.ConfGenerator == nil {
		return false
	}
	return true
}

func (app *ApplicationConfig) InvokeConfGenerator(name, index string) string {
	return app.ConfGenerator(name, index)
}

func (app *ApplicationConfig) SetLanguage(language string) {
	app.Language = language
}

func (app *ApplicationConfig) SetInstanceType(instanceType string) {
	app.InstanceType = instanceType
}

func (app *ApplicationConfig) HasRebuildEnabled() bool {
	return app.Rebuild
}

func (app *ApplicationConfig) DisableRebuild() {
	app.Rebuild = false
}

func (app *ApplicationConfig) SetCloudflareID(cloudflareID string) {
	app.CloudflareID = cloudflareID
}

func (app *ApplicationConfig) SetAppURL(appURL string) {
	app.AppURL = appURL
}

func (app *ApplicationConfig) SetSuccess(success bool) {
	app.Success = success
}

func (app *ApplicationConfig) SetHostIP(IP string) {
	app.HostIP = IP
}
