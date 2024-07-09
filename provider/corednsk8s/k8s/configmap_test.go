package k8s

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type CoreDNSConfigMapTestSuite struct {
	suite.Suite
	client    corev1.CoreV1Interface
	configMap *CoreDNSConfigMap
}

func (s *CoreDNSConfigMapTestSuite) SetupTest() {

}

// BeforeTest is a function to be executed right before the test
func (s *CoreDNSConfigMapTestSuite) BeforeTest(suiteName, testName string) {
	s.client = fake.NewSimpleClientset().CoreV1()
	s.configMap = NewConfigMap(s.client, "kube-system", "coredns")
}

func (s *CoreDNSConfigMapTestSuite) TestGetCoreDNSConfig_Success() {

	_, err := s.client.ConfigMaps("kube-system").Create(context.TODO(), &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "coredns", Namespace: "kube-system"},
		Data:       map[string]string{CorefileKey: "CoreDNS config data"},
	}, metav1.CreateOptions{})
	s.Require().NoError(err)

	result, err := s.configMap.GetCoreDNSConfig(context.TODO())
	s.Require().NoError(err)
	s.Equal("CoreDNS config data", result)
}

func (s *CoreDNSConfigMapTestSuite) TestGetCoreDNSConfig_Failure_KeyNotFound() {
	_, err := s.client.ConfigMaps("kube-system").Create(context.TODO(), &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "coredns", Namespace: "kube-system"},
		Data:       map[string]string{"Invalid": "CoreDNS config data"},
	}, metav1.CreateOptions{})
	s.Require().NoError(err)

	_, err = s.configMap.GetCoreDNSConfig(context.TODO())
	s.Require().Error(err)
	s.Equal("Corefile not found", err.Error())
}

func (s *CoreDNSConfigMapTestSuite) TestGetCoreDNSConfig_Failure_ConfigMapNotFound() {
	_, err := s.configMap.GetCoreDNSConfig(context.TODO())
	s.Require().Error(err)
	s.Equal("configmaps \"coredns\" not found", err.Error())
}

func (s *CoreDNSConfigMapTestSuite) TestUpdateCoreDNSConfig_Success() {
	corefileContent := "Updated CoreDNS config data"
	_, err := s.client.ConfigMaps("kube-system").Create(context.TODO(), &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "coredns", Namespace: "kube-system"},
		Data:       map[string]string{CorefileKey: "Initial CoreDNS config data"},
	}, metav1.CreateOptions{})
	s.Require().NoError(err)

	err = s.configMap.UpdateCoreDNSConfig(context.Background(), corefileContent)
	s.Require().NoError(err)

	updatedCfgMap, err := s.client.ConfigMaps("kube-system").Get(context.Background(), "coredns", metav1.GetOptions{})
	s.Require().NoError(err)
	s.Equal(corefileContent, updatedCfgMap.Data[CorefileKey], "Corefile should be updated with new content")
}

func (s *CoreDNSConfigMapTestSuite) TestUpdateCoreDNSConfig_Failure_NotFound() {
	corefileContent := "Updated CoreDNS config data"

	err := s.configMap.UpdateCoreDNSConfig(context.Background(), corefileContent)
	s.Require().Error(err)
	s.Contains(err.Error(), "not found", "Should return an error for non-existent CoreDNSConfigMap")
}

func (s *CoreDNSConfigMapTestSuite) TestGetZone_Success() {
	zoneContent := "example.com."
	_, err := s.client.ConfigMaps("kube-system").Create(context.TODO(), &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "coredns", Namespace: "kube-system"},
		Data:       map[string]string{ZoneKey: zoneContent},
	}, metav1.CreateOptions{})
	s.Require().NoError(err)

	result, err := s.configMap.GetZone(context.TODO())
	s.Require().NoError(err)
	s.Equal(zoneContent, result)
}

func (s *CoreDNSConfigMapTestSuite) TestGetZone_Failure_ZoneNotFound() {
	_, err := s.client.ConfigMaps("kube-system").Create(context.TODO(), &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "coredns", Namespace: "kube-system"},
		Data:       map[string]string{"InvalidKey": "example.com."},
	}, metav1.CreateOptions{})
	s.Require().NoError(err)

	_, err = s.configMap.GetZone(context.TODO())
	s.Require().Error(err)
	s.Equal("zone not found", err.Error())
}

func (s *CoreDNSConfigMapTestSuite) TestUpdateZone_Success() {
	initialZoneContent := "initial.example.com."
	updatedZoneContent := "updated.example.com."
	_, err := s.client.ConfigMaps("kube-system").Create(context.TODO(), &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "coredns", Namespace: "kube-system"},
		Data:       map[string]string{ZoneKey: initialZoneContent},
	}, metav1.CreateOptions{})
	s.Require().NoError(err)

	err = s.configMap.UpdateZone(context.Background(), updatedZoneContent)
	s.Require().NoError(err)

	updatedCfgMap, err := s.client.ConfigMaps("kube-system").Get(context.Background(), "coredns", metav1.GetOptions{})
	s.Require().NoError(err)
	s.Equal(updatedZoneContent, updatedCfgMap.Data[ZoneKey], "Zone should be updated with new content")
}

func (s *CoreDNSConfigMapTestSuite) TestUpdateZone_Failure_NotFound() {
	updatedZoneContent := "updated.example.com."

	err := s.configMap.UpdateZone(context.Background(), updatedZoneContent)
	s.Require().Error(err)
	s.Contains(err.Error(), "not found", "Should return an error for non-existent CoreDNSConfigMap")
}

func TestConfigMapTestSuite(t *testing.T) {
	suite.Run(t, new(CoreDNSConfigMapTestSuite))
}
