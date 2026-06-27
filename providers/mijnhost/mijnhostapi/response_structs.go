package mijnhostapi

type UpdateDNSRecordForDomainResponse struct {
	StatusWithDescription
}

type StatusWithDescription struct {
	Status               int64  `json:"status"`
	StatusDescription    string `json:"status_description"`
}
