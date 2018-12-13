package apphandler

import (
	"fmt"
	"regexp"

	"github.com/fiveateooate/deployinator/helmbuddy"
	"github.com/fiveateooate/deployinator/k8sbuddy"
	"github.com/fiveateooate/deployinator/types"
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

func checkVersion(k8sVersion string, helmVersion string, version string) bool {
	fmt.Printf("Checking version\n")
	fmt.Printf("Helm Verison: %s, k8sVersion: %s, selected Version: %s\n", helmVersion, k8sVersion, version)
	if k8sVersion == helmVersion && k8sVersion == version {
		return false
	}
	return true
}

// ManageHelmApp do stuff for a single app
func ManageHelmApp(helmInfo types.HelmInfo, clientset *kubernetes.Clientset) {
	var (
		version string
	)
	fmt.Printf("Getting info for deployment %s\n", helmInfo.AppName)
	deployment, err := k8sbuddy.GetDeployment(helmInfo.AppName, helmInfo.Namespace, clientset)
	if err == nil {
		fmt.Printf("Found k8s deployment: %s\n", deployment.Name)
	} else {
		fmt.Println(err)
	}
	helmbuddy.GetRelease(&helmInfo)
	if helmInfo.ReleaseName != "" {
		fmt.Printf("Found helm release: %s\n", helmInfo.ReleaseName)
		if deployment != nil {
			version = selectVersion(helmInfo.Chart)
			checkVersion(getK8sVersion(deployment, helmInfo.AppName), helmInfo.ReleaseVersion, version)
			fmt.Printf("Upgrading release %s\n", helmInfo.ReleaseName)
			helmbuddy.HelmUpgrade(helmInfo, version)
		}
	} else {
		fmt.Printf("Installing %s\n", helmInfo.AppName)
		helmbuddy.GetPkgs(helmInfo.Chart)
	}
}
