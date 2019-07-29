package apphandler

// import (
// 	"fmt"
// 	"os"
// 	"regexp"

// 	"github.com/fiveateooate/deployinator/internal/helmbuddy"
// 	"github.com/wsxiaoys/terminal/color"
// 	appsv1 "k8s.io/api/apps/v1"
// )

// // DaemonsetHandler do daemonset specific stuff
// type DaemonsetHandler struct {
// 	Handler AppHandler
// 	App     *App
// }

// func (ds *DaemonsetHandler) getVersion(daemonset *appsv1.DaemonSet, appName string) string {
// 	var (
// 		k8sVersion string
// 		re         = regexp.MustCompile(fmt.Sprintf(".*%s:(.*)$", appName))
// 	)
// 	for _, container := range daemonset.Spec.Template.Spec.Containers {
// 		k8sVersion = re.FindStringSubmatch(container.Image)[1]
// 		break
// 	}
// 	return k8sVersion
// }

// //ManageApp do stuff
// func (ds *DaemonsetHandler) ManageApp() {
// 	fmt.Println("dsh")
// }

// // ManageHelmApp do stuff for a single app
// func (ds *DaemonsetHandler) ManageHelmApp() {
// 	var (
// 		version         string
// 		deployedVersion string
// 	)
// 	ds.App.HelmInfo.RepoUpdate()
// 	if ds.App.HelmInfo.ReleaseExists && ds.App.K8sApp.DS != nil {
// 		deployedVersion = ds.getVersion(ds.App.K8sApp.DS, ds.App.HelmInfo.AppName)
// 		version = ds.App.HelmInfo.Version
// 		if version == "" {
// 			version = selectVersion(ds.App.HelmInfo.Chart)
// 		}
// 		if !checkVersion(deployedVersion, ds.App.HelmInfo.ReleaseVersion, version) {
// 			color.Printf("@yVersion %s already running\n", version)
// 			return
// 		}
// 		color.Printf("@yUpgrading release %s\n", ds.App.HelmInfo.ReleaseName)
// 		ds.App.HelmInfo.HelmUpgrade(version)
// 	} else if !ds.App.HelmInfo.ReleaseExists && ds.App.K8sApp.DS == nil {
// 		version = selectVersion(ds.App.HelmInfo.Chart)
// 		fmt.Printf("Installing %s\n", ds.App.HelmInfo.Chart)
// 		helmbuddy.HelmInstall(ds.App.HelmInfo, version)
// 	} else {
// 		color.Printf("@rDIE DIE DIE bad helm or k8s state")
// 		os.Exit(2)
// 	}
// }
