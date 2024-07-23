package manager

import (
	"context"
	"fmt"
	"strings"

	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
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
func NewRFC1035Manager(client kubernetes.Interface, deployment, configMap, ns string) (*RFC1035Manager, error) {
	m := &RFC1035Manager{
		client:            client,
		zoneEditor:        editor.NewZoneEditor(),
		coreDNSEditor:     editor.NewCoreDNSConfigEditor(),
		coreDNSConfigMap:  k8s.NewConfigMap(client.CoreV1(), ns, configMap),
		coreDNSDeployment: k8s.NewDeployment(client.AppsV1(), ns, deployment, configMap),
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

func (m *RFC1035Manager) Records(ctx context.Context) ([]*endpoint.Endpoint, error) {
	// Get Zone from ConfigMap
	zoneCfg, err := m.coreDNSConfigMap.GetZone(ctx)
	if err != nil {
		return nil, err
	}

	// Parse it
	if err = m.zoneEditor.LoadZone(zoneCfg); err != nil {
		return nil, err
	}
	// Get all records for every zone
	records := make([]*endpoint.Endpoint, 0)
	for _, zone := range m.zoneEditor.GetZones() {
		for _, record := range m.zoneEditor.GetAllRecords(zone) {
			switch r := record.(type) {
			case *dns.A:
				records = append(records, endpoint.NewEndpointWithTTL(r.Header().Name, endpoint.RecordTypeA, endpoint.TTL(r.Hdr.Ttl), r.A.String()))
			case *dns.CNAME:
				records = append(records, endpoint.NewEndpointWithTTL(r.Header().Name, endpoint.RecordTypeCNAME, endpoint.TTL(r.Hdr.Ttl), r.Target))
			case *dns.TXT:
				records = append(records, endpoint.NewEndpointWithTTL(r.Header().Name, endpoint.RecordTypeTXT, endpoint.TTL(r.Hdr.Ttl), r.Txt...))
			case *dns.SRV:
				records = append(records, endpoint.NewEndpointWithTTL(r.Header().Name, endpoint.RecordTypeSRV, endpoint.TTL(r.Hdr.Ttl), r.Target))
			case *dns.PTR:
				records = append(records, endpoint.NewEndpointWithTTL(r.Header().Name, endpoint.RecordTypePTR, endpoint.TTL(r.Hdr.Ttl), r.Ptr))
			case *dns.MX:
				records = append(records, endpoint.NewEndpointWithTTL(r.Header().Name, endpoint.RecordTypeMX, endpoint.TTL(r.Hdr.Ttl), r.Mx))
			case *dns.AAAA:
				records = append(records, endpoint.NewEndpointWithTTL(r.Header().Name, endpoint.RecordTypeAAAA, endpoint.TTL(r.Hdr.Ttl), r.AAAA.String()))
			case *dns.NS:
				records = append(records, endpoint.NewEndpointWithTTL(r.Header().Name, endpoint.RecordTypeNS, endpoint.TTL(r.Hdr.Ttl), r.Ns))
			default:
				//TODO implement me
			}
		}
	}
	return records, nil
}

func (m *RFC1035Manager) ApplyChanges(ctx context.Context, changes *plan.Changes) error {
	// Reload the Corefile and Zone file
	if err := m.reload(ctx); err != nil {
		return err
	}
	// Create new records
	if createChanges := changes.Create; len(createChanges) > 0 {
		if err := m.applyCreate(createChanges); err != nil {
			return err
		}
	}
	// Update records
	if updateChanges := changes.UpdateNew; len(updateChanges) > 0 {
		if err := m.applyUpdate(changes.UpdateOld, updateChanges); err != nil {
			return err
		}
	}
	// Delete records
	if deleteChanges := changes.Delete; len(deleteChanges) > 0 {
		if err := m.applyDelete(deleteChanges); err != nil {
			return err
		}
	}
	// Apply changes
	return m.commit(ctx)
}

// applyCreate applies the create changes from the endpoints array
func (m *RFC1035Manager) applyCreate(endpoints []*endpoint.Endpoint) error {
	for _, endPt := range endpoints {
		logFields := log.Fields{
			"record": endPt.DNSName,
			"type":   endPt.RecordType,
		}
		log.WithFields(logFields).Debug("Creating record")
		zone := getZone(endPt.DNSName)
		record := endpointToRecord(endPt)
		if record == nil {
			log.WithFields(logFields).Warn("Could not map record to DNS entry")
			continue
		}
		m.zoneEditor.AddRecord(zone, record)
	}
	return nil
}

// applyUpdate applies the update changes from the endpoints array
func (m *RFC1035Manager) applyUpdate(endpointsOld, endpointsNew []*endpoint.Endpoint) error {
	if len(endpointsOld) != len(endpointsNew) {
		return fmt.Errorf("old and new endpoints length mismatch")
	}
	for i, endPtOld := range endpointsOld {
		endPtNew := endpointsNew[i]
		logFields := log.Fields{
			"old":    map[string]interface{}{"record": endPtOld.DNSName, "type": endPtOld.RecordType},
			"new":    map[string]interface{}{"record": endPtNew.DNSName, "type": endPtNew.RecordType},
			"record": endPtOld.DNSName,
			"type":   endPtOld.RecordType,
		}
		log.WithFields(logFields).Debug("Updating record")
		zone := getZone(endPtOld.DNSName)
		recordOld := endpointToRecord(endPtOld)
		recordNew := endpointToRecord(endPtNew)
		if recordOld == nil || recordNew == nil {
			log.WithFields(logFields).Warn("Could not map record to DNS entry")
			continue
		}
		m.zoneEditor.UpdateRecord(zone, recordOld, recordNew)
	}
	return nil
}

// applyDelete applies the delete changes from the endpoints array
func (m *RFC1035Manager) applyDelete(endpoints []*endpoint.Endpoint) error {
	for _, endPt := range endpoints {
		logFields := log.Fields{
			"record": endPt.DNSName,
			"type":   endPt.RecordType,
		}
		log.WithFields(logFields).Debug("Deleting record")
		zone := getZone(endPt.DNSName)
		record := endpointToRecord(endPt)
		if record == nil {
			log.WithFields(logFields).Warn("Could not map record to DNS entry")
			continue
		}
		m.zoneEditor.DeleteRecord(zone, record)
	}
	return nil
}

// reload reloads the Corefile and Zone file
func (m *RFC1035Manager) reload(ctx context.Context) error {
	// Reload the Corefile
	if corefile, err := m.coreDNSConfigMap.GetCoreDNSConfig(ctx); err == nil {
		if err = m.coreDNSEditor.LoadCorefile(corefile); err != nil {
			return err
		}
	} else {
		return err
	}

	// Reload the Zone file
	if zone, err := m.coreDNSConfigMap.GetZone(ctx); err == nil {
		if err = m.zoneEditor.LoadZone(zone); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

// commit saves the changes to the Corefile and Zone file
func (m *RFC1035Manager) commit(ctx context.Context) error {
	// Update zones in corefile
	zones := m.zoneEditor.GetZones()
	if err := m.coreDNSEditor.SetZones(zones); err != nil {
		return err
	}

	// Update records in zone file
	if err := m.coreDNSConfigMap.UpdateCoreDNSConfig(ctx, m.coreDNSEditor.GetConfig()); err != nil {
		return err
	}
	// Update zones in zone file
	if err := m.coreDNSConfigMap.UpdateZone(ctx, m.zoneEditor.RenderZone()); err != nil {
		return err
	}
	return nil
}

// getZone returns the n-1 parts of the TLD (e.g. "test.example.com" -> "example.com.")
func getZone(domain string) string {
	parts := dns.SplitDomainName(domain)
	if len(parts) < 3 {
		panic("invalid domain")
	}

	zone := strings.Join(parts[1:], ".")
	zone = fmt.Sprintf("%s.", zone)
	return zone
}
