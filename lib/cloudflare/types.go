package cloudflare

type errorResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

type zoneRecord struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Status              string   `json:"status"`
	NameServers         []string `json:"name_servers"`
	OriginalNameServers []string `json:"original_name_servers"`
}

type dnsRecord struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	ZoneID   string `json:"zone_id"`
	ZoneName string `json:"zone_name"`
}

// GenericResponse is the common response from Cloudflare API
type GenericResponse struct {
	Success bool            `json:"success"`
	Errors  []errorResponse `json:"errors"`
}

// ZoneResponse stores the details of a zone
type ZoneResponse struct {
	Result []zoneRecord `json:"result"`
	GenericResponse
}

// MultiResponse stores details of multiple DNS records
type MultiResponse struct {
	Result []dnsRecord `json:"result"`
	GenericResponse
}

// SingleResponse stores details of a single DNS record
type SingleResponse struct {
	Result dnsRecord `json:"result"`
	GenericResponse
}

// singlePayload is the request body for creating a new DNS record
type singlePayload struct {
	// DNS record type. A in our case
	Type string `json:"type,omitempty"`
	// Name of the application
	Name string `json:"name,omitempty"`
	// IP address of the deployed application
	Content string `json:"content,omitempty"`
}
