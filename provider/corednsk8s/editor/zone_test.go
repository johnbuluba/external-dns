package editor

import (
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

func (s *ZoneEditorTestSuite) TestAddARecord() {
	zone := "example.com."
	name := "test.example.com."
	ip := "192.0.2.1"

	s.editor.AddARecord(zone, name, ip)

	records, found := s.editor.entriesByZone[zone]
	s.Require().True(found, "Zone should exist after adding an A record")

	foundRecord := false
	for _, record := range records {
		if aRecord, ok := record.(*dns.A); ok {
			if aRecord.Hdr.Name == name && aRecord.A.String() == ip {
				foundRecord = true
				break
			}
		}
	}

	s.Require().True(foundRecord, "A record should be found in the zone")
}

func (s *ZoneEditorTestSuite) TestAddARecord_SecondTime() {
	zone := "example.com."
	name := "test.example.com."
	ip := "192.0.2.1"

	// Add the first time
	s.editor.AddARecord(zone, name, ip)
	// Add a second time with different IP
	ip = "192.0.2.2"
	s.editor.AddARecord(zone, name, ip)

	records, found := s.editor.entriesByZone[zone]
	s.Require().True(found, "Zone should exist after adding an A record")

	foundRecord := false
	for _, record := range records {
		if aRecord, ok := record.(*dns.A); ok {
			if aRecord.Hdr.Name == name && aRecord.A.String() == ip {
				foundRecord = true
				break
			}
		}
	}

	s.Require().True(foundRecord, "A record should be found in the zone")
}

func TestZoneEditorTestSuite(t *testing.T) {
	suite.Run(t, new(ZoneEditorTestSuite))
}
