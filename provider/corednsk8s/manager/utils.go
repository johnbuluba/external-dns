package manager

import (
	"net"
	"strings"

	"github.com/miekg/dns"
	"sigs.k8s.io/external-dns/endpoint"
)

// EndpointToRecord converts an endpoint to a DNS record.
func endpointToRecord(endpt *endpoint.Endpoint) dns.RR {
	switch endpt.RecordType {
	case endpoint.RecordTypeA:
		return &dns.A{
			Hdr: dns.RR_Header{
				Name:   dns.Fqdn(endpt.DNSName),
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    uint32(endpt.RecordTTL),
			},
			A: net.ParseIP(endpt.Targets[0]),
		}
	case endpoint.RecordTypeCNAME:
		return &dns.CNAME{
			Hdr: dns.RR_Header{
				Name:   dns.Fqdn(endpt.DNSName),
				Rrtype: dns.TypeCNAME,
				Class:  dns.ClassINET,
				Ttl:    uint32(endpt.RecordTTL),
			},
			Target: endpt.Targets[0],
		}
	case endpoint.RecordTypeTXT:
		// Remove quotes from targets
		trimmedTargets := make([]string, len(endpt.Targets))
		for i, target := range endpt.Targets {
			trimmedTargets[i] = strings.Trim(target, "\"")
		}
		return &dns.TXT{
			Hdr: dns.RR_Header{
				Name:   dns.Fqdn(endpt.DNSName),
				Rrtype: dns.TypeTXT,
				Class:  dns.ClassINET,
				Ttl:    uint32(endpt.RecordTTL),
			},
			Txt: trimmedTargets,
		}
	case endpoint.RecordTypeSRV:
		return &dns.SRV{
			Hdr: dns.RR_Header{
				Name:   dns.Fqdn(endpt.DNSName),
				Rrtype: dns.TypeSRV,
				Class:  dns.ClassINET,
				Ttl:    uint32(endpt.RecordTTL),
			},
			Target: endpt.Targets[0],
		}
	case endpoint.RecordTypePTR:
		return &dns.PTR{
			Hdr: dns.RR_Header{
				Name:   dns.Fqdn(endpt.DNSName),
				Rrtype: dns.TypePTR,
				Class:  dns.ClassINET,
				Ttl:    uint32(endpt.RecordTTL),
			},
			Ptr: endpt.Targets[0],
		}
	case endpoint.RecordTypeMX:
		return &dns.MX{
			Hdr: dns.RR_Header{
				Name:   dns.Fqdn(endpt.DNSName),
				Rrtype: dns.TypeMX,
				Class:  dns.ClassINET,
				Ttl:    uint32(endpt.RecordTTL),
			},
			Mx: endpt.Targets[0],
		}
	case endpoint.RecordTypeAAAA:
		return &dns.AAAA{
			Hdr: dns.RR_Header{
				Name:   dns.Fqdn(endpt.DNSName),
				Rrtype: dns.TypeAAAA,
				Class:  dns.ClassINET,
				Ttl:    uint32(endpt.RecordTTL),
			},
			AAAA: net.ParseIP(endpt.Targets[0]),
		}
	case endpoint.RecordTypeNS:
		return &dns.NS{
			Hdr: dns.RR_Header{
				Name:   dns.Fqdn(endpt.DNSName),
				Rrtype: dns.TypeNS,
				Class:  dns.ClassINET,
				Ttl:    uint32(endpt.RecordTTL),
			},
			Ns: endpt.Targets[0],
		}
	}
	return nil
}
