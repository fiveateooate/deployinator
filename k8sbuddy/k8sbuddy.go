package k8sbuddy

import (
	"fmt"

	"github.com/wsxiaoys/terminal/color"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetDaemonset return info on single deployment
func GetDaemonset(appName string, namespaceName string, clientset *kubernetes.Clientset) (*appsv1.DaemonSet, error) {
	color.Printf("@cSearching for daemonset ...")
	daemonsetClient := clientset.AppsV1().DaemonSets(namespaceName)
	tmp, err := daemonsetClient.Get(appName, metav1.GetOptions{})
	if err == nil {
		color.Printf("found %s\n", tmp.Name)
	} else {
		color.Printf("not found\n")
		tmp = nil
	}
	return tmp, err
}

// GetDeployment load info for single deployment
func GetDeployment(appName string, namespaceName string, clientset *kubernetes.Clientset) (*appsv1.Deployment, error) {
	color.Printf("@cSearching for deployment ...")
	deploymentsClient := clientset.AppsV1().Deployments(namespaceName)
	tmp, err := deploymentsClient.Get(appName, metav1.GetOptions{})
	if err == nil {
		color.Printf("found %s\n", tmp.Name)
	} else {
		color.Printf("not found\n")
		tmp = nil
	}
	return tmp, err
}

// GetStatefulset return info on single deployment
func GetStatefulset(appName string, namespaceName string, clientset *kubernetes.Clientset) (*appsv1.StatefulSet, error) {
	color.Printf("@cSearching for statefulset...")
	statefulsetClient := clientset.AppsV1().StatefulSets(namespaceName)
	tmp, err := statefulsetClient.Get(appName, metav1.GetOptions{})
	if err == nil {
		color.Printf("found %s\n", tmp.Name)
	} else {
		color.Printf("not found\n")
		tmp = nil
	}
	return tmp, err
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

// GetAppInfo returns something about appName/namespaceName
func (k8sapp *K8sApp) GetAppInfo(appName string, namespaceName string, clientset *kubernetes.Clientset) {

	if deployment, err := GetDeployment(appName, namespaceName, clientset); err == nil {
		k8sapp.DP = deployment
		k8sapp.Kind = "deployment"
		return
	}
	if statefulset, err := GetStatefulset(appName, namespaceName, clientset); err == nil {
		k8sapp.SS = statefulset
		k8sapp.Kind = "statefulset"
		return
	}
	if daemonset, err := GetDaemonset(appName, namespaceName, clientset); err == nil {
		k8sapp.DS = daemonset
		k8sapp.Kind = "daemonset"
		return
	}
	k8sapp.Kind = ""
}
