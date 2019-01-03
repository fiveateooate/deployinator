package apphandler

import (
	"fmt"
	"os"
	"regexp"

	"github.com/fiveateooate/deployinator/helmbuddy"
	"github.com/fiveateooate/deployinator/k8sbuddy"
	"github.com/fiveateooate/deployinator/model"
	"github.com/wsxiaoys/terminal/color"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
)

// DaemonsetHandler do daemonset specific stuff
type DaemonsetHandler struct {
}

func (ds *DaemonsetHandler) getVersion(daemonset *appsv1.DaemonSet, appName string) string {
	var (
		k8sVersion string
		re         = regexp.MustCompile(fmt.Sprintf(".*%s:(.*)$", appName))
	)
	for _, container := range daemonset.Spec.Template.Spec.Containers {
		k8sVersion = re.FindStringSubmatch(container.Image)[1]
		break
	}
	return k8sVersion
}

// ManageHelmApp do stuff for a single app
func (ds *DaemonsetHandler) ManageHelmApp(helmInfo model.HelmInfo, clientset *kubernetes.Clientset) {
	var (
		version         string
		daemonset       *appsv1.DaemonSet
		deployedVersion string
	)
	helmbuddy.RepoUpdate(helmInfo)
	helmbuddy.GetRelease(&helmInfo)
	daemonset = k8sbuddy.GetDaemonset(helmInfo.AppName, helmInfo.Namespace, clientset)
	if helmInfo.ReleaseExists && daemonset != nil {
		deployedVersion = ds.getVersion(daemonset, helmInfo.AppName)
		version = selectVersion(helmInfo.Chart)
		if !checkVersion(deployedVersion, helmInfo.ReleaseVersion, version) {
			color.Printf("@yVersion %s already running\n", version)
			return
		}
		color.Printf("@yUpgrading release %s\n", helmInfo.ReleaseName)
		helmbuddy.HelmUpgrade(helmInfo, version)
	} else if !helmInfo.ReleaseExists && daemonset == nil {
		version = selectVersion(helmInfo.Chart)
		fmt.Printf("Installing %s\n", helmInfo.Chart)
		helmbuddy.HelmInstall(helmInfo, version)
	} else {
		color.Printf("@rDIE DIE DIE bad helm or k8s state")
		os.Exit(2)
	}
}
