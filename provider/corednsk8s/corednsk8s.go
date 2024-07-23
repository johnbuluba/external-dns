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

type CoreDNSConfig struct {
	CoreDNSDeployment string
	CoreDNSConfigMap  string
	CoreDNSNamespace  string
}

// NewCoreDNSProvider creates a new CoreDNS provider.
func NewCoreDNSProvider(domainFilter endpoint.DomainFilter, c source.ClientGenerator, cfg CoreDNSConfig, dryRun bool) (*coreDNSk8sProvider, error) {
	client, err := c.KubeClient()
	if err != nil {
		return nil, err
	}
	p := &coreDNSk8sProvider{
		dryRun:     dryRun,
		client:     client,
		zoneEditor: editor.NewZoneEditor(),
	}
	//if mgr, err := manager.NewRFC1035Manager(client, cfg.CoreDNSDeployment, cfg.CoreDNSConfigMap, cfg.CoreDNSNamespace); err == nil {
	//	p.manager = mgr
	//} else {
	//	return nil, err
	//}
	mgr, err := manager.NewHostsManager(client, cfg.CoreDNSConfigMap, cfg.CoreDNSNamespace)
	if err != nil {
		return nil, err
	}
	p.manager = mgr
	return p, nil
}

func (c *coreDNSk8sProvider) Records(ctx context.Context) ([]*endpoint.Endpoint, error) {
	return c.manager.Records(ctx)
}

func (c *coreDNSk8sProvider) ApplyChanges(ctx context.Context, changes *plan.Changes) error {
	return c.manager.ApplyChanges(ctx, changes)
}
