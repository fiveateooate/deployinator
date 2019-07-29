package apphandler

// import (
// 	"fmt"
// 	"os"
// 	"regexp"

// 	"github.com/fiveateooate/deployinator/internal/helmbuddy"
// 	"github.com/wsxiaoys/terminal/color"
// 	appsv1 "k8s.io/api/apps/v1"
// )

// // StatefulsetHandler do statefulset specific stuff
// type StatefulsetHandler struct {
// 	Handler AppHandler
// 	App     *App
// }

// func (ss *StatefulsetHandler) getVersion(statefulset *appsv1.StatefulSet, appName string) string {
// 	var (
// 		k8sVersion string
// 		re         = regexp.MustCompile(fmt.Sprintf(".*%s:(.*)$", appName))
// 	)
// 	for _, container := range statefulset.Spec.Template.Spec.Containers {
// 		k8sVersion = re.FindStringSubmatch(container.Image)[1]
// 		break
// 	}
// 	return k8sVersion
// }

// //ManageApp manage an app
// func (ss *StatefulsetHandler) ManageApp() {
// 	fmt.Println(ss.App)
// }

// // ManageHelmApp do stuff for a single app
// func (ss *StatefulsetHandler) ManageHelmApp() {
// 	var (
// 		version         string
// 		deployedVersion string
// 	)
// 	ss.App.HelmInfo.RepoUpdate()
// 	if ss.App.HelmInfo.ReleaseExists && ss.App.K8sApp.SS != nil {
// 		deployedVersion = ss.getVersion(ss.App.K8sApp.SS, ss.App.HelmInfo.AppName)
// 		version = ss.App.HelmInfo.Version
// 		if version == "" {
// 			version = selectVersion(ss.App.HelmInfo.Chart)
// 		}
// 		if !checkVersion(deployedVersion, ss.App.HelmInfo.ReleaseVersion, version) {
// 			color.Printf("@yVersion %s already running\n", version)
// 			return
// 		}
// 		color.Printf("@yUpgrading release %s\n", ss.App.HelmInfo.ReleaseName)
// 		ss.App.HelmInfo.HelmUpgrade(version)
// 	} else if !ss.App.HelmInfo.ReleaseExists && ss.App.K8sApp.SS == nil {
// 		version = selectVersion(ss.App.HelmInfo.Chart)
// 		fmt.Printf("Installing %s\n", ss.App.HelmInfo.Chart)
// 		helmbuddy.HelmInstall(ss.App.HelmInfo, version)
// 	} else {
// 		color.Printf("@rDIE DIE DIE bad helm or k8s state")
// 		os.Exit(2)
// 	}
// }
