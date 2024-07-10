package editor

import (
	"net"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/suite"
)

type ZoneEditorTestSuite struct {
	suite.Suite
	editor *ZoneEditor
}

func (s *ZoneEditorTestSuite) SetupTest() {
	s.editor = NewZoneEditor()
}

func (s *ZoneEditorTestSuite) TestLoadZone() {
	zoneFile := `$ORIGIN example.com.
@ 3600 IN SOA ns.example.com. hostmaster.example.com. (
              2023041501 ; serial
              7200       ; refresh (2 hours)
              3600       ; retry (1 hour)
              1209600    ; expire (2 weeks)
              3600       ; minimum (1 hour)
              )
    3600 IN NS ns.example.com.
ns  3600 IN A  192.0.2.1

$ORIGIN example.net.
@ 3600 IN SOA ns.example.net. hostmaster.example.net. (
              2023041501 ; serial
              7200       ; refresh (2 hours)
              3600       ; retry (1 hour)
              1209600    ; expire (2 weeks)
              3600       ; minimum (1 hour)
              )
    3600 IN NS ns.example.net.
ns  3600 IN A  192.0.2.1

`

	err := s.editor.LoadZone(zoneFile)
	s.Require().NoError(err, "LoadZone should not return an error")
	s.Require().Contains(s.editor.entriesByZone, "example.com.", "entriesByZone should contain 'example.com.'")
}

func (s *ZoneEditorTestSuite) TestLoadZone_InvalidZone() {
	zoneFile := `This is not a valid zone file.`

	err := s.editor.LoadZone(zoneFile)
	s.Require().Error(err, "LoadZone should return an error for invalid zone file")
}

func (s *ZoneEditorTestSuite) TestAddRecord_NewRecord() {
	zone := "example.com."
	record := &dns.A{
		Hdr: dns.RR_Header{Name: zone, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
		A:   net.ParseIP("192.0.2.1"),
	}

	s.editor.AddRecord(zone, record)
	records := s.editor.GetAllRecords(zone)

	s.Require().Contains(records, record, "Newly added record should be present in the zone")
}

func (s *ZoneEditorTestSuite) TestUpdateRecord_ExistingRecord() {
	zone := "example.com."
	oldRecord := &dns.A{
		Hdr: dns.RR_Header{Name: "ns." + zone, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
		A:   net.ParseIP("192.0.2.1"),
	}
	newRecord := &dns.A{
		Hdr: dns.RR_Header{Name: "ns." + zone, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
		A:   net.ParseIP("192.0.2.2"),
	}

	// Assuming oldRecord is already added
	s.editor.AddRecord(zone, oldRecord)
	s.editor.UpdateRecord(zone, oldRecord, newRecord)

	updatedRecords := s.editor.GetAllRecords(zone)
	s.Require().Contains(updatedRecords, newRecord, "Updated record should be present in the zone")
	s.Require().NotContains(updatedRecords, oldRecord, "Old record should not be present in the zone after update")
}

func (s *ZoneEditorTestSuite) TestDeleteRecord() {
	zone := "example.com."
	otherRecord := &dns.A{
		Hdr: dns.RR_Header{Name: "www." + zone, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
		A:   net.ParseIP("192.0.2.5"),
	}
	s.editor.AddRecord(zone, otherRecord)
	recordToDelete := &dns.A{
		Hdr: dns.RR_Header{Name: "ns." + zone, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
		A:   net.ParseIP("192.0.2.1"),
	}

	// Assuming recordToDelete is already added
	s.editor.AddRecord(zone, recordToDelete)
	s.editor.DeleteRecord(zone, recordToDelete)

	recordsAfterDeletion := s.editor.GetAllRecords(zone)
	s.Require().NotContains(recordsAfterDeletion, recordToDelete, "Deleted record should not be present in the zone")
}

func (s *ZoneEditorTestSuite) TestDeleteRecord_LastRecord() {
	zone := "example.com."
	recordToDelete := &dns.A{
		Hdr: dns.RR_Header{Name: "ns." + zone, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
		A:   net.ParseIP("192.0.2.1"),
	}

	// Assuming recordToDelete is already added
	s.editor.AddRecord(zone, recordToDelete)
	s.editor.DeleteRecord(zone, recordToDelete)

	recordsAfterDeletion := s.editor.GetAllRecords(zone)
	s.Require().NotContains(recordsAfterDeletion, recordToDelete, "Deleted record should not be present in the zone")
	s.Require().NotContains(s.editor.GetZones(), zone, "The zone should be removed if it has no records left")
}

func (s *ZoneEditorTestSuite) TestRenderZone() {
	zoneFile := `$ORIGIN example.com.
@ 3600 IN SOA ns.example.com. hostmaster.example.com. (
              2023041501 ; serial
              7200       ; refresh (2 hours)
              3600       ; retry (1 hour)
              1209600    ; expire (2 weeks)
              3600       ; minimum (1 hour)
              )
    3600 IN NS ns.example.com.
ns  3600 IN A  192.0.2.1`

	expectedRender := `$ORIGIN example.com.
example.com.	3600	IN	SOA	ns.example.com. hostmaster.example.com. 2023041501 7200 3600 1209600 3600
example.com.	3600	IN	NS	ns.example.com.
ns.example.com.	3600	IN	A	192.0.2.1
`

	err := s.editor.LoadZone(zoneFile)
	s.Require().NoError(err, "LoadZone should not return an error")

	renderedZone := s.editor.RenderZone()
	s.Require().Contains(expectedRender, renderedZone, "RenderZone output should match the expected render")
}

func TestZoneEditorTestSuite(t *testing.T) {
	suite.Run(t, new(ZoneEditorTestSuite))
}
