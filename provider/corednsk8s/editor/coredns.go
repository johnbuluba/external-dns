package editor

import (
	"fmt"
	"math"
	"strings"

	"github.com/coredns/caddy/caddyfile"
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
		cleanUpZones[i] = cleanUpZoneName(z)
	}
	// Update the file configuration
	appendToLine := c.getLineToAppend()
	entry := fmt.Sprintf("    file %s %s", ZoneFilePath, strings.Join(cleanUpZones, " "))
	return c.addInLine(appendToLine, entry)
}

// AddZone adds a new zone in the file plugin
func (c *CoreDNSConfigEditor) AddZone(zone string) error {
	cleanZoneName := cleanUpZoneName(zone)
	// Get existing zones
	zones := c.GetZones()
	// Check if the zone already exists
	for _, z := range zones {
		if z == cleanZoneName {
			return nil
		}
	}
	// Remove the file entry
	err := c.removeFileEntry()
	if err != nil {
		return err
	}
	// Update the file configuration
	zones = append(zones, cleanZoneName)
	appendToLine := c.getLineToAppend()
	entry := fmt.Sprintf("    file %s %s", ZoneFilePath, strings.Join(zones, " "))
	return c.addInLine(appendToLine, entry)
}

// RemoveZone removes a zone from the file plugin
func (c *CoreDNSConfigEditor) RemoveZone(zone string) error {
	cleanZoneName := cleanUpZoneName(zone)
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
	return c.getMaxLine() - 2
}

// cleanUpZoneName removes the trailing dot from the zone name
func cleanUpZoneName(zone string) string {
	if strings.HasSuffix(zone, ".") {
		return zone[:len(zone)-1]
	}
	return zone
}
