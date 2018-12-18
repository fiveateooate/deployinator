package apphandler

import (
	"fmt"
	"os"
	"regexp"

	"github.com/fiveateooate/deployinator/helmbuddy"
	"github.com/fiveateooate/deployinator/k8sbuddy"
	"github.com/fiveateooate/deployinator/model"
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
	fmt.Println("Get version")
	for _, container := range deployment.Spec.Template.Spec.Containers {
		fmt.Println(container)
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
	fmt.Printf("Getting info for %s\n", helmInfo.AppName)
	deployment, err = k8sbuddy.GetDeployment(helmInfo.AppName, helmInfo.Namespace, clientset)
	if err == nil {
		fmt.Printf("Found k8s deployment: %s\n", deployment.Name)
		deployedVersion = dp.getVersion(deployment, helmInfo.AppName)
	} else {
		fmt.Println(err)
	}
	helmbuddy.RepoUpdate(helmInfo)
	helmbuddy.GetRelease(&helmInfo)
	if helmInfo.ReleaseName != "" {
		fmt.Printf("Found helm release: %s\n", helmInfo.ReleaseName)
		if deployment != nil {
			version = selectVersion(helmInfo.Chart)
			if !checkVersion(deployedVersion, helmInfo.ReleaseVersion, version) {
				fmt.Printf("Version %s is already installed\n", version)
				return
			}
			fmt.Printf("Upgrading release %s\n", helmInfo.ReleaseName)
			helmbuddy.HelmUpgrade(helmInfo, version)
		} else {
			fmt.Println("Something is not right DIE DIE DIE")
			os.Exit(2)
		}
	} else {
		fmt.Printf("Installing %s\n", helmInfo.AppName)
		version = selectVersion(helmInfo.Chart)
		fmt.Printf("Installing %s\n", helmInfo.Chart)
		helmbuddy.HelmInstall(helmInfo, version)
	}
}
