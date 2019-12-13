package cloudflare

import "github.com/sdslabs/gasper/configs"

const (
	// ApplicationInstance is the identifier attached with an application's DNS entry in cloudflare
	ApplicationInstance = "app"

	// DatabaseInstance is the identifier attached with a database's DNS entry in cloudflare
	DatabaseInstance = "db"

	baseEndpoint         = "https://api.cloudflare.com/client/v4"
	listZonesEndpoint    = baseEndpoint + "/zones"
	fetchRecordEndpoint  = listZonesEndpoint + "/%s/dns_records"
	createRecordEndpoint = listZonesEndpoint + "/%s/dns_records"
	updateRecordEndpoint = listZonesEndpoint + "/%s/dns_records/%s"
	deleteRecordEndpoint = listZonesEndpoint + "/%s/dns_records/%s"
)

var (
	token    = configs.CloudflareConfig.Token
	domain   = configs.GasperConfig.Domain
	publicIP = configs.CloudflareConfig.PublicIP
)
