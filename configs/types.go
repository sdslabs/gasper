package configs

import (
	"time"

	"github.com/sdslabs/gasper/types"
)

// JWT is the configuration for auth token
type JWT struct {
	Timeout    time.Duration `toml:"timeout"`
	MaxRefresh time.Duration `toml:"max_refresh"`
}

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

// GenericService is the default configuration for all services
type GenericService struct {
	Deploy bool `toml:"deploy"`
	Port   int  `toml:"port"`
}

// AppMakerService is the default configuration for appmaker microservice
type AppMakerService struct {
	GenericService
	MetricsInterval time.Duration `toml:"metrics_interval"`
	HealthInterval  time.Duration `toml:"health_interval"`
	AppLimit        int           `toml:"app_limit"`
}

// MasterService is the default configuration for Master microservice
type MasterService struct {
	GenericService
	CleanupInterval time.Duration   `toml:"cleanup_interval"`
	MongoDB         DatabaseService `toml:"mongodb"`
	Redis           DatabaseService `toml:"redis"`
}

// GenSSHService is the configuration for GenSSH microservice
type GenSSHService struct {
	GenericService
	HostSigners     []string `toml:"host_signers"`
	UsingPassphrase bool     `toml:"using_passphrase"`
	Passphrase      string   `toml:"passphrase"`
	EntrypointIP    string   `toml:"entrypoint_ip"`
}

// SSLConfig is the configuration for SSL in GenProxy microservice
type SSLConfig struct {
	PlugIn      bool   `toml:"plugin"`
	Port        int    `toml:"port"`
	Certificate string `toml:"certificate"`
	PrivateKey  string `toml:"private_key"`
}

// GenProxyService is the configuration for GenProxy microservice
type GenProxyService struct {
	GenericService
	SSL                  SSLConfig     `toml:"ssl"`
	RecordUpdateInterval time.Duration `toml:"record_update_interval"`
}

// GenDNSService is the configuration for GenDNS microservice
type GenDNSService struct {
	GenericService
	RecordUpdateInterval time.Duration `toml:"record_update_interval"`
}

// DatabaseService is the configuration for database servers
type DatabaseService struct {
	PlugIn        bool    `toml:"plugin"`
	ContainerPort int     `toml:"container_port"`
	Env           types.M `toml:"env"`
	Password      string  `toml:"password"`
}

// DbMakerService is the configuration for DbMaker microservice
type DbMakerService struct {
	GenericService
	MySQL      DatabaseService `toml:"mysql"`
	MongoDB    DatabaseService `toml:"mongodb"`
	PostgreSQL DatabaseService `toml:"postgresql"`
	Redis      DatabaseService `toml:"redis"`
	DBLimit    int             `toml:"db_limit"`
}

// JikanService is the configuration for Jikan microservice
type JikanService struct {
	GenericService
}

// Images is the configuration for the docker images in use
type Images struct {
	Static     string `toml:"static"`
	Php        string `toml:"php"`
	Nodejs     string `toml:"nodejs"`
	Python2    string `toml:"python2"`
	Python3    string `toml:"python3"`
	Golang     string `toml:"golang"`
	Ruby       string `toml:"ruby"`
	Rust       string `toml:"rust"`
	Mysql      string `toml:"mysql"`
	Mongodb    string `toml:"mongodb"`
	Postgresql string `toml:"postgresql"`
	Redis      string `toml:"redis"`
}

// Services is the configuration for all Services
type Services struct {
	ExposureInterval time.Duration   `toml:"exposure_interval"`
	RateInterval     time.Duration   `toml:"rate_interval"`
	RateLimit        int             `toml:"rate_limit"`
	Master           MasterService   `toml:"master"`
	AppMaker         AppMakerService `toml:"appmaker"`
	GenSSH           GenSSHService   `toml:"genssh"`
	GenProxy         GenProxyService `toml:"genproxy"`
	GenDNS           GenDNSService   `toml:"gendns"`
	DbMaker          DbMakerService  `toml:"dbmaker"`
	Jikan            JikanService    `toml:"jikan"`
}

type Github struct {
	Username string `toml:"username"`
	Email    string `toml:"email"`
	PAT      string `toml:"pat"`
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
	JWT         JWT        `toml:"jwt"`
	Admin       Admin      `toml:"admin"`
	Cloudflare  Cloudflare `toml:"cloudflare"`
	Mongo       Mongo      `toml:"mongo"`
	Redis       Redis      `toml:"redis"`
	Images      Images     `toml:"images"`
	Services    Services   `toml:"services"`
	Github      Github     `toml:"github"`
}
