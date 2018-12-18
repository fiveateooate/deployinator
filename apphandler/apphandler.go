package apphandler

import (
	"fmt"

	"github.com/fiveateooate/deployinator/helmbuddy"
	"github.com/fiveateooate/deployinator/model"
	"k8s.io/client-go/kubernetes"
)

// AppHandler commect to get rid of lint
type AppHandler interface {
	ManageHelmApp(helmInfo model.HelmInfo, clientset *kubernetes.Clientset)
}

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

func checkVersion(k8sVersion string, helmVersion string, version string) bool {
	fmt.Printf("Checking version\n")
	fmt.Printf("Helm Verison: %s, k8sVersion: %s, selected Version: %s\n", helmVersion, k8sVersion, version)
	if k8sVersion == helmVersion && k8sVersion == version {
		return false
	}
	return true
}
