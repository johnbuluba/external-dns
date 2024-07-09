package editor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CoreDNSConfigEditorTestSuite struct {
	suite.Suite
	config     *CoreDNSConfigEditor
	testConfig string
}

func (s *CoreDNSConfigEditorTestSuite) SetupTest() {
	s.config = NewCoreDNSConfigEditor()
	s.testConfig = `
.:53 {
	errors
	health {
	   lameduck 5s
	}
    # Test comment
	ready
	kubernetes cluster.local in-addr.arpa ip6.arpa {
	   pods insecure
	   fallthrough in-addr.arpa ip6.arpa
	   ttl 30
	}
	prometheus :9153
	forward . /etc/resolv.conf {
	   max_concurrent 1000
	}
	cache 30
	loop
	reload
	loadbalance
}
`
}

func (s *CoreDNSConfigEditorTestSuite) TestNewCoreDNSConfig() {
	assert.NotNil(s.T(), s.config, "NewCoreDNSConfigEditor should return a non-nil CoreDNSConfigEditor instance")
}

func (s *CoreDNSConfigEditorTestSuite) TestLoadCorefile_Success() {
	err := s.config.LoadCorefile(s.testConfig)
	assert.Nil(s.T(), err, "LoadCorefile should not return an error for valid input")
	assert.NotNil(s.T(), s.config.serverBlocks, "LoadCorefile should populate the serverBlocks field on success")
}

func (s *CoreDNSConfigEditorTestSuite) TestLoadCorefile_Error() {
	config := "invalid {:}} invalid {"
	err := s.config.LoadCorefile(config)
	assert.NotNil(s.T(), err, "LoadCorefile should return an error for invalid input")
}

func (s *CoreDNSConfigEditorTestSuite) TestGetZones() {
	// Step 2: Invoke LoadCorefile
	err := s.config.LoadCorefile(`
.:53 {
    errors
	health {
	   lameduck 5s
	}
	file /etc/coredns/Zone test1.com test2.com test3.com
    file /unrelated/file google.com
    file /root
}
`)
	s.Require().Nil(err, "LoadCorefile should not return an error")

	// Step 3: Invoke AddFile
	zones := s.config.GetZones()

	// Step 4: Assert Conditions
	s.Equal(zones, []string{"test1.com", "test2.com", "test3.com"}, "GetZones should return the correct zones")
}

func (s *CoreDNSConfigEditorTestSuite) TestAddZone_NewZone() {
	// Step 1: Add a new zone that does not exist
	err := s.config.LoadCorefile(`
.:53 {
	errors
	health {
	   lameduck 5s
	}

    file /unrelated/file google.com
    file /root
}
`)
	s.Require().Nil(err, "LoadCorefile should not return an error")
	err = s.config.AddZone("newzone.com")
	s.Require().Nil(err, "AddZone should not return an error for a new zone")

	// Step 2: Assert the zone was added
	zones := s.config.GetZones()
	s.Contains(zones, "newzone.com", "Newly added zone should be present in the zones list")
}

func (s *CoreDNSConfigEditorTestSuite) TestAddZone_AdditionalZone() {
	// Step 1: Add a new zone that does not exist
	err := s.config.LoadCorefile(`
.:53 {
	errors
	health {
	   lameduck 5s
	}

    file /etc/coredns/Zone test.com
    file /unrelated/file google.com
    file /root
}
`)
	s.Require().Nil(err, "LoadCorefile should not return an error")
	err = s.config.AddZone("newzone.com")
	s.Require().Nil(err, "AddZone should not return an error for a new zone")

	// Step 2: Assert the zone was added
	zones := s.config.GetZones()
	s.Contains(zones, "newzone.com", "Newly added zone should be present in the zones list")
}

func (s *CoreDNSConfigEditorTestSuite) TestAddZone_ExistingZone() {
	// Step 1: Add an existing zone
	err := s.config.LoadCorefile(`
.:53 {
	errors
	health {
	   lameduck 5s
	}
    file /etc/coredns/Zone test.com
    file /unrelated/file google.com
    file /root
}
`)
	s.Require().Nil(err, "LoadCorefile should not return an error")
	err = s.config.AddZone("test.com")
	s.Require().Nil(err, "AddZone should not return an error for an existing zone")

	// Step 2: Assert the zone is not duplicated
	zones := s.config.GetZones()
	count := 0
	for _, zone := range zones {
		if zone == "test.com" {
			count++
		}
	}
	s.Equal(1, count, "Existing zone should not be duplicated")
}

func (s *CoreDNSConfigEditorTestSuite) TestRemoveZone_ExistingZone() {
	// Step 1: Remove an existing zone
	err := s.config.LoadCorefile(`
.:53 {
	errors
	health {
	   lameduck 5s
	}
    file /etc/coredns/Zone test.com test2.com
    file /unrelated/file google.com
    file /root
}
`)
	s.Require().Nil(err, "LoadCorefile should not return an error")
	err = s.config.RemoveZone("test.com")
	s.Require().Nil(err, "RemoveZone should not return an error for an existing zone")

	// Step 2: Assert the zone was removed
	zones := s.config.GetZones()
	s.NotContains(zones, "test.com", "Removed zone should not be present in the zones list")
}

func (s *CoreDNSConfigEditorTestSuite) TestRemoveZone_LastZone() {
	// Step 1: Remove an existing zone
	err := s.config.LoadCorefile(`
.:53 {
	errors
	health {
	   lameduck 5s
	}

    file /etc/coredns/Zone test.com
    file /unrelated/file google.com
    file /root
}
`)
	s.Require().Nil(err, "LoadCorefile should not return an error")
	err = s.config.RemoveZone("test.com")
	s.Require().Nil(err, "RemoveZone should not return an error for an existing zone")

	// Step 2: Assert the zone was removed
	zones := s.config.GetZones()
	s.NotContains(zones, "test.com", "Removed zone should not be present in the zones list")
	s.NotContains(s.config.GetConfig(), "file /etc/coredns/Zone", "file plugin should not be present in the config")
}

func (s *CoreDNSConfigEditorTestSuite) TestRemoveZone_NonExistingZone() {
	// Step 1: Attempt to remove a zone that does not exist
	err := s.config.LoadCorefile(`
.:53 {
	errors
	health {
	   lameduck 5s
	}
    file /etc/coredns/Zone test.com test2.com
    file /unrelated/file google.com
    file /root
}
`)
	s.Require().Nil(err, "LoadCorefile should not return an error")
	err = s.config.RemoveZone("nonexistingzone.com")
	s.Require().Nil(err, "RemoveZone should not return an error for a non-existing zone")

	// Step 2: Assert the zones list remains unchanged
	zonesBefore := s.config.GetZones()
	err = s.config.RemoveZone("nonexistingzone.com")
	s.Require().Nil(err, "RemoveZone should not affect zones list for a non-existing zone")
	zonesAfter := s.config.GetZones()
	s.Equal(zonesBefore, zonesAfter, "Zones list should remain unchanged after attempting to remove a non-existing zone")
}

func (s *CoreDNSConfigEditorTestSuite) TestGetConfig() {
	// Step 1: Invoke LoadCorefile
	err := s.config.LoadCorefile(s.testConfig)
	s.Require().Nil(err, "loadCorefile should not return an error")

	// Step 2: Invoke GetConfig
	configContent := s.config.GetConfig()

	// Step 3: Assert Conditions
	s.Equal(s.testConfig, configContent, "GetConfig should return the correct Corefile content")
}

func TestCoreDNSConfigEditorTestSuite(t *testing.T) {
	suite.Run(t, new(CoreDNSConfigEditorTestSuite))
}
