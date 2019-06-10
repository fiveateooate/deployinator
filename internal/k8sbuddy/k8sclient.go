package k8sbuddy

import (
	"fmt"
	"os"

	"github.com/wsxiaoys/terminal/color"
	"k8s.io/client-go/kubernetes"
	// needed to do gcp auth
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Connect return a connection to k8s api
func Connect(incluster bool, context string) *kubernetes.Clientset {
	var (
		clientset *kubernetes.Clientset
		err       error
	)
	if incluster {
		fmt.Println("Using incluster config")
	} else {
		clientset, err = ExternalClient(context)
		if err != nil {
			color.Printf("@rFailed to connect to k8s: %s\n", err)
			os.Exit(1)
		}
	}
	return clientset
}

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
