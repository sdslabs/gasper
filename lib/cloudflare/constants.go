package cloudflare

import "github.com/sdslabs/gasper/configs"

const (
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
