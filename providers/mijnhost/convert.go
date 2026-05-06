package mijnhost

import (
	"net/netip"

	"github.com/DNSControl/dnscontrol/v4/models"
	"github.com/DNSControl/dnscontrol/v4/providers/mijnhost/mijnhostapi"
)

const MAX_UINT32 = ^uint32(0)

func nativeToRecordConfig(domain string, native mijnhostapi.Record) *models.RecordConfig {
	rc := &models.RecordConfig{
		TTL: convertTTL(native.TTL),
	}
	// native.Name contains the label, a dot, the domain name, a dot
	rc.SetLabelFromFQDN(native.Name, domain)

	switch native.Type {
	case "A":
		tgt, err := netip.ParseAddr(native.Value)
		if err != nil {
			panic("nativeToRecordConfig: native A value does not parse to a valid address")
		}
		if !tgt.Is4() {
			panic("nativeToRecordConfig: native A value is not an IPv4 address")
		}
		rc.Type = "A"
		rc.SetTargetIP(tgt)
	case "AAAA":
		tgt, err := netip.ParseAddr(native.Value)
		if err != nil {
			panic("nativeToRecordConfig: native AAAA value does not parse to a valid address")
		}
		if !tgt.Is6() {
			panic("nativeToRecordConfig: native AAAA value is not an IPv6 address")
		}
		rc.Type = "AAAA"
		rc.SetTargetIP(tgt)
	case "TXT":
		rc.SetTargetTXT(native.Value)
	case "MX":
		rc.SetTargetMXString(native.Value)
	default:
		panic("not implemented")
	}

	return rc
}

// convertTTL checks that the native value fits within the recordConfig constraints
func convertTTL(native int64) uint32 {
	if native < 0 {
		panic("convertTTL: native value negative")
	}
	if native > int64(MAX_UINT32) {
		panic("convertTTL: native value too large for uint32")
	}
	return uint32(native)
}
