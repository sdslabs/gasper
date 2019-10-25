package configs

import (
	"time"

	"github.com/sdslabs/gasper/types"
)

// Admin is the configuration for the default Admin
type Admin struct {
	Email    string `toml:"email"`
	Password string `toml:"password"`
	Username string `toml:"username"`
}

// Cloudflare is the configuration for cloudflare API
type Cloudflare struct {
	PlugIn   bool   `toml:"plugin"`
	PublicIP string `toml:"public_ip"`
	Token    string `toml:"api_token"`
}

// Mongo is the configuration for mongodb storage
type Mongo struct {
	URL string `toml:"url"`
}

// Redis is the configuration for redis storage
type Redis struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

// Falcon is the configuration for SDSLabs oauth2
type Falcon struct {
	PlugIn                        bool   `toml:"plugin"`
	FalconClientID                string `toml:"falcon_client_id"`
	FalconClientSecret            string `toml:"falcon_client_secret"`
	FalconURLAccessToken          string `toml:"falcon_access_token_url"`
	FalconURLResourceOwnerDetails string `toml:"falcon_resource_owner_url"`
	FalconAccountsURL             string `toml:"falcon_accounts_url"`
	RedirectURI                   string `toml:"redirect_uri"`
}

// GenericService is the default configuration for all services
type GenericService struct {
	Deploy bool `toml:"deploy"`
	Port   int  `toml:"port"`
}

// DominusService is the default configuration for Dominus microservice
type DominusService struct {
	GenericService
	CleanupInterval time.Duration `toml:"cleanup_interval"`
}

// SSHProxyCfg is the configuration for SSH_Proxy plugin
type SSHProxyCfg struct {
	PlugIn bool `toml:"plugin"`
	Port   int  `toml:"port"`
}

// SSHService is the configuration for SSH microservice
type SSHService struct {
	GenericService
	HostSigners     []string    `toml:"host_signers"`
	UsingPassphrase bool        `toml:"using_passphrase"`
	Passphrase      string      `toml:"passphrase"`
	Proxy           SSHProxyCfg `toml:"proxy"`
}

// SSLConfig is the configuration for SSL in Enrai microservice
type SSLConfig struct {
	PlugIn      bool   `toml:"plugin"`
	Port        int    `toml:"port"`
	Certificate string `toml:"certificate"`
	PrivateKey  string `toml:"private_key"`
}

// EnraiService is the configuration for Enrai microservice
type EnraiService struct {
	GenericService
	SSL SSLConfig `toml:"ssl"`
}

// HikariService is the configuration for Hikari microservice
type HikariService struct {
	GenericService
	RecordUpdateInterval time.Duration `toml:"record_update_interval"`
}

// MysqlService is the configuration for Mysql microservice
type MysqlService struct {
	GenericService
	ContainerPort int     `toml:"container_port"`
	Env           types.M `toml:"env"`
}

// MongodbService is the configuration for Mongodb microservice
type MongodbService struct {
	GenericService
	ContainerPort int     `toml:"container_port"`
	Env           types.M `toml:"env"`
}

// Images is the configuration for the docker images in use
type Images struct {
	Static  string `toml:"static"`
	Php     string `toml:"php"`
	Nodejs  string `toml:"nodejs"`
	Python2 string `toml:"python2"`
	Python3 string `toml:"python3"`
	Mysql   string `toml:"mysql"`
	Mongodb string `toml:"mongodb"`
}

// Services is the configuration for all Services
type Services struct {
	ExposureInterval time.Duration  `toml:"exposure_interval"`
	Dominus          DominusService `toml:"dominus"`
	Mizu             GenericService `toml:"mizu"`
	SSH              SSHService     `toml:"ssh"`
	Enrai            EnraiService   `toml:"enrai"`
	Hikari           HikariService  `toml:"hikari"`
	Mysql            MysqlService   `toml:"mysql"`
	Mongodb          MongodbService `toml:"mongodb"`
}

// GasperCfg is the configuration for the entire project
type GasperCfg struct {
	Debug       bool       `toml:"debug"`
	Domain      string     `toml:"domain"`
	Secret      string     `toml:"secret"`
	ProjectRoot string     `toml:"project_root"`
	RcFile      string     `toml:"rc_file"`
	OfflineMode bool       `toml:"offline_mode"`
	DNSServers  []string   `toml:"dns_servers"`
	Admin       Admin      `toml:"admin"`
	Cloudflare  Cloudflare `toml:"cloudflare"`
	Mongo       Mongo      `toml:"mongo"`
	Redis       Redis      `toml:"redis"`
	Falcon      Falcon     `toml:"falcon"`
	Images      Images     `toml:"images"`
	Services    Services   `toml:"services"`
}
