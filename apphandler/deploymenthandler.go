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

// DeploymentHandler do stuff specific to deployment types
type DeploymentHandler struct {
}

func (dp *DeploymentHandler) getVersion(deployment *appsv1.Deployment, appName string) string {
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

// ManageHelmApp do stuff for a single app
func (dp *DeploymentHandler) ManageHelmApp(helmInfo model.HelmInfo, clientset *kubernetes.Clientset) {
	var (
		version         string
		deployment      *appsv1.Deployment
		deployedVersion string
	)
	helmbuddy.RepoUpdate(helmInfo)
	helmbuddy.GetRelease(&helmInfo)
	deployment = k8sbuddy.GetDeployment(helmInfo.AppName, helmInfo.Namespace, clientset)
	if helmInfo.ReleaseExists && deployment != nil {
		deployedVersion = dp.getVersion(deployment, helmInfo.AppName)
		version = selectVersion(helmInfo.Chart)
		if !checkVersion(deployedVersion, helmInfo.ReleaseVersion, version) {
			color.Printf("@yVersion %s already running\n", version)
			return
		}
		color.Printf("@yUpgrading release %s\n", helmInfo.ReleaseName)
		helmbuddy.HelmUpgrade(helmInfo, version)
	} else if !helmInfo.ReleaseExists && deployment == nil {
		version = selectVersion(helmInfo.Chart)
		fmt.Printf("Installing %s\n", helmInfo.Chart)
		helmbuddy.HelmInstall(helmInfo, version)
	} else {
		color.Printf("@rDIE DIE DIE bad helm or k8s state")
		os.Exit(2)
	}
}
