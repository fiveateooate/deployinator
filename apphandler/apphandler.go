package apphandler

import (
	"fmt"
	"regexp"

	"github.com/fiveateooate/deployinator/helmbuddy"
	"github.com/fiveateooate/deployinator/k8sbuddy"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
)

func selectVersion(chart string) string {
	var (
		selectedVersion int
	)
	pkgs := helmbuddy.GetPkgs(chart)
	fmt.Println("Select Pkg Version:")
	for idx, version := range pkgs {
		fmt.Printf("  %d. %s\n", idx+1, version)
	}
	fmt.Print("Selection: ")
	fmt.Scanln(&selectedVersion)
	fmt.Printf("Using version %s\n", pkgs[selectedVersion-1])
	return pkgs[selectedVersion-1]
}

func getK8sVersion(deployment *appsv1.Deployment, appName string) string {
	var (
		k8sVersion string
		re         = regexp.MustCompile(fmt.Sprintf(".*%s:(.*)$", appName))
	)
	for _, container := range deployment.Spec.Template.Spec.Containers {
		k8sVersion = re.FindStringSubmatch(container.Image)[1]
		break
	}
	return k8sVersion
}

func checkVersion(deployment *appsv1.Deployment, helmRelease helmbuddy.HelmRelease, version string, appName string) bool {
	k8sVersion := getK8sVersion(deployment, appName)
	fmt.Printf("Checking version\n")
	fmt.Printf("Helm Verison: %s, k8sVersion: %s, selected Version: %s\n", helmRelease.AppVersion, k8sVersion, version)
	if k8sVersion == helmRelease.AppVersion && k8sVersion == version {
		return false
	}
	return true
}

// ManageApp do stuff for a single app
func ManageApp(appName string, namespace string, kubeContext string, helmRepo string, clientset *kubernetes.Clientset) {
	var (
		version string
		chart   = fmt.Sprintf("%s/%s", helmRepo, appName)
	)
	fmt.Printf("Getting info for deployment %s\n", appName)
	deployment, err := k8sbuddy.GetDeployment(appName, namespace, clientset)
	if err == nil {
		fmt.Printf("Found k8s deployment: %s\n", deployment.Name)
	} else {
		fmt.Println(err)
	}
	helmRelease := helmbuddy.GetRelease(appName, namespace, kubeContext)
	if helmRelease.Name != "" {
		fmt.Printf("Found helm release: %s\n", helmRelease.Name)
		if deployment != nil {
			// do something with version checking ?
			fmt.Printf("Upgrading release %s\n", helmRelease.Name)
			version = selectVersion(chart)
			checkVersion(deployment, helmRelease, version, appName)
			helmbuddy.HelmUpgrade(helmRelease.Name, namespace, chart, version, kubeContext)
		}
	} else {
		fmt.Printf("Installing %s\n", appName)
		helmbuddy.GetPkgs(chart)
	}
}
