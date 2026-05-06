package mijnhostapi

// Contents of the structs below were generated from <https://mijn.host/api/doc/api-3563906>, their names have been adapted.

type GetDNSRecordsForDomainResponse struct {
	Data GetDNSRecordsForDomainResponseData `json:"data"`
	// Result code
	Status int64 `json:"status"`
	// Result information
	StatusDescription string `json:"status_description"`
}

type GetDNSRecordsForDomainResponseData struct {
	// Domain name
	Domain string `json:"domain"`
	// DNS records
	Records []Record `json:"records"`
}

type Record struct {
	Name  string `json:"name"`
	TTL   int64  `json:"ttl"`
	Type  string `json:"type"`
	Value string `json:"value"`
}
