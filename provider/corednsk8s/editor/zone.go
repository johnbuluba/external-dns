package editor

import (
	"fmt"
	"net"
	"strings"

	"github.com/miekg/dns"
)

type ZoneEditor struct {
	entriesByZone map[string][]dns.RR
}

// NewZoneEditor creates a new ZoneEditor
func NewZoneEditor() *ZoneEditor {
	return &ZoneEditor{
		entriesByZone: make(map[string][]dns.RR),
	}
}

// LoadZone loads a zone file
func (z *ZoneEditor) LoadZone(zoneFile string) error {
	parser := dns.NewZoneParser(strings.NewReader(zoneFile), "", "")
	zone := ""
	for rr, ok := parser.Next(); ok; rr, ok = parser.Next() {
		// If the record is an SOA record, update the current zone
		if soa, ok := rr.(*dns.SOA); ok {
			zone = soa.Hdr.Name
		} else if zone == "" {
			return fmt.Errorf("SOA record not found in zone file")
		}
		// Get the records for the current zone
		zoneRecords, found := z.entriesByZone[zone]
		if !found {
			// If the zone doesn't exist, create it
			zoneRecords = make([]dns.RR, 0)
		}
		// Append the record to the zone
		zoneRecords = append(zoneRecords, rr)
		z.entriesByZone[zone] = zoneRecords
	}
	return parser.Err()
}

// AddARecord adds an A record to the zone
func (z *ZoneEditor) AddARecord(zone, name, ip string) {
	zoneRecords := z.GetOrCreateZone(zone)
	// Check if record already exists
	if r := z.GetRecordByTypeAndName(zone, name, dns.TypeA); r != nil {
		// Check if the IP is the same
		if aRecord, ok := r.(*dns.A); ok {
			if aRecord.A.String() == ip {
				return
			}
		}
		// If the IP is different, remove the record
		z.RemoveRecord(zone, name, dns.TypeA)
	}

	zoneRecords = append(zoneRecords, &dns.A{
		Hdr: dns.RR_Header{Name: name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
		A:   net.ParseIP(ip),
	})
	z.entriesByZone[zone] = zoneRecords
}

// GetZones returns all zones
func (z *ZoneEditor) GetZones() []string {
	zones := make([]string, 0)
	for zone := range z.entriesByZone {
		zones = append(zones, zone)
	}
	return zones
}

// GetAllRecords returns all records for a zone
func (z *ZoneEditor) GetAllRecords(zone string) []dns.RR {
	records, found := z.entriesByZone[zone]
	if !found {
		return nil
	}
	return records
}

// GetRecordByTypeAndName returns a record by type and name
func (z *ZoneEditor) GetRecordByTypeAndName(zone, name string, recordType uint16) dns.RR {
	records, found := z.entriesByZone[zone]
	if !found {
		return nil
	}
	for _, record := range records {
		if record.Header().Name == name && record.Header().Rrtype == recordType {
			return record
		}
	}
	return nil
}

// RenderZone prints the zone file
func (z *ZoneEditor) RenderZone() string {
	var sb strings.Builder

	for zone, records := range z.entriesByZone {
		sb.WriteString(fmt.Sprintf("$ORIGIN %s\n", zone))
		for _, record := range records {
			sb.WriteString(record.String())
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// RemoveRecord removes a record from the zone
func (z *ZoneEditor) RemoveRecord(zone, name string, recordType uint16) {
	zoneRecords := z.GetOrCreateZone(zone)
	newRecords := make([]dns.RR, 0)
	for _, record := range zoneRecords {
		if record.Header().Name == name && record.Header().Rrtype == recordType {
			continue
		}
		newRecords = append(newRecords, record)
	}
	z.entriesByZone[zone] = newRecords
}

// GetOrCreateZone returns the records for a zone, creating it if it doesn't exist
func (z *ZoneEditor) GetOrCreateZone(zone string) []dns.RR {
	records, found := z.entriesByZone[zone]
	if !found {
		records = make([]dns.RR, 0)
		// TODO: Make SOA configurable
		soa := dns.SOA{
			Hdr: dns.RR_Header{
				Name:     zone,
				Rrtype:   dns.TypeSOA,
				Class:    dns.ClassINET,
				Ttl:      3600,
				Rdlength: 0,
			},
			Ns:      fmt.Sprintf("ns.%s", zone),
			Mbox:    fmt.Sprintf("hostmaster.%s", zone),
			Serial:  1,
			Refresh: 3600,
			Retry:   3600,
			Expire:  3600,
			Minttl:  3600,
		}
		records = append(records, &soa)
		z.entriesByZone[zone] = records
	}
	return records
}
