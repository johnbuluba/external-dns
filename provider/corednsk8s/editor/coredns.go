package editor

import (
	"fmt"
	"math"
	"net"
	"strings"

	"github.com/coredns/caddy/caddyfile"
	"github.com/miekg/dns"
)

const ZoneFilePath = "/etc/coredns/Zone"

type CoreDNSConfigEditor struct {
	serverBlocks []caddyfile.ServerBlock
	text         string
}

// NewCoreDNSConfigEditor Creates a new CoreDNSConfigEditor
func NewCoreDNSConfigEditor() *CoreDNSConfigEditor {
	return &CoreDNSConfigEditor{}
}

// SetZones sets the zones in the file plugin in the Corefile
func (c *CoreDNSConfigEditor) SetZones(zones []string) error {
	// Remove the file entry
	err := c.removeFileEntry()
	if err != nil {
		return err
	}
	// If there are no zones, return
	if len(zones) == 0 {
		return nil
	}
	// Cleanup the zone names
	cleanUpZones := make([]string, len(zones))
	for i, z := range zones {
		cleanUpZones[i] = removeTrailingDot(z)
	}
	// Update the file configuration
	appendToLine := c.getLineToAppend()
	entry := fmt.Sprintf("    file %s %s", ZoneFilePath, strings.Join(cleanUpZones, " "))
	return c.addInLine(appendToLine, entry)
}

// AddZone adds a new zone in the file plugin
func (c *CoreDNSConfigEditor) AddZone(zone string) error {
	cleanZoneName := removeTrailingDot(zone)
	// Get existing zones
	zones := c.GetZones()
	// Check if the zone already exists
	for _, z := range zones {
		if z == cleanZoneName {
			return nil
		}
	}
	zones = append(zones, cleanZoneName)
	return c.SetZones(zones)
}

// RemoveZone removes a zone from the file plugin
func (c *CoreDNSConfigEditor) RemoveZone(zone string) error {
	cleanZoneName := removeTrailingDot(zone)
	// Get existing zones
	zones := c.GetZones()
	// Check if the zone exists
	found := false
	for i, z := range zones {
		if z == cleanZoneName {
			found = true
			zones = append(zones[:i], zones[i+1:]...)
			break
		}
	}
	if !found {
		return nil
	}
	// Remove the file entry
	err := c.removeFileEntry()
	if err != nil {
		return err
	}
	// If there are no more zones, remove the file entry
	if len(zones) == 0 {
		return nil
	}
	// Update the file configuration
	appendToLine := c.getLineToAppend()
	entry := fmt.Sprintf("    file %s %s", ZoneFilePath, strings.Join(zones, " "))
	return c.addInLine(appendToLine, entry)
}

// GetZones returns the zones from the file plugin in the Corefile
func (c *CoreDNSConfigEditor) GetZones() []string {
	block := c.get53Block()
	zones := make([]string, 0)
	if tokens, ok := block.Tokens["file"]; ok {
		ourFileEntryFound := false
		for _, line := range tokens {
			// We only want the entries in our file plugin
			if strings.Contains(line.Text, ZoneFilePath) {
				ourFileEntryFound = true
				continue
			} else if !ourFileEntryFound {
				continue
			}
			// If a new file is found means that a new file plugin is configured, break
			if strings.Contains(line.Text, "file") {
				break
			}
			zones = append(zones, strings.TrimSpace(line.Text))
		}
	}
	return zones
}

// SetHosts sets the hosts in the Corefile
func (c *CoreDNSConfigEditor) SetHosts(hosts []dns.A) error {
	// Remove the hosts
	err := c.removeHosts()
	if err != nil {
		return err
	}
	// If there are no hosts, return
	if len(hosts) == 0 {
		return nil
	}
	// Update the hosts configuration
	sb := strings.Builder{}
	sb.WriteString("  hosts {\n")
	for _, h := range hosts {
		sb.WriteString(fmt.Sprintf("    %s %s\n", h.A.String(), removeTrailingDot(h.Hdr.Name)))
	}
	sb.WriteString("    fallthrough\n")
	sb.WriteString("  }")

	appendToLine := c.getLineToAppendHosts()
	return c.addInLine(appendToLine, sb.String())
}

// AddHost adds a new host in the hosts plugin
func (c *CoreDNSConfigEditor) AddHost(host dns.A) error {
	// Get existing hosts
	hosts := c.GetHosts()
	// Check if the host already exists
	for _, h := range hosts {
		if h.A.Equal(host.A) && h.Hdr.Name == host.Hdr.Name {
			return nil
		}
	}
	hosts = append(hosts, host)
	return c.SetHosts(hosts)
}

// UpdateHost updates a host in the hosts plugin
func (c *CoreDNSConfigEditor) UpdateHost(oldHost, newHost dns.A) error {
	// Get existing hosts
	hosts := c.GetHosts()
	// Update the host
	newHosts := make([]dns.A, 0)
	for _, h := range hosts {
		if h.A.Equal(oldHost.A) && h.Hdr.Name == oldHost.Hdr.Name {
			newHosts = append(newHosts, newHost)
		} else {
			newHosts = append(newHosts, h)
		}
	}
	return c.SetHosts(newHosts)
}

// RemoveHost removes a host from the hosts plugin
func (c *CoreDNSConfigEditor) RemoveHost(host dns.A) error {
	// Get existing hosts
	hosts := c.GetHosts()
	// Remove the host
	newHosts := make([]dns.A, 0)
	for _, h := range hosts {
		if h.A.Equal(host.A) && h.Hdr.Name == host.Hdr.Name {
			continue
		}
		newHosts = append(newHosts, h)
	}
	return c.SetHosts(newHosts)
}

// GetHosts returns the hosts from the hosts plugin in the Corefile
func (c *CoreDNSConfigEditor) GetHosts() []dns.A {
	block := c.get53Block()
	hosts := make([]dns.A, 0)
	if tokens, ok := block.Tokens["hosts"]; ok {
		ourHostsEntryFound := false
		for i := 0; i < len(tokens); i += 1 {
			line := tokens[i]
			// We only want the entries in our hosts plugin
			if strings.Contains(line.Text, "hosts") {
				ourHostsEntryFound = true
				continue
			} else if strings.Contains(line.Text, "{") {
				continue
			} else if strings.Contains(line.Text, "fallthrough") {
				continue
			} else if !ourHostsEntryFound {
				continue
			}
			// If a new hosts is found means that a new hosts plugin is configured, break
			if strings.Contains(line.Text, "}") {
				break
			}
			// Parse the line
			ip := tokens[i].Text
			i += 1
			name := tokens[i].Text
			hosts = append(hosts, dns.A{Hdr: dns.RR_Header{Name: addTrailingDot(name)}, A: net.ParseIP(ip)})
		}
	}
	return hosts
}

// GetConfig returns the Corefile as a string
func (c *CoreDNSConfigEditor) GetConfig() string {
	return c.text
}

// getMaxLine returns the maximum line number in the Corefile
func (c *CoreDNSConfigEditor) getMaxLine() int {
	maxLine := 0
	block := c.get53Block()
	for _, token := range block.Tokens {
		for _, line := range token {
			if line.Line > maxLine {
				maxLine = line.Line
			}
		}
	}
	return maxLine
}

// getMinLine returns the minimum line number in the Corefile
func (c *CoreDNSConfigEditor) getMinLine() int {
	minLine := math.MaxInt
	block := c.get53Block()
	for _, token := range block.Tokens {
		for _, line := range token {
			if line.Line < minLine {
				minLine = line.Line
			}
		}
	}
	return minLine
}

// get53Block returns the .:53 block from the Corefile
func (c *CoreDNSConfigEditor) get53Block() *caddyfile.ServerBlock {
	for _, block := range c.serverBlocks {
		if block.Keys[0] == ".:53" {
			return &block
		}
	}
	return nil
}

// LoadCorefile Loads the Corefile from string
func (c *CoreDNSConfigEditor) LoadCorefile(corefile string) error {
	// Load the Corefile
	parsed, err := caddyfile.Parse(
		"config",
		strings.NewReader(corefile),
		nil,
	)
	if err != nil {
		return err
	}
	c.serverBlocks = parsed
	c.text = corefile
	return nil
}

// addInLine adds a new line in the Corefile at the specified line number
func (c *CoreDNSConfigEditor) addInLine(line int, text string) error {
	lines := strings.Split(c.text, "\n")
	lines = append(lines[:line-1], append([]string{text}, lines[line-1:]...)...)
	newText := strings.Join(lines, "\n")
	return c.LoadCorefile(newText)
}

// removeFileEntry removes all file entries from the Corefile
func (c *CoreDNSConfigEditor) removeFileEntry() error {
	if !strings.HasSuffix(c.text, "\n") {
		c.text += "\n"
	}

	lines := strings.Split(c.text, "\n")
	newLines := make([]string, len(lines))
	toRemove := fmt.Sprintf("file %s", ZoneFilePath)
	for _, line := range lines {
		if !strings.Contains(line, toRemove) {
			newLines = append(newLines, line)
		}
	}

	c.text = strings.Join(newLines, "\n")
	c.text = strings.TrimSpace(c.text)
	c.text += "\n"
	return c.LoadCorefile(c.text)
}

// getLineToAppend returns the line number to append the new entry
func (c *CoreDNSConfigEditor) getLineToAppend() int {
	// TODO: Do something smarter (e.g. after a specific plugin)
	return c.getMaxLine() - 1
}

// getLineToAppendHosts returns the line number to append the new entry in the hosts plugin
//
// The line number is calculated as follows:
// 1. If the forward plugin is present, the line number is the line number of the forward plugin -1
// 2. If the forward plugin is not present, the line number is the line number of the last line in the Corefile
func (c *CoreDNSConfigEditor) getLineToAppendHosts() int {
	// Initialize the line number to the maximum possible to ensure it's always after the current content
	lineToAppend := c.getMaxLine()
	// Iterate through each server block to find the 'forward' plugin
	for _, block := range c.serverBlocks {
		if tokens, ok := block.Tokens["forward"]; ok {
			// If the 'forward' plugin is found, update the line number to the line before the 'forward' plugin starts
			for _, token := range tokens {
				if token.Line < lineToAppend {
					lineToAppend = token.Line
					break // Assuming only one 'forward' entry is relevant
				}
			}
		}
	}
	return lineToAppend
}

func (c *CoreDNSConfigEditor) removeHosts() error {
	if !strings.HasSuffix(c.text, "\n") {
		c.text += "\n"
	}

	lines := strings.Split(c.text, "\n")
	newLines := make([]string, len(lines))
	found := false
	for _, line := range lines {
		if strings.Contains(line, "hosts {") {
			found = true
		}
		if found && strings.Contains(line, "}") {
			found = false
			continue
		}
		if !found {
			newLines = append(newLines, line)
		}
	}

	c.text = strings.Join(newLines, "\n")
	c.text = strings.TrimSpace(c.text)
	c.text += "\n"
	return c.LoadCorefile(c.text)
}

// removeTrailingDot removes the trailing dot from the zone name
func removeTrailingDot(zone string) string {
	if strings.HasSuffix(zone, ".") {
		return zone[:len(zone)-1]
	}
	return zone
}

// addTrailingDot adds a trailing dot to the zone name
func addTrailingDot(zone string) string {
	if !strings.HasSuffix(zone, ".") {
		return zone + "."
	}
	return zone
}
