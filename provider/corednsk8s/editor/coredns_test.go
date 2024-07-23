package editor

import (
	"fmt"
	"net"
	"testing"

	"github.com/miekg/dns"
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
	s.Require().NoError(s.config.LoadCorefile(s.testConfig))
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

func (s *CoreDNSConfigEditorTestSuite) TestSetZones_MultipleZones() {
	zones := []string{"example.com", "test.com"}
	err := s.config.SetZones(zones)
	s.Require().Nil(err, "SetZones should not return an error for valid zones")

	// Assert the zones were added
	actualZones := s.config.GetZones()
	for _, zone := range zones {
		s.Contains(actualZones, zone, fmt.Sprintf("Zone %s should be present in the zones list", zone))
	}
}

func (s *CoreDNSConfigEditorTestSuite) TestSetZones_EmptyZones() {
	err := s.config.SetZones([]string{})
	s.Require().Nil(err, "SetZones should not return an error for an empty zones list")

	// Assert no zones are present
	actualZones := s.config.GetZones()
	s.Len(actualZones, 0, "No zones should be present in the zones list")
}

func (s *CoreDNSConfigEditorTestSuite) TestSetZones_RemoveExistingZones() {
	// Initially set some zones
	initialZones := []string{"example.com", "test.com"}
	s.config.SetZones(initialZones)

	// Now remove all zones
	err := s.config.SetZones([]string{})
	s.Require().Nil(err, "SetZones should not return an error when removing all zones")

	// Assert no zones are present
	actualZones := s.config.GetZones()
	s.Len(actualZones, 0, "No zones should be present after removing them")
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

func (s *CoreDNSConfigEditorTestSuite) TestSetHosts_MultipleHosts() {
	err := s.config.LoadCorefile(`
.:53 {
	errors
	health {
	   lameduck 5s
	}

    file /etc/coredns/Zone test.com
    file /unrelated/file google.com
    hosts {
       192.168.10.1 test.com
    }
    file /root
}
`)
	s.Require().Nil(err, "LoadCorefile should not return an error")
	hosts := []dns.A{
		{
			Hdr: dns.RR_Header{Name: "example.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
			A:   net.ParseIP("192.0.2.1"),
		},
		{
			Hdr: dns.RR_Header{Name: "test.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
			A:   net.ParseIP("192.0.2.2"),
		},
	}

	err = s.config.SetHosts(hosts)
	s.Require().Nil(err, "SetHosts should not return an error for valid hosts")

	// Assert the hosts were added
	configContent := s.config.GetConfig()
	for _, host := range hosts {
		s.Contains(configContent, host.A.String(), fmt.Sprintf("Host %s should be present in the config", host.A.String()))
		s.Contains(configContent, host.Hdr.Name, fmt.Sprintf("Host name %s should be present in the config", host.Hdr.Name))
	}
}

func (s *CoreDNSConfigEditorTestSuite) TestSetHosts_EmptyHosts() {
	err := s.config.SetHosts([]dns.A{})
	s.Require().Nil(err, "SetHosts should not return an error for an empty hosts list")

	// Assert no hosts are present
	configContent := s.config.GetConfig()
	s.NotContains(configContent, "hosts {", "No hosts block should be present in the config")
}

func (s *CoreDNSConfigEditorTestSuite) TestGetHosts_WithHosts() {
	// Load a Corefile with predefined hosts
	err := s.config.LoadCorefile(`
.:53 {
    errors
    health {
       lameduck 5s
    }
    hosts {
       192.168.10.1 example.com
       192.168.10.2 test.com
    }
    file /root
}
`)
	s.Require().Nil(err, "LoadCorefile should not return an error")

	// Call GetHosts
	hosts := s.config.GetHosts()

	// Assert the hosts are correctly retrieved
	expectedHosts := []dns.A{
		{
			Hdr: dns.RR_Header{Name: "example.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
			A:   net.ParseIP("192.168.10.1"),
		},
		{
			Hdr: dns.RR_Header{Name: "test.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
			A:   net.ParseIP("192.168.10.2"),
		},
	}

	for _, expectedHost := range expectedHosts {
		found := false
		for _, host := range hosts {
			if host.Hdr.Name == expectedHost.Hdr.Name && host.A.String() == expectedHost.A.String() {
				found = true
				break
			}
		}
		s.True(found, fmt.Sprintf("Host %s should be present in the hosts", expectedHost.Hdr.Name))
	}
}

func (s *CoreDNSConfigEditorTestSuite) TestGetHosts_NoHosts() {
	// Load a Corefile without hosts
	err := s.config.LoadCorefile(`
.:53 {
    errors
    health {
       lameduck 5s
    }
    file /root
}
`)
	s.Require().Nil(err, "LoadCorefile should not return an error")

	// Call GetHosts
	hosts := s.config.GetHosts()

	// Assert no hosts are returned
	s.Len(hosts, 0, "No hosts should be present")
}

func (s *CoreDNSConfigEditorTestSuite) TestAddHost_NewHost() {
	host := dns.A{
		Hdr: dns.RR_Header{Name: "newhost.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
		A:   net.ParseIP("192.0.2.3"),
	}

	// Step 1: Add a new host
	err := s.config.AddHost(host)
	s.Require().Nil(err, "AddHost should not return an error for a new host")

	// Step 2: Assert the host was added
	hosts := s.config.GetHosts()
	found := false
	for _, h := range hosts {
		if h.Hdr.Name == host.Hdr.Name && h.A.Equal(host.A) {
			found = true
			break
		}
	}
	s.True(found, "Newly added host should be present in the hosts list")
}

func (s *CoreDNSConfigEditorTestSuite) TestAddHost_ExistingHost() {
	host := dns.A{
		Hdr: dns.RR_Header{Name: "existinghost.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
		A:   net.ParseIP("192.0.2.4"),
	}

	// Assuming host is already added
	s.config.AddHost(host)

	// Step 1: Add the same host again
	err := s.config.AddHost(host)
	s.Require().Nil(err, "AddHost should not return an error for an existing host")

	// Step 2: Assert the host is not duplicated
	hosts := s.config.GetHosts()
	count := 0
	for _, h := range hosts {
		if h.Hdr.Name == host.Hdr.Name && h.A.Equal(host.A) {
			count++
		}
	}
	s.Equal(1, count, "Existing host should not be duplicated")
}

func (s *CoreDNSConfigEditorTestSuite) TestUpdateHost_ExistingHost() {
	oldHost := dns.A{
		Hdr: dns.RR_Header{Name: "oldhost.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
		A:   net.ParseIP("192.0.2.3"),
	}
	newHost := dns.A{
		Hdr: dns.RR_Header{Name: "newhost.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
		A:   net.ParseIP("192.0.2.4"),
	}

	// Assuming oldHost is already added
	s.config.AddHost(oldHost)
	err := s.config.UpdateHost(oldHost, newHost)
	s.Require().Nil(err, "UpdateHost should not return an error for updating an existing host")

	// Assert the host was updated
	hosts := s.config.GetHosts()
	foundNewHost := false
	foundOldHost := false
	for _, h := range hosts {
		if h.Hdr.Name == newHost.Hdr.Name && h.A.Equal(newHost.A) {
			foundNewHost = true
		}
		if h.Hdr.Name == oldHost.Hdr.Name && h.A.Equal(oldHost.A) {
			foundOldHost = true
		}
	}
	s.True(foundNewHost, "Updated host should be present in the hosts list")
	s.False(foundOldHost, "Old host should not be present in the hosts list after update")
}

func (s *CoreDNSConfigEditorTestSuite) TestUpdateHost_NonExistingHost() {
	nonExistingHost := dns.A{
		Hdr: dns.RR_Header{Name: "nonexistinghost.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
		A:   net.ParseIP("192.0.2.5"),
	}
	updatedHost := dns.A{
		Hdr: dns.RR_Header{Name: "updatedhost.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
		A:   net.ParseIP("192.0.2.6"),
	}

	// Attempt to update a non-existing host
	err := s.config.UpdateHost(nonExistingHost, updatedHost)
	s.Require().Nil(err, "UpdateHost should not return an error even if the host does not exist")

	// Assert the updated host is not added
	hosts := s.config.GetHosts()
	foundUpdatedHost := false
	for _, h := range hosts {
		if h.Hdr.Name == updatedHost.Hdr.Name && h.A.Equal(updatedHost.A) {
			foundUpdatedHost = true
			break
		}
	}
	s.False(foundUpdatedHost, "Non-existing host should not be added to the hosts list after update attempt")
}

func (s *CoreDNSConfigEditorTestSuite) TestRemoveHost_ExistingHost() {
	host := dns.A{
		Hdr: dns.RR_Header{Name: "removehost.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
		A:   net.ParseIP("192.0.2.5"),
	}

	// Assuming host is already added
	s.config.AddHost(host)

	// Step 1: Remove the existing host
	err := s.config.RemoveHost(host)
	s.Require().Nil(err, "RemoveHost should not return an error for an existing host")

	// Step 2: Assert the host was removed
	hosts := s.config.GetHosts()
	found := false
	for _, h := range hosts {
		if h.Hdr.Name == host.Hdr.Name && h.A.Equal(host.A) {
			found = true
			break
		}
	}
	s.False(found, "Removed host should not be present in the hosts list")
}

func (s *CoreDNSConfigEditorTestSuite) TestRemoveHost_NonExistingHost() {
	host := dns.A{
		Hdr: dns.RR_Header{Name: "nonexistinghost.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
		A:   net.ParseIP("192.0.2.6"),
	}

	// Step 1: Attempt to remove a host that does not exist
	err := s.config.RemoveHost(host)
	s.Require().Nil(err, "RemoveHost should not return an error for a non-existing host")

	// Step 2: Assert the hosts list remains unchanged
	hostsBefore := s.config.GetHosts()
	err = s.config.RemoveHost(host)
	s.Require().Nil(err, "RemoveHost should not affect hosts list for a non-existing host")
	hostsAfter := s.config.GetHosts()
	s.Equal(hostsBefore, hostsAfter, "Hosts list should remain unchanged after attempting to remove a non-existing host")
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
