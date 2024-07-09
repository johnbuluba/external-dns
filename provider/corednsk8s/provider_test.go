package corednsk8s

import (
	"errors"
	"flag"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func createRealKubernetesClient() kubernetes.Interface {
	// Create Kubernetes client
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

type CoreDNSk8sProviderTestSuite struct {
	suite.Suite
	provider *coreDNSk8sProvider
}

func (s *CoreDNSk8sProviderTestSuite) SetupTest() {
	// Initialize provider with dryRun mode enabled to avoid making actual changes
	//provider, err := NewCoreDNSProvider(true, createRealKubernetesClient())
	//s.Require().NoError(err)
	//s.provider = provider
}

func (s *CoreDNSk8sProviderTestSuite) TestInitialize_Success() {
	err := errors.New("not implemented")
	s.Require().NoError(err, "Initialize should not return an error")
}

// Add more tests here
func TestCoreDNSk8sProviderTestSuite(t *testing.T) {
	suite.Run(t, new(CoreDNSk8sProviderTestSuite))
}
