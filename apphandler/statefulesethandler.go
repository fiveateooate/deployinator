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

// StatefulsetHandler do statefulset specific stuff
type StatefulsetHandler struct {
}

func (ss *StatefulsetHandler) getVersion(statefulset *appsv1.StatefulSet, appName string) string {
	var (
		k8sVersion string
		re         = regexp.MustCompile(fmt.Sprintf(".*%s:(.*)$", appName))
	)
	for _, container := range statefulset.Spec.Template.Spec.Containers {
		k8sVersion = re.FindStringSubmatch(container.Image)[1]
		break
	}
	return k8sVersion
}

// ManageHelmApp do stuff for a single app
func (ss *StatefulsetHandler) ManageHelmApp(helmInfo model.HelmInfo, clientset *kubernetes.Clientset) {
	var (
		version         string
		statefulset     *appsv1.StatefulSet
		deployedVersion string
	)
	helmbuddy.RepoUpdate(helmInfo)
	helmbuddy.GetRelease(&helmInfo)
	statefulset = k8sbuddy.GetStatefulset(helmInfo.AppName, helmInfo.Namespace, clientset)
	if helmInfo.ReleaseExists && statefulset != nil {
		deployedVersion = ss.getVersion(statefulset, helmInfo.AppName)
		version = selectVersion(helmInfo.Chart)
		if !checkVersion(deployedVersion, helmInfo.ReleaseVersion, version) {
			color.Printf("@yVersion %s already running\n", version)
			return
		}
		color.Printf("@yUpgrading release %s\n", helmInfo.ReleaseName)
		helmbuddy.HelmUpgrade(helmInfo, version)
	} else if !helmInfo.ReleaseExists && statefulset == nil {
		version = selectVersion(helmInfo.Chart)
		fmt.Printf("Installing %s\n", helmInfo.Chart)
		helmbuddy.HelmInstall(helmInfo, version)
	} else {
		color.Printf("@rDIE DIE DIE bad helm or k8s state")
		os.Exit(2)
	}
}
