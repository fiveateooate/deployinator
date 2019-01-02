package apphandler

import (
	"fmt"
	"os"

	"github.com/fiveateooate/deployinator/helmbuddy"
	"github.com/fiveateooate/deployinator/model"
	"github.com/wsxiaoys/terminal/color"
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
	if len(pkgs) == 0 {
		fmt.Printf("No helm packages found, please check repo")
		os.Exit(0)
	}
	color.Println("@cSelect version:")
	for idx, version := range pkgs {
		fmt.Printf("  %d. %s\n", idx+1, version)
	}
	fmt.Print("select: ")
	fmt.Scanln(&selectedVersion)
	return pkgs[selectedVersion-1]
}

func checkVersion(k8sVersion string, helmVersion string, version string) bool {
	if k8sVersion == helmVersion && k8sVersion == version {
		return false
	}
	return true
}
