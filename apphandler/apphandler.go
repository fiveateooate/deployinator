package apphandler

import (
	"fmt"

	"github.com/fiveateooate/deployinator/artifactsbuddy"
	"github.com/fiveateooate/deployinator/helmbuddy"
	"github.com/fiveateooate/deployinator/k8sbuddy"
	"k8s.io/client-go/kubernetes"
)

// ManageApp do stuff for a single app
func ManageApp(appName string, namespace string, kubecontext string, clientset *kubernetes.Clientset) {
	fmt.Printf("Getting info for deployment %s\n", appName)
	deployment, err := k8sbuddy.GetDeployment(appName, namespace, clientset)
	if err == nil {
		fmt.Printf("Found k8s deployment: %s\n", deployment.Name)
	} else {
		fmt.Println(err)
	}
	helmRelease := helmbuddy.GetRelease(appName, namespace, kubecontext)
	fmt.Printf("Found release: %s\n", helmRelease.Name)

	if helmRelease.Name != "" && deployment != nil {
		fmt.Printf("Upgrading %s\n", appName)
		artifactsbuddy.GetPkgs(appName)
	}
}
