package mijnhost

import (
	"fmt"
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

func recordConfigToNative(rc *models.RecordConfig) mijnhostapi.Record {
	rec := mijnhostapi.Record{
		Name: rc.GetLabel(), // may be overwritten below
		TTL:  int64(rc.TTL),
		Type: rc.Type,
		// Value is set below
	}

	switch rc.Type {
	case "A", "AAAA":
		rec.Value = rc.GetTargetIP().WithZone("").String()
	case "TXT":
		rec.Value = rc.GetTargetTXTJoined()
	case "MX":
		rec.Value = fmt.Sprintf("%d %s", rc.MxPreference, rc.GetTargetCombinedFunc(nil))
	default:
		panic(fmt.Sprintf("recordConfigToNative for type %q not implemented", rc.Type))
	}

	return rec
}

// convertTTL checks that the native value fits within the recordConfig constraints or panics otherwise
func convertTTL(native int64) uint32 {
	if err := ttlAllowed(native); err != nil {
		panic(fmt.Errorf("convertTTL: %w", err))
	}
	return uint32(native)
}

// ttlAllowed returns non-nil error if the given TTL is out of bounds [0, 2**32)
func ttlAllowed(native int64) error {
	if native < 0 {
		return fmt.Errorf("negative TTL not allowed")
	}
	if native > int64(MAX_UINT32) {
		panic("TTL out of bounds")
	}
	return nil
}
