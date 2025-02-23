package app

import (
	"os"

	"github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// ProvideWorkflowClient tries in-cluster config first, then falls back to
// local kubeconfig. Returns an Argo Workflows clientset.
func ProvideWorkflowClient() (versioned.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		// If not in a cluster, try KUBECONFIG or default kubeconfig file
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = clientcmd.RecommendedHomeFile
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}
	return versioned.NewForConfig(config)
}
