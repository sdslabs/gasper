package types

import (
	"crypto/rsa"
	"fmt"
	"math"
	"time"
)

// Application is the interface for creating an application
type Application interface {
	GetName() string
	GetGitRepositoryURL() string
	GetGitRepositoryBranch() string
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
	GetNameServers() []string
	GetDockerImage() string
	SetContainerID(id string)
	GetContainerID() string
	SetContainerPort(port int)
	GetContainerPort() int
	HasConfGenerator() bool
	InvokeConfGenerator(name, index string) string
}

// Git stores the information related to the application's git repository
type Git struct {
	RepoURL     string `json:"repo_url" bson:"repo_url" valid:"required~Field 'repo_url' inside field 'git' is required but was not provided,url~Field 'repo_url' inside field 'git' is not a valid URL"`
	AccessToken string `json:"access_token,omitempty" bson:"access_token,omitempty"`
	Branch      string `json:"branch,omitempty" bson:"branch,omitempty"`
}

// Context stores the information related to building and running an application
type Context struct {
	Index  string   `json:"index" bson:"index" valid:"required~Field 'index' inside field 'context' was required but was not provided"`
	Port   int      `json:"port" bson:"port" valid:"port~Field 'port' inside field 'context' is not a valid port"`
	RcFile bool     `json:"rc_file" bson:"rc_file"`
	Build  []string `json:"build,omitempty" bson:"build,omitempty"`
	Run    []string `json:"run,omitempty" bson:"run,omitempty"`
}

// Resources defines the resources requested by an application
type Resources struct {
	// Memory limits in GB
	Memory float64 `json:"memory" bson:"memory" valid:"float~Field 'memory' inside field 'resources' should be of type float"`

	// CPU quota in units of CPUs
	CPU float64 `json:"cpu" bson:"cpu" valid:"float~Field 'cpu' inside field 'resources' should be of type float"`
}

type RepositoryRequest struct {
	Name string `json:"name" bson:"name" valid:"required~Field 'name' is required but was not provided,alphanum~Field 'name' should only have alphanumeric characters,stringlength(3|40)~Field 'name' should have length between 3 to 40 characters,lowercase~Field 'name' should have only lowercase characters"`
}

type EncryptKey struct {
	PublicKey rsa.PublicKey `json:"public_key"`
}

// ApplicationConfig is the configuration required for creating an application
type ApplicationConfig struct {
	Name          string                      `json:"name" bson:"name" valid:"required~Field 'name' is required but was not provided,alphanum~Field 'name' should only have alphanumeric characters,stringlength(3|40)~Field 'name' should have length between 3 to 40 characters,lowercase~Field 'name' should have only lowercase characters"`
	Password      string                      `json:"password" bson:"password" valid:"required~Field 'password' is required but was not provided"`
	Git           Git                         `json:"git" bson:"git"`
	Context       Context                     `json:"context" bson:"context"`
	Resources     Resources                   `json:"resources,omitempty" bson:"resources,omitempty"`
	Env           M                           `json:"env,omitempty" bson:"env,omitempty"`
	NameServers   []string                    `json:"name_servers,omitempty" bson:"name_servers,omitempty"`
	DockerImage   string                      `json:"docker_image" bson:"docker_image"`
	ContainerID   string                      `json:"container_id" bson:"container_id"`
	ContainerPort int                         `json:"container_port" bson:"container_port"`
	ConfGenerator func(string, string) string `json:"-" bson:"-"`
	Language      string                      `json:"language" bson:"language"`
	InstanceType  string                      `json:"instance_type" bson:"instance_type"`
	CloudflareID  string                      `json:"cloudflare_id,omitempty" bson:"cloudflare_id,omitempty"`
	AppURL        string                      `json:"app_url,omitempty" bson:"app_url,omitempty"`
	HostIP        string                      `json:"host_ip,omitempty" bson:"host_ip,omitempty"`
	PublicIP      string                      `json:"public_ip,omitempty" bson:"public_ip,omitempty"`
	SSHCmd        string                      `json:"ssh_cmd,omitempty" bson:"ssh_cmd,omitempty"`
	Owner         string                      `json:"owner,omitempty" bson:"owner,omitempty"`
	Datetime      time.Time                   `json:"datetime" bson:"datetime"`
	Success       bool                        `json:"success,omitempty" bson:"-"`
}

type ApplicationRemote struct {
	GitURL string `json:"giturl" bson:"giturl"`
}

type AccessToken struct {
	// PAT for pushing code to repository
	PAT string `json:"pat" bson:"pat"`
	// Username of Gasper Github user
	Username string `json:"username" bson:"username"`
	// Email id of Gasper Github user
	Email string `json:"email" bson:"email"`
}

// GetName returns the application's name
func (app *ApplicationConfig) GetName() string {
	return app.Name
}

// GetGitRepositoryURL returns the application's git repository URL
func (app *ApplicationConfig) GetGitRepositoryURL() string {
	return app.Git.RepoURL
}

// GetGitRepositoryBranch returns the branch to clone from the application's git repository
// Default branch is `master`
func (app *ApplicationConfig) GetGitRepositoryBranch() string {
	if app.Git.Branch == "" {
		return "master"
	}
	return app.Git.Branch
}

// HasGitAccessToken checks whether access token is required for cloning
// the application's git repository
func (app *ApplicationConfig) HasGitAccessToken() bool {
	if app.Git.AccessToken == "" {
		return false
	}
	return true
}

// GetGitAccessToken returns the application's git access token
func (app *ApplicationConfig) GetGitAccessToken() string {
	return app.Git.AccessToken
}

// GetIndex returns the index file required for starting the application
func (app *ApplicationConfig) GetIndex() string {
	return app.Context.Index
}

// GetApplicationPort returns the port on which the application runs
func (app *ApplicationConfig) GetApplicationPort() int {
	if app.Context.Port == 0 {
		app.Context.Port = 80
	}
	return app.Context.Port
}

// HasRcFile checks if a Run Commands file is required for building and
// running the application
func (app *ApplicationConfig) HasRcFile() bool {
	return app.Context.RcFile
}

// GetBuildCommands returns the shell commands used for building the application's dependencies
func (app *ApplicationConfig) GetBuildCommands() []string {
	return app.Context.Build
}

// GetRunCommands returns the shell commands used for running the application
func (app *ApplicationConfig) GetRunCommands() []string {
	return app.Context.Run
}

// GetCPULimit returns application's CPU Limit in units of NanoCPUs
func (app *ApplicationConfig) GetCPULimit() int64 {
	if app.Resources.CPU == 0 {
		app.Resources.CPU = DefaultCPUs
	}
	return int64(app.Resources.CPU * math.Pow(10, 9))
}

// GetMemoryLimit returns application's Memory Limit in units of bytes
func (app *ApplicationConfig) GetMemoryLimit() int64 {
	if app.Resources.Memory == 0 {
		app.Resources.Memory = DefaultMemory
	}
	return int64(app.Resources.Memory * math.Pow(1024, 3))
}

// GetEnvVars returns the environment variables to be used inside the docker container
func (app *ApplicationConfig) GetEnvVars() map[string]interface{} {
	return app.Env
}

// SetNameServers sets the DNS NameServers to be used by the application's docker container
// in the application's context
func (app *ApplicationConfig) SetNameServers(servers []string) {
	app.NameServers = servers
}

// AddNameServers adds a DNS NameServer to be used by the application's docker container
// in the application's context
func (app *ApplicationConfig) AddNameServers(servers ...string) {
	app.NameServers = append(app.NameServers, servers...)
}

// GetNameServers returns the DNS NameServers to be used by the application's docker container
func (app *ApplicationConfig) GetNameServers() []string {
	return app.NameServers
}

// SetDockerImage defines the docker image to be used for creating the container
func (app *ApplicationConfig) SetDockerImage(image string) {
	app.DockerImage = image
}

// GetDockerImage returns the docker image used for creating container
func (app *ApplicationConfig) GetDockerImage() string {
	return app.DockerImage
}

// SetContainerID sets docker container ID in the application's context
func (app *ApplicationConfig) SetContainerID(id string) {
	app.ContainerID = id
}

// GetContainerID returns the docker container ID in the application's context
func (app *ApplicationConfig) GetContainerID() string {
	return app.ContainerID
}

// SetContainerPort sets the port to which the container will be bound to
// in the host system
func (app *ApplicationConfig) SetContainerPort(port int) {
	app.ContainerPort = port
}

// GetContainerPort returns the port to which the container is bound in the
// host system
func (app *ApplicationConfig) GetContainerPort() int {
	return app.ContainerPort
}

// SetConfGenerator defines a config generator used for applications using nginx
// Ex :- PHP and Static applications
func (app *ApplicationConfig) SetConfGenerator(gen func(string, string) string) {
	app.ConfGenerator = gen
}

// HasConfGenerator checks whether a config generator is required for bootstraping
// the application
func (app *ApplicationConfig) HasConfGenerator() bool {
	if app.ConfGenerator == nil {
		return false
	}
	return true
}

// InvokeConfGenerator invokes the config generator
func (app *ApplicationConfig) InvokeConfGenerator(name, index string) string {
	return app.ConfGenerator(name, index)
}

// SetLanguage sets the application's language in its context
func (app *ApplicationConfig) SetLanguage(language string) {
	app.Language = language
}

// SetInstanceType sets the application's type of instance in its context
func (app *ApplicationConfig) SetInstanceType(instanceType string) {
	app.InstanceType = instanceType
}

// SetCloudflareID sets the application's cloudflare record ID in its context
func (app *ApplicationConfig) SetCloudflareID(cloudflareID string) {
	app.CloudflareID = cloudflareID
}

// SetAppURL sets the application's domain URL in its context
func (app *ApplicationConfig) SetAppURL(appURL string) {
	app.AppURL = appURL
}

// SetSuccess defines the success of deploying the application
func (app *ApplicationConfig) SetSuccess(success bool) {
	app.Success = success
}

// SetHostIP sets the IP address of the host in which the application is deployed
// in its context
func (app *ApplicationConfig) SetHostIP(IP string) {
	app.HostIP = IP
}

// SetPublicIP sets the public IP address of the host in which the application is deployed
// in its context
func (app *ApplicationConfig) SetPublicIP(IP string) {
	app.PublicIP = IP
}

// SetSSHCmd generates the command to SSH into an application's docker container
// for the information of the client
func (app *ApplicationConfig) SetSSHCmd(port int, appName, IP string) {
	app.SSHCmd = fmt.Sprintf("ssh -p %d %s@%s", port, appName, IP)
}

// SetOwner sets the owner of the application in its context
// The owner is referenced by his/her email ID
func (app *ApplicationConfig) SetOwner(owner string) {
	app.Owner = owner
}

// SetDateTime sets the date on which the db was created
func (app *ApplicationConfig) SetDateTime() {
	app.Datetime = time.Now()
}
