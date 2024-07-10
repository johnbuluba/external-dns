package corednsk8s

import (
	"context"

	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/plan"
	"sigs.k8s.io/external-dns/provider"
	"sigs.k8s.io/external-dns/provider/corednsk8s/editor"
	"sigs.k8s.io/external-dns/provider/corednsk8s/k8s"
	"sigs.k8s.io/external-dns/provider/corednsk8s/manager"
	"sigs.k8s.io/external-dns/source"
)

type coreDNSk8sProvider struct {
	provider.BaseProvider
	dryRun            bool
	client            kubernetes.Interface
	zoneEditor        *editor.ZoneEditor
	coreDNSConfig     *editor.CoreDNSConfigEditor
	coreDNSConfigMap  *k8s.CoreDNSConfigMap
	coreDNSDeployment *k8s.CoreDNSDeployment
	manager           CoreDNSManager
}

// NewCoreDNSProvider creates a new CoreDNS provider.
func NewCoreDNSProvider(domainFilter endpoint.DomainFilter, c source.ClientGenerator, dryRun bool) (*coreDNSk8sProvider, error) {
	client, err := c.KubeClient()
	if err != nil {
		return nil, err
	}
	p := &coreDNSk8sProvider{
		dryRun:     dryRun,
		client:     client,
		zoneEditor: editor.NewZoneEditor(),
	}
	if mgr, err := manager.NewRFC1035Manager(client); err == nil {
		p.manager = mgr
	} else {
		return nil, err
	}
	return p, nil
}

func (c *coreDNSk8sProvider) Records(ctx context.Context) ([]*endpoint.Endpoint, error) {
	return c.manager.Records(ctx)
}

func (c *coreDNSk8sProvider) ApplyChanges(ctx context.Context, changes *plan.Changes) error {
	return c.manager.ApplyChanges(ctx, changes)
}

//// Initializes the provider
//func (c *coreDNSk8sProvider) initialize() error {
//	// Create CoreDNSConfigMap client
//	c.coreDNSConfigMap = k8s.NewConfigMap(c.client.CoreV1(), "kube-system", "coredns")
//	// Configure the file plugin
//	config, err := c.coreDNSConfigMap.GetCoreDNSConfig(context.Background())
//	if err != nil {
//		return err
//	}
//	c.coreDNSConfig = editor.NewCoreDNSConfigEditor()
//	err = c.coreDNSConfig.LoadCorefile(config)
//	if err != nil {
//		return err
//	}
//	err = c.coreDNSConfig.AddZone("test.com")
//	err = c.coreDNSConfig.AddZone("test1.com")
//	if err != nil {
//		return err
//	}
//	err = c.coreDNSConfigMap.UpdateCoreDNSConfig(context.Background(), c.coreDNSConfig.GetConfig())
//	if err != nil {
//		return err
//	}
//
//	// Verify that the Zone file is mounted
//	c.coreDNSDeployment = k8s.NewDeployment(c.client.AppsV1(), "kube-system", "coredns", "coredns")
//	err = c.coreDNSDeployment.MountZoneFile(context.Background())
//	if err != nil {
//		return err
//	}
//
//	// Create a new ZoneEditor
//	c.zoneEditor = editor.NewZoneEditor()
//	// Load the zone
//	if err := c.zoneEditor.LoadZone(""); err != nil {
//		return err
//	}
//	c.zoneEditor.GetOrCreateZone("test.com.")
//	c.zoneEditor.GetOrCreateZone("test1.com.")
//	c.zoneEditor.AddARecord("test.com.", "test", "192.168.1.1")
//	c.zoneEditor.AddARecord("test1.com.", "test", "192.168.2.1")
//
//	// Initialize an empty zone from the CoreDNSConfigMap
//	err = c.coreDNSConfigMap.UpdateZone(context.Background(), c.zoneEditor.RenderZone())
//	if err != nil {
//		return err
//	}
//	return nil
//}
