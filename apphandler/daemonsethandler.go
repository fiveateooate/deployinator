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
	fmt.Println(k8sVersion)
	return k8sVersion
}

// ManageHelmApp do stuff for a single app
func (ds *DaemonsetHandler) ManageHelmApp(helmInfo model.HelmInfo, clientset *kubernetes.Clientset) {
	var (
		version         string
		daemonset       *appsv1.DaemonSet
		err             error
		deployedVersion string
	)
	fmt.Printf("Getting info for %s\n", helmInfo.AppName)
	daemonset, err = k8sbuddy.GetDaemonset(helmInfo.AppName, helmInfo.Namespace, clientset)
	if err == nil {
		fmt.Printf("Found k8s daemonset: %s\n", daemonset.Name)
		deployedVersion = ds.getVersion(daemonset, helmInfo.AppName)
	} else {
		fmt.Println(err)
	}
	helmbuddy.RepoUpdate(helmInfo)
	helmbuddy.GetRelease(&helmInfo)
	if helmInfo.ReleaseExists {
		fmt.Printf("Found helm release: %s\n", helmInfo.ReleaseName)
		if daemonset != nil {
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
