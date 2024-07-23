package manager

import (
	"context"
	"fmt"

	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/plan"
	"sigs.k8s.io/external-dns/provider/corednsk8s/editor"
	"sigs.k8s.io/external-dns/provider/corednsk8s/k8s"
)

type HostsManager struct {
	client           kubernetes.Interface
	coreDNSEditor    *editor.CoreDNSConfigEditor
	coreDNSConfigMap *k8s.CoreDNSConfigMap
}

// NewHostsManager creates a new RFC1035 manager.
func NewHostsManager(client kubernetes.Interface, configMap, ns string) (*HostsManager, error) {
	m := &HostsManager{
		client:           client,
		coreDNSEditor:    editor.NewCoreDNSConfigEditor(),
		coreDNSConfigMap: k8s.NewConfigMap(client.CoreV1(), ns, configMap),
	}
	err := m.reload(context.Background())
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *HostsManager) Records(ctx context.Context) ([]*endpoint.Endpoint, error) {
	// Get all records
	records := make([]*endpoint.Endpoint, 0)
	for _, host := range m.coreDNSEditor.GetHosts() {
		records = append(records, endpoint.NewEndpoint(host.Header().Name, endpoint.RecordTypeA, host.A.String()))
	}
	return records, nil
}

func (m *HostsManager) ApplyChanges(ctx context.Context, changes *plan.Changes) error {
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
func (m *HostsManager) applyCreate(endpoints []*endpoint.Endpoint) error {
	for _, endPt := range endpoints {
		logFields := log.Fields{
			"record": endPt.DNSName,
			"type":   endPt.RecordType,
		}
		log.WithFields(logFields).Debug("Creating record")
		record, ok := endpointToRecord(endPt).(*dns.A)
		if !ok {
			log.WithFields(logFields).Warn("Could not map record to DNS entry")
			continue
		}
		err := m.coreDNSEditor.AddHost(*record)
		if err != nil {
			return err
		}
	}
	return nil
}

// applyUpdate applies the update changes from the endpoints array
func (m *HostsManager) applyUpdate(endpointsOld, endpointsNew []*endpoint.Endpoint) error {
	if len(endpointsOld) != len(endpointsNew) {
		return fmt.Errorf("old and new endpoints length mismatch")
	}
	for i, endPtOld := range endpointsOld {
		endPtNew := endpointsNew[i]
		logFields := log.Fields{
			"old": map[string]interface{}{"record": endPtOld.DNSName, "type": endPtOld.RecordType},
			"new": map[string]interface{}{"record": endPtNew.DNSName, "type": endPtNew.RecordType},
		}
		log.WithFields(logFields).Debug("Updating record")
		recordOld, ok := endpointToRecord(endPtOld).(*dns.A)
		if !ok {
			log.WithFields(logFields).Warn("Could not map record to DNS entry")
			continue
		}
		recordNew, ok := endpointToRecord(endPtNew).(*dns.A)
		if !ok {
			log.WithFields(logFields).Warn("Could not map record to DNS entry")
			continue
		}
		err := m.coreDNSEditor.UpdateHost(*recordOld, *recordNew)
		if err != nil {
			return err
		}
	}
	return nil
}

// applyDelete applies the delete changes from the endpoints array
func (m *HostsManager) applyDelete(endpoints []*endpoint.Endpoint) error {
	for _, endPt := range endpoints {
		logFields := log.Fields{
			"record": endPt.DNSName,
			"type":   endPt.RecordType,
		}
		log.WithFields(logFields).Debug("Deleting record")
		record, ok := endpointToRecord(endPt).(*dns.A)
		if !ok {
			log.WithFields(logFields).Warn("Could not map record to DNS entry")
			continue
		}
		err := m.coreDNSEditor.RemoveHost(*record)
		if err != nil {
			return err
		}
	}
	return nil
}

// reload reloads the Corefile
func (m *HostsManager) reload(ctx context.Context) error {
	// Reload the Corefile
	if corefile, err := m.coreDNSConfigMap.GetCoreDNSConfig(ctx); err == nil {
		if err = m.coreDNSEditor.LoadCorefile(corefile); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

// commit saves the changes to the Corefile and Zone file
func (m *HostsManager) commit(ctx context.Context) error {
	// Update Corefile
	if err := m.coreDNSConfigMap.UpdateCoreDNSConfig(ctx, m.coreDNSEditor.GetConfig()); err != nil {
		return err
	}
	return nil
}
