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

// StatefulsetHandler do daemonset specific stuff
type StatefulsetHandler struct {
}

func (ds *StatefulsetHandler) getVersion(statefulset *appsv1.StatefulSet, appName string) string {
	var (
		k8sVersion string
		re         = regexp.MustCompile(fmt.Sprintf(".*%s:(.*)$", appName))
	)
	for _, container := range statefulset.Spec.Template.Spec.Containers {
		k8sVersion = re.FindStringSubmatch(container.Image)[1]
		break
	}
	fmt.Println(k8sVersion)
	return k8sVersion
}

// ManageHelmApp do stuff for a single app
func (ds *StatefulsetHandler) ManageHelmApp(helmInfo model.HelmInfo, clientset *kubernetes.Clientset) {
	var (
		version         string
		statefulset     *appsv1.StatefulSet
		err             error
		deployedVersion string
	)
	fmt.Printf("Getting info for %s\n", helmInfo.AppName)
	statefulset, err = k8sbuddy.GetStatefulset(helmInfo.AppName, helmInfo.Namespace, clientset)
	if err == nil {
		fmt.Printf("Found k8s statefulset: %s\n", statefulset.Name)
		deployedVersion = ds.getVersion(statefulset, helmInfo.AppName)
	} else {
		fmt.Println(err)
	}
	helmbuddy.RepoUpdate(helmInfo)
	helmbuddy.GetRelease(&helmInfo)
	if helmInfo.ReleaseExists {
		fmt.Printf("Found helm release: %s\n", helmInfo.ReleaseName)
		if statefulset != nil {
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
