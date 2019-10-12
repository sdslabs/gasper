package configs

// Admin is the configuration for the default Admin
type Admin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

// Cloudflare is the configuration for cloudflare API
type Cloudflare struct {
	PlugIn   bool   `json:"plugIn"`
	PublicIP string `json:"publicIP"`
	Token    string `json:"token"`
}

// Cron is the configuration for cronjobs
type Cron struct {
	CleanupInterval  int `json:"cleanupInterval"`
	ExposureInterval int `json:"exposureInterval"`
}

// Mongo is the configuration for mongodb storage
type Mongo struct {
	URL string `json:"url"`
}

// Redis is the configuration for redis storage
type Redis struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"DB"`
}

// Falcon is the configuration for SDSLabs oauth2
type Falcon struct {
	PlugIn                        bool   `json:"plugIn"`
	FalconClientID                string `json:"falconClientId"`
	FalconClientSecret            string `json:"falconClientSecret"`
	FalconURLAccessToken          string `json:"falconUrlAccessToken"`
	FalconURLResourceOwnerDetails string `json:"falconUrlResourceOwnerDetails"`
	FalconAccountsURL             string `json:"falconAccountsUrl"`
	RedirectURI                   string `json:"redirectUri"`
}

// GenericService is the default configuration for all services
type GenericService struct {
	Deploy bool `json:"deploy"`
	Port   int  `json:"port"`
}

// SSH is the configuration for SSH and SSH_Proxy service
type SSH struct {
	GenericService
	HostSigners     []string `json:"host_signers"`
	UsingPassphrase bool     `json:"using_passphrase"`
	Passphrase      string   `json:"passphrase"`
}

// Mysql is the configuration for Mysql service
type Mysql struct {
	GenericService
	ContainerPort int                    `json:"container_port"`
	Env           map[string]interface{} `json:"env"`
}

// Mongodb is the configuration for Mongodb service
type Mongodb struct {
	GenericService
	ContainerPort int                    `json:"container_port"`
	Env           map[string]interface{} `json:"env"`
}

// Images is the configuration for the images in use
type Images struct {
	Static  string `json:"static"`
	Php     string `json:"php"`
	Nodejs  string `json:"nodejs"`
	Python2 string `json:"python2"`
	Python3 string `json:"python3"`
	Mysql   string `json:"mysql"`
	Mongodb string `json:"mongodb"`
}

// Services is the configuration for all Services
type Services struct {
	Dominus  GenericService `json:"dominus"`
	Mizu     GenericService `json:"mizu"`
	SSH      SSH            `json:"ssh"`
	SSHProxy SSH            `json:"ssh_proxy"`
	Enrai    GenericService `json:"enrai"`
	Mysql    Mysql          `json:"mysql"`
	Mongodb  Mongodb        `json:"mongodb"`
}

// GasperCfg is the configuration for the entire project
type GasperCfg struct {
	Debug       bool       `json:"debug"`
	Domain      string     `json:"domain"`
	Secret      string     `json:"secret"`
	ProjectRoot string     `json:"projectRoot"`
	RcFile      string     `json:"rcFile"`
	OfflineMode bool       `json:"offlineMode"`
	DNSServers  []string   `json:"dnsServers"`
	Admin       Admin      `json:"admin"`
	Cloudflare  Cloudflare `json:"cloudflare"`
	Cron        Cron       `json:"cron"`
	Mongo       Mongo      `json:"mongo"`
	Redis       Redis      `json:"redis"`
	Falcon      Falcon     `json:"falcon"`
	Images      Images     `json:"images"`
	Services    Services   `json:"services"`
}
