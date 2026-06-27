package mijnhost

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/DNSControl/dnscontrol/v4/models"
	"github.com/DNSControl/dnscontrol/v4/pkg/diff2"
	"github.com/DNSControl/dnscontrol/v4/pkg/providers"

	"github.com/DNSControl/dnscontrol/v4/providers/mijnhost/mijnhostapi"
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
	// mijn.host does not include NS records, strip them from the target (dc.Records)
	removeDomainNameserversFromDomainRecords(dc)

	changelist, actualChangeCount, err := diff2.ByRecord(existing, dc, nil)
	if err != nil {
		return nil, actualChangeCount, err
	}

	if len(changelist) == 0 {
		return []*models.Correction{}, 0, nil
	}

	corrections := make([]*models.Correction, 0, len(changelist))
	for _, change := range changelist {
		correction := change.CreateCorrection(mhp.constructChangeOperation(dc, change))
		if change.Type != diff2.REPORT && correction == nil {
			panic("constructChangeOperation did not provide a correction for non-REPORT change")
		}
		corrections = append(corrections, correction)
	}

	return corrections, actualChangeCount, nil
}

// removeDomainNameserversFromDomainRecords removes the nameserver records from the dc.Records which are already defined as the Domain nameservers.
//
// carefully copied from providers/transip
func removeDomainNameserversFromDomainRecords(dc *models.DomainConfig) {
	nameserverLookup := map[string]any{}
	for _, nameserver := range dc.Nameservers {
		nameserverLookup[nameserver.Name] = nil
	}

	newList := make([]*models.RecordConfig, 0, len(dc.Records))
	for _, rec := range dc.Records {

		dotLessNameFQDN := strings.TrimRight(rec.GetTargetField(), ".")
		_, recordInDCNameservers := nameserverLookup[dotLessNameFQDN]

		if rec.Type == "NS" && recordInDCNameservers {
			continue
		}

		newList = append(newList, rec)
	}
	dc.Records = newList
}

func (mhp *mijnhostProvider) constructChangeOperation(dc *models.DomainConfig, change diff2.Change) func() error {
	switch change.Type {
	case diff2.CREATE:
		if len(change.New) != 1 {
			panic("change CREATE did not contain exactly 1 New record")
		}
		if len(change.Old) != 0 {
			panic("change CREATE did not contain exactly 0 Old records")
		}
	case diff2.CHANGE:
		if len(change.New) != 1 {
			panic("change CHANGE did not contain exactly 1 New record")
		}
		if len(change.Old) != 1 {
			panic("change CHANGE did not contain exactly 1 Old record")
		}
	case diff2.DELETE:
		if len(change.New) != 0 {
			panic("change DELETE did not contain exactly 0 New records")
		}
		if len(change.Old) != 1 {
			panic("change DELETE did not contain exactly 1 Old record")
		}
	case diff2.REPORT:
		return nil
	default:
		panic("change contained unknown diff2.Type")
	}

	return func() error {
		var nativeRecord mijnhostapi.Record
		switch change.Type {
		case diff2.CHANGE:
			nativeRecord = recordConfigToNative(change.New[0])
			fallthrough
		case diff2.CREATE:
			_, err := mhp.client.UpdateDNSRecordForDomain(dc.Name, nativeRecord)
			if err != nil {
				return err
			}
		case diff2.DELETE:
			//nativeRecord = recordConfigToNative(change.Old[0])
		default:
			panic("unsupported change.Type")
		}

		return nil
	}
}
