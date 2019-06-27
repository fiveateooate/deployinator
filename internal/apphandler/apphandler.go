package apphandler

// import (
// 	"fmt"
// 	"os"

// 	"github.com/fiveateooate/deployinator/internal/helmbuddy"
// 	"github.com/fiveateooate/deployinator/internal/k8sbuddy"
// 	"github.com/wsxiaoys/terminal/color"
// )

// //App holds info for an app
// type App struct {
// 	HelmInfo     *helmbuddy.HelmInfo
// 	K8sApp       *k8sbuddy.K8sApp
// 	DeployerType string
// }

// // AppHandler comment to get rid of lint
// type AppHandler interface {
// 	ManageApp()
// }

// // AbstractHandler - just for whatever
// type AbstractHandler struct {
// 	Handler AppHandler
// 	App     *App
// }

// func selectVersion(chart string) string {
// 	var (
// 		selectedVersion int
// 		pkgCount        int
// 	)
// 	pkgs := helmbuddy.GetPkgs(chart)
// 	pkgCount = len(pkgs)
// 	if pkgCount < 1 {
// 		os.Exit(1)
// 	}
// 	color.Println("@cSelect version:")
// 	for idx, version := range pkgs {
// 		fmt.Printf("  %d. %s\n", idx+1, version)
// 	}
// 	fmt.Print("select: ")
// 	fmt.Scanln(&selectedVersion)
// 	if selectedVersion-1 >= pkgCount || selectedVersion <= 0 {
// 		color.Printf("@rInvalid selection\n")
// 		os.Exit(1)
// 	}
// 	return pkgs[selectedVersion-1]
// }

// func checkVersion(k8sVersion string, helmVersion string, version string) bool {
// 	if k8sVersion == helmVersion && k8sVersion == version {
// 		return false
// 	}
// 	return true
// }
