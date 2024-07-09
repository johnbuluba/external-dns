package k8s

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	apiv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

type CoreDNSDeploymentTestSuite struct {
	suite.Suite
	deployment *CoreDNSDeployment
	client     apiv1.AppsV1Interface
}

func (s *CoreDNSDeploymentTestSuite) SetupTest() {
	client := fake.NewSimpleClientset().AppsV1()
	s.client = client
	s.deployment = NewDeployment(client, "kube-system", "coredns", "coredns")
}

func (s *CoreDNSDeploymentTestSuite) TestMountZoneFile_Success() {
	// Create a deployment to be updated
	_, err := s.client.Deployments("kube-system").Create(context.TODO(), &appsv1.Deployment{
		ObjectMeta: v1.ObjectMeta{Name: "coredns", Namespace: "kube-system"},
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "coredns",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config-volume",
									MountPath: "/etc/coredns",
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "config-volume",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "coredns",
									},
									Items: []corev1.KeyToPath{{Key: "Corefile", Path: "Corefile"}},
								},
							},
						},
					},
				},
			},
		},
	}, v1.CreateOptions{})
	s.Require().NoError(err)

	err = s.deployment.MountZoneFile(context.TODO())
	s.Require().NoError(err)

	// Verify the deployment has been updated
	deployment, err := s.client.Deployments("kube-system").Get(context.TODO(), "coredns", v1.GetOptions{})
	s.Require().NoError(err)
	s.Len(deployment.Spec.Template.Spec.Volumes, 1, "Expected 1 volume in the deployment")
	s.Len(deployment.Spec.Template.Spec.Volumes[0].ConfigMap.Items, 2, "Expected 2 items in the volume")
	s.Len(deployment.Spec.Template.Spec.Containers[0].VolumeMounts, 1, "Expected 1 volume mount in the CoreDNS container")
}

func (s *CoreDNSDeploymentTestSuite) TestMountZoneFile_Failure_DeploymentNotFound() {
	err := s.deployment.MountZoneFile(context.TODO())
	s.Require().Error(err)
	s.Contains(err.Error(), "not found")
}

func TestCoreDNSDeploymentTestSuite(t *testing.T) {
	suite.Run(t, new(CoreDNSDeploymentTestSuite))
}
