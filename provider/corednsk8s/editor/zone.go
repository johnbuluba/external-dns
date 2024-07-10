package editor

import (
	"fmt"
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
	z.entriesByZone = make(map[string][]dns.RR)
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

// AddRecord adds a record to the zone
func (z *ZoneEditor) AddRecord(zone string, record dns.RR) {
	zoneRecords := z.GetOrCreateZone(zone)
	zoneRecords = append(zoneRecords, record)
	z.entriesByZone[zone] = zoneRecords
	z.increaseZoneSerial(zone)
}

// UpdateRecord updates a record in the zone
func (z *ZoneEditor) UpdateRecord(zone string, oldRecord, newRecord dns.RR) {
	zoneRecords := z.GetOrCreateZone(zone)
	newRecords := make([]dns.RR, 0)
	for _, r := range zoneRecords {
		if r.String() == oldRecord.String() {
			newRecords = append(newRecords, newRecord)
		} else {
			newRecords = append(newRecords, r)
		}
	}
	z.entriesByZone[zone] = newRecords
	z.increaseZoneSerial(zone)
}

// DeleteRecord deletes a record from the zone
func (z *ZoneEditor) DeleteRecord(zone string, record dns.RR) {
	zoneRecords := z.GetOrCreateZone(zone)
	newRecords := make([]dns.RR, 0)
	for _, r := range zoneRecords {
		if !(r.Header().Name == record.Header().Name && r.Header().Rrtype == record.Header().Rrtype) {
			newRecords = append(newRecords, r)
		}
	}
	// If only the SOA record is left, remove the zone
	if len(newRecords) == 1 {
		delete(z.entriesByZone, zone)
		return
	}
	z.entriesByZone[zone] = newRecords
	z.increaseZoneSerial(zone)
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

// increaseZoneSerial increases the serial number of the SOA record for the zone
//
// This is required to notify the CoreDNS server that the zone has been updated.
func (z *ZoneEditor) increaseZoneSerial(zone string) {
	zoneRecords := z.GetOrCreateZone(zone)
	for _, record := range zoneRecords {
		if soa, ok := record.(*dns.SOA); ok {
			soa.Serial++
			z.entriesByZone[zone] = zoneRecords
			return
		}
	}
}
