package mijnhostapi

// Contents of the structs below were generated from <https://mijn.host/api/doc/api-3563900>, their names have been adapted.

type GetDomainResponse struct {
	// Result data
	Data GetDomainResponseData `json:"data"`
	// Result code
	Status int64 `json:"status"`
	// Result information
	StatusDescription string `json:"status_description"`
}

type GetDomainResponseData struct {
	// DNSSEC enabled
	DnssecEnabled int64 `json:"dnssec_enabled"`
	// DNSSEC keys
	DnssecKeys []DnssecKey `json:"dnssec_keys"`
	// Domain name
	Domain string `json:"domain"`
	// Single site-wide redirect when the forwarder is enabled and configured; `null` when
	// disabled or not set
	Forwarder *Forwarder `json:"forwarder"`
	// True when regular hosting is active so the URL forwarder cannot be managed via this API
	ForwarderBlockedByHosting bool `json:"forwarder_blocked_by_hosting"`
	// Whether the URL forwarder (DirectAdmin redirect) product is active for this domain
	ForwarderEnabled bool `json:"forwarder_enabled"`
	// Domain profile handles
	Handles Handles `json:"handles"`
	// Domain lock possible
	IsLockable bool `json:"is_lockable"`
	// Domain lock status
	IsLocked bool `json:"is_locked"`
	// DNS managed by mijn.host
	ManagedDNS bool `json:"managed_dns"`
	// Optional messages regarding the domain status
	Messages []string `json:"messages"`
	// Domain nameservers
	Nameservers []string `json:"nameservers"`
	// Renewal date (yyyy-mm-dd)
	RenewalDate string `json:"renewal_date"`
	// Domain status
	Status Status `json:"status"`
	// Domain tags (same format as `GET /domains/` list items)
	Tags []string `json:"tags"`
	// Whitelabel nameservers active
	WhitelabelNS bool `json:"whitelabel_ns"`
}

type DnssecKey struct {
	Alg    *string `json:"alg,omitempty"`
	Flags  *string `json:"flags,omitempty"`
	PubKey *string `json:"pubKey,omitempty"`
}

type Forwarder struct {
	// HTTP redirect status
	Type int64 `json:"type"`
	// Redirect target URL
	URL string `json:"url"`
}

// Domain profile handles
type Handles struct {
	Admin    Admin    `json:"admin"`
	Owner    Owner    `json:"owner"`
	Reseller Reseller `json:"reseller"`
	Tech     Tech     `json:"tech"`
}

type Admin struct {
	HandleID int64  `json:"handle_id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}

type Owner struct {
	HandleID int64  `json:"handle_id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}

type Reseller struct {
	HandleID int64  `json:"handle_id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}

type Tech struct {
	HandleID int64  `json:"handle_id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}

// Domain status
type Status string

const (
	Active                        Status = "active"
	ActiveOutgoingTransferPending Status = "Active, outgoing transfer pending"
	Cancelled                     Status = "Cancelled"
	Failed                        Status = "Failed"
	IncomingTransfer              Status = "Incoming transfer"
	NotFinishedOrder              Status = "Not finished order"
	PendingTransfer               Status = "pending_transfer"
)
