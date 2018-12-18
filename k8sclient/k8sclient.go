package k8sclient

import (
	"k8s.io/client-go/kubernetes"
	// needed to do gcp auth
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//InClusterClient - return a k8s client using incluster config
func InClusterClient() (*kubernetes.Clientset, error) {
	var clientset *kubernetes.Clientset
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return clientset, err
	}
	// creates the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return clientset, err
	}
	return clientset, err
}

// ExternalClient - use out of cluster config
func ExternalClient(context string) (*kubernetes.Clientset, error) {
	var clientset *kubernetes.Clientset
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{CurrentContext: context}).ClientConfig()
	if err != nil {
		return clientset, err
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return clientset, err
	}
	return clientset, nil
}
