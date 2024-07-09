package manager

import (
	"context"

	"github.com/miekg/dns"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/plan"
	"sigs.k8s.io/external-dns/provider/corednsk8s/editor"
	"sigs.k8s.io/external-dns/provider/corednsk8s/k8s"
)

type RFC1035Manager struct {
	client            kubernetes.Interface
	zoneEditor        *editor.ZoneEditor
	coreDNSEditor     *editor.CoreDNSConfigEditor
	coreDNSConfigMap  *k8s.CoreDNSConfigMap
	coreDNSDeployment *k8s.CoreDNSDeployment
}

// NewRFC1035Manager creates a new RFC1035 manager.
func NewRFC1035Manager(client kubernetes.Interface) (*RFC1035Manager, error) {
	m := &RFC1035Manager{
		client:            client,
		zoneEditor:        editor.NewZoneEditor(),
		coreDNSEditor:     editor.NewCoreDNSConfigEditor(),
		coreDNSConfigMap:  k8s.NewConfigMap(client.CoreV1(), "kube-system", "coredns"),
		coreDNSDeployment: k8s.NewDeployment(client.AppsV1(), "kube-system", "coredns", "coredns"),
	}

	// Verify that the Zone file is exist in ConfigMap
	_, err := m.coreDNSConfigMap.GetZone(context.Background())
	switch {
	case err == nil:
		// Zone file exists
	case err.Error() == "zone not found":
		// Create the Zone file in the ConfigMap
		err = m.coreDNSConfigMap.UpdateZone(context.Background(), "")
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
	}

	// Verify that the Zone file is mounted
	err = m.coreDNSDeployment.MountZoneFile(context.Background())
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m RFC1035Manager) Records(ctx context.Context) ([]*endpoint.Endpoint, error) {
	// Get Zone from ConfigMap
	zone, err := m.coreDNSConfigMap.GetZone(ctx)
	if err != nil {
		return nil, err
	}
	// Parse it
	err = m.zoneEditor.LoadZone(zone)
	// Get all records for every zone
	records := make([]*endpoint.Endpoint, 0)
	for _, zone := range m.zoneEditor.GetZones() {
		for _, record := range m.zoneEditor.GetAllRecords(zone) {
			switch r := record.(type) {
			case *dns.A:
				records = append(records, endpoint.NewEndpoint(r.Header().Name, endpoint.RecordTypeA, r.A.String()))
			case *dns.CNAME:
				records = append(records, endpoint.NewEndpoint(r.Header().Name, endpoint.RecordTypeCNAME, r.Target))
			case *dns.TXT:
				records = append(records, endpoint.NewEndpoint(r.Header().Name, endpoint.RecordTypeTXT, r.Txt...))
			case *dns.SRV:
				records = append(records, endpoint.NewEndpoint(r.Header().Name, endpoint.RecordTypeSRV, r.Target))
			case *dns.PTR:
				records = append(records, endpoint.NewEndpoint(r.Header().Name, endpoint.RecordTypePTR, r.Ptr))
			case *dns.MX:
				records = append(records, endpoint.NewEndpoint(r.Header().Name, endpoint.RecordTypeMX, r.Mx))
			case *dns.AAAA:
				records = append(records, endpoint.NewEndpoint(r.Header().Name, endpoint.RecordTypeAAAA, r.AAAA.String()))
			case *dns.NS:
				records = append(records, endpoint.NewEndpoint(r.Header().Name, endpoint.RecordTypeNS, r.Ns))
			default:
				//TODO implement me
			}
		}
	}
	return records, nil
}

func (m RFC1035Manager) ApplyChanges(ctx context.Context, changes *plan.Changes) error {
	//TODO implement me
	panic("implement me")
}
