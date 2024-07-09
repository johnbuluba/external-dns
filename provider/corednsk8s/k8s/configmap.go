package k8s

import (
	"context"
	"errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	CorefileKey = "Corefile"
	ZoneKey     = "Zone"
)

// CoreDNSConfigMap is a representation of a Kubernetes Coredns CoreDNSConfigMap
type CoreDNSConfigMap struct {
	ns     string
	name   string
	client apiv1.CoreV1Interface
}

// NewConfigMap creates a new CoreDNSConfigMap
func NewConfigMap(client apiv1.CoreV1Interface, ns, name string) *CoreDNSConfigMap {
	return &CoreDNSConfigMap{
		ns:     ns,
		name:   name,
		client: client,
	}
}

// GetCoreDNSConfig returns the CoreDNSConfigMap string
func (c *CoreDNSConfigMap) GetCoreDNSConfig(ctx context.Context) (string, error) {
	cfgMap, err := c.client.ConfigMaps(c.ns).Get(ctx, c.name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	if value, ok := cfgMap.Data[CorefileKey]; ok {
		return value, nil
	}
	return "", errors.New("Corefile not found")
}

// UpdateCoreDNSConfig updates the CoreDNSConfigMap with the new Corefile
func (c *CoreDNSConfigMap) UpdateCoreDNSConfig(ctx context.Context, corefile string) error {
	cfgMap, err := c.client.ConfigMaps(c.ns).Get(ctx, c.name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	cfgMap.Data[CorefileKey] = corefile
	_, err = c.client.ConfigMaps(c.ns).Update(ctx, cfgMap, metav1.UpdateOptions{})
	return err
}

// GetZone returns the Zone from the CoreDNSConfigMap
func (c *CoreDNSConfigMap) GetZone(ctx context.Context) (string, error) {
	cfgMap, err := c.client.ConfigMaps(c.ns).Get(ctx, c.name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	if value, ok := cfgMap.Data[ZoneKey]; ok {
		return value, nil
	}
	return "", errors.New("zone not found")
}

// UpdateZone updates the CoreDNSConfigMap with the new Zone
func (c *CoreDNSConfigMap) UpdateZone(ctx context.Context, zone string) error {
	cfgMap, err := c.client.ConfigMaps(c.ns).Get(ctx, c.name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	cfgMap.Data[ZoneKey] = zone
	_, err = c.client.ConfigMaps(c.ns).Update(ctx, cfgMap, metav1.UpdateOptions{})
	return err
}
