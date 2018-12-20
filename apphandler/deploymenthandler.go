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
		err             error
		deployedVersion string
	)
	color.Printf("@cSearching for deployment ...")
	deployment, err = k8sbuddy.GetDeployment(helmInfo.AppName, helmInfo.Namespace, clientset)
	if err == nil {
		color.Printf("found %s\n", deployment.Name)
		deployedVersion = dp.getVersion(deployment, helmInfo.AppName)
	} else {
		fmt.Println(err)
	}
	helmbuddy.RepoUpdate(helmInfo)
	helmbuddy.GetRelease(&helmInfo)
	if helmInfo.ReleaseName != "" {
		if deployment != nil {
			version = selectVersion(helmInfo.Chart)
			if !checkVersion(deployedVersion, helmInfo.ReleaseVersion, version) {
				fmt.Printf("Version %s is already installed\n", version)
				return
			}
			color.Printf("@yUpgrading release %s\n", helmInfo.ReleaseName)
			helmbuddy.HelmUpgrade(helmInfo, version)
		} else {
			color.Println("@rSomething is not right DIE DIE DIE")
			os.Exit(2)
		}
	} else {
		helmInfo.ReleaseName = fmt.Sprintf("%s-%s", helmInfo.AppName, RandStringBytes(5))
		color.Printf("@yInstalling %s\n", helmInfo.AppName)
		version = selectVersion(helmInfo.Chart)
		fmt.Printf("Installing %s\n", helmInfo.Chart)
		helmbuddy.HelmInstall(helmInfo, version)
	}
}
