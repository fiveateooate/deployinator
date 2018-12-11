package k8sbuddy

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetDeployment return info on single deployment
func GetDeployment(appName string, namespaceName string, clientset *kubernetes.Clientset) (*appsv1.Deployment, error) {
	deploymentsClient := clientset.AppsV1().Deployments(namespaceName)
	result, err := deploymentsClient.Get(appName, metav1.GetOptions{})
	return result, err
}

// GetDeployments return list of all deployments in all namespaces
func GetDeployments(clientset *kubernetes.Clientset, namespace apiv1.Namespace) {
	deploymentsClient := clientset.AppsV1().Deployments(namespace.GetName())
	deployments, err := deploymentsClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listing deployments in namespace %q:\n", namespace.GetName())
	for _, d := range deployments.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}
}

// GetNamespaces return list of all namespaces
func GetNamespaces(clientset *kubernetes.Clientset) *apiv1.NamespaceList {
	namespaces, _ := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	return namespaces
}
