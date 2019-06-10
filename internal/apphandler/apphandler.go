package apphandler

import (
	"fmt"
	"os"

	"github.com/fiveateooate/deployinator/k8sbuddy"

	"github.com/fiveateooate/deployinator/helmbuddy"
	"github.com/wsxiaoys/terminal/color"
)

//App holds info for an app
type App struct {
	HelmInfo     *helmbuddy.HelmInfo
	K8sApp       *k8sbuddy.K8sApp
	DeployerType string
}

// AppHandler commect to get rid of lint
type AppHandler interface {
	ManageApp()
}

// AbstractHandler - just for whatever
type AbstractHandler struct {
	Handler AppHandler
	App     *App
}

// // HandleApp do stuff with an app
// func (ah *AbstractHandler) HandleApp() {
// 	fmt.Println("HandleApp")
// 	switch ah.App.DeployerType {
// 	case "helm":
// 		if ah.App.K8sApp.DP != nil {
// 			DeploymentHandler{}.Handler.ManageApp()
// 		} else if ah.App.K8sApp.DS != nil {
// 			DaemonsetHandler{}.Handler.ManageApp()
// 		} else if ah.App.K8sApp.SS != nil {
// 			StatefulsetHandler{}.Handler.ManageApp()
// 		} else {
// 			NullHandler{}.Handler.ManageApp()
// 		}
// 	}
// }

func selectVersion(chart string) string {
	var (
		selectedVersion int
		pkgCount        int
	)
	pkgs := helmbuddy.GetPkgs(chart)
	pkgCount = len(pkgs)
	if pkgCount < 1 {
		os.Exit(1)
	}
	color.Println("@cSelect version:")
	for idx, version := range pkgs {
		fmt.Printf("  %d. %s\n", idx+1, version)
	}
	fmt.Print("select: ")
	fmt.Scanln(&selectedVersion)
	if selectedVersion-1 >= pkgCount || selectedVersion <= 0 {
		color.Printf("@rInvalid selection\n")
		os.Exit(1)
	}
	return pkgs[selectedVersion-1]
}

func checkVersion(k8sVersion string, helmVersion string, version string) bool {
	if k8sVersion == helmVersion && k8sVersion == version {
		return false
	}
	return true
}
