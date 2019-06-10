package apphandler

import (
	"fmt"
	"os"
	"regexp"

	"github.com/fiveateooate/deployinator/helmbuddy"
	"github.com/wsxiaoys/terminal/color"
	appsv1 "k8s.io/api/apps/v1"
)

// DeploymentHandler do stuff specific to deployment types
type DeploymentHandler struct {
	Handler AppHandler
	App     *App
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

//ManageApp do stuff
func (dp *DeploymentHandler) ManageApp() {
	fmt.Println(dp.App)
}

// ManageHelmApp do stuff for a single app
func (dp *DeploymentHandler) ManageHelmApp() {
	var (
		version         string
		deployedVersion string
	)
	dp.App.HelmInfo.RepoUpdate()
	// helmbuddy.GetRelease(dp.App.HelmInfo)
	if dp.App.HelmInfo.ReleaseExists && dp.App.K8sApp.DP != nil {
		deployedVersion = dp.getVersion(dp.App.K8sApp.DP, dp.App.HelmInfo.AppName)
		version = dp.App.HelmInfo.Version
		if version == "" {
			version = selectVersion(dp.App.HelmInfo.Chart)
		}
		if !checkVersion(deployedVersion, dp.App.HelmInfo.ReleaseVersion, version) {
			color.Printf("@yVersion %s already running\n", version)
			return
		}
		color.Printf("@yUpgrading release %s\n", dp.App.HelmInfo.ReleaseName)
		dp.App.HelmInfo.HelmUpgrade(version)
	} else if !dp.App.HelmInfo.ReleaseExists && dp.App.K8sApp.DP == nil {
		version = selectVersion(dp.App.HelmInfo.Chart)
		fmt.Printf("Installing %s\n", dp.App.HelmInfo.Chart)
		helmbuddy.HelmInstall(dp.App.HelmInfo, version)
	} else {
		color.Printf("@rDIE DIE DIE bad helm or k8s state")
		os.Exit(2)
	}
}
