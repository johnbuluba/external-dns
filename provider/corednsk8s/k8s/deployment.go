package k8s

import (
	"context"
	"errors"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

const (
	ZonePath = "/etc/coredns/zone.db"
)

type CoreDNSDeployment struct {
	ns        string
	name      string
	configMap string
	client    apiv1.AppsV1Interface
}

// NewDeployment creates a new CoreDNSDeployment
func NewDeployment(client apiv1.AppsV1Interface, ns, name, configMap string) *CoreDNSDeployment {
	return &CoreDNSDeployment{
		ns:        ns,
		name:      name,
		configMap: configMap,
		client:    client,
	}
}

// MountZoneFile mounts the Zone file in the CoreDNS Deployment
func (c *CoreDNSDeployment) MountZoneFile(ctx context.Context) error {
	deployment, err := c.client.Deployments(c.ns).Get(ctx, c.name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Add the volume
	added, err := c.addVolume(deployment)
	if err != nil {
		return err
	}
	if !added {
		return nil
	}

	// Update the deployment
	_, err = c.client.Deployments(c.ns).Update(ctx, deployment, metav1.UpdateOptions{})
	return err
}

// addVolume adds volume for the Zone file in the CoreDNS Deployment
//
// Returns true  if the volume was added
func (c *CoreDNSDeployment) addVolume(deployment *appsv1.Deployment) (bool, error) {
	// Add the Zone file volume
	for i, volume := range deployment.Spec.Template.Spec.Volumes {
		if cfg := volume.ConfigMap; cfg != nil && cfg.Name == c.configMap {
			// If there is an existing item skip
			for _, item := range cfg.Items {
				if item.Key == ZoneKey {
					return false, nil
				}
			}
			// Append the new item
			added := append(deployment.Spec.Template.Spec.Volumes[i].VolumeSource.ConfigMap.Items, corev1.KeyToPath{
				Key:  ZoneKey,
				Path: ZoneKey,
			})
			deployment.Spec.Template.Spec.Volumes[i].VolumeSource.ConfigMap.Items = added
			return true, nil
		}
	}
	return false, errors.New("coreDNS configMap not found in volumes")
}
