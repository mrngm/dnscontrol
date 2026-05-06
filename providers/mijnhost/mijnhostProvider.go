package mijnhost

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/DNSControl/dnscontrol/v4/models"
	"github.com/DNSControl/dnscontrol/v4/pkg/providers"
)

/*
Mijn.host DNS provider

Info required in `creds.json`:
   - apikey
*/

var features = providers.DocumentationNotes{
	providers.CanAutoDNSSEC:          providers.Unimplemented("Not yet implemented"),
	providers.CanConcur:              providers.Unimplemented("Not yet implemented"),
	providers.CanGetZones:            providers.Unimplemented("Not yet implemented"),
	providers.CanOnlyDiff1Features:   providers.Unimplemented("Not yet implemented"),
	providers.CanUseAKAMAICDN:        providers.Unimplemented("Not yet implemented"),
	providers.CanUseAKAMAITLC:        providers.Unimplemented("Not yet implemented"),
	providers.CanUseAlias:            providers.Unimplemented("Not yet implemented"),
	providers.CanUseAzureAlias:       providers.Unimplemented("Not yet implemented"),
	providers.CanUseCAA:              providers.Unimplemented("Not yet implemented"),
	providers.CanUseDHCID:            providers.Unimplemented("Not yet implemented"),
	providers.CanUseDNAME:            providers.Unimplemented("Not yet implemented"),
	providers.CanUseDNSKEY:           providers.Unimplemented("Not yet implemented"),
	providers.CanUseDSForChildren:    providers.Unimplemented("Not yet implemented"),
	providers.CanUseDS:               providers.Unimplemented("Not yet implemented"),
	providers.CanUseHTTPS:            providers.Unimplemented("Not yet implemented"),
	providers.CanUseLOC:              providers.Unimplemented("Not yet implemented"),
	providers.CanUseNAPTR:            providers.Unimplemented("Not yet implemented"),
	providers.CanUseOPENPGPKEY:       providers.Unimplemented("Not yet implemented"),
	providers.CanUsePTR:              providers.Unimplemented("Not yet implemented"),
	providers.CanUseRoute53Alias:     providers.Unimplemented("Not yet implemented"),
	providers.CanUseRP:               providers.Unimplemented("Not yet implemented"),
	providers.CanUseSMIMEA:           providers.Unimplemented("Not yet implemented"),
	providers.CanUseSOA:              providers.Unimplemented("Not yet implemented"),
	providers.CanUseSRV:              providers.Unimplemented("Not yet implemented"),
	providers.CanUseSSHFP:            providers.Unimplemented("Not yet implemented"),
	providers.CanUseSVCB:             providers.Unimplemented("Not yet implemented"),
	providers.CanUseTLSA:             providers.Unimplemented("Not yet implemented"),
	providers.DocCreateDomains:       providers.Unimplemented("Not yet implemented"),
	providers.DocDualHost:            providers.Unimplemented("Not yet implemented"),
	providers.DocOfficiallySupported: providers.Unimplemented("Not yet implemented"),
}

func init() {
	const providerName = "MIJNHOST"
	const providerMaintainer = "@mrngm"
	fns := providers.DspFuncs{
		Initializer:   newMijnhost,
		RecordAuditor: AuditRecords,
	}
	providers.RegisterDomainServiceProviderType(providerName, fns, features)
	providers.RegisterMaintainer(providerName, providerMaintainer)

}

func newMijnhost(conf map[string]string, metadata json.RawMessage) (providers.DNSServiceProvider, error) {
	apikey, ok := conf["apikey"]
	if !ok {
		return nil, fmt.Errorf("required setting 'apikey' missing")
	}
	if apikey == "" {
		return nil, fmt.Errorf("required setting 'apikey' empty")
	}

	mhp := &mijnhostProvider{
		client: newAPIClient(apikey, 60*time.Second),
	}

	return mhp, nil
}

type mijnhostProvider struct {
	client *mijnhostAPIClient
}

func (mhp *mijnhostProvider) GetNameservers(domain string) ([]*models.Nameserver, error) {
	response, err := mhp.client.GetDomain(domain)
	if err != nil {
		return nil, fmt.Errorf("could not get domain information for %s: %w", domain, err)
	}

	// Sanity checks
	if response.Domain != domain {
		return nil, fmt.Errorf("requested domain differs from domain in response: %s != %s", response.Domain, domain)
	}

	return models.ToNameservers(response.Nameservers)
}

func (mhp *mijnhostProvider) GetZoneRecords(dc *models.DomainConfig) (models.Records, error) {
	response, err := mhp.client.GetDNSRecordsForDomain(dc.Name)
	if err != nil {
		return nil, fmt.Errorf("could not get zone records for %s: %w", dc.Name, err)
	}

	// Sanity checks
	if response.Domain != dc.Name {
		return nil, fmt.Errorf("requested domain differs from domain in response: %s != %s", response.Domain, dc.Name)
	}

	ret := make(models.Records, 0, len(response.Records))
	for _, record := range response.Records {
		ret = append(ret, nativeToRecordConfig(response.Domain, record))
	}

	return ret, nil
}

func (mhp *mijnhostProvider) GetZoneRecordsCorrections(dc *models.DomainConfig, existing models.Records) ([]*models.Correction, int, error) {
	return nil, 0, fmt.Errorf("unimplemented")
}
