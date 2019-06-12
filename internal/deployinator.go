package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fiveateooate/deployinator/apphandler"
	"github.com/fiveateooate/deployinator/clusterconfig"
	"github.com/fiveateooate/deployinator/envfilehandler"
	"github.com/fiveateooate/deployinator/helmbuddy"
	"github.com/fiveateooate/deployinator/k8sbuddy"
	"github.com/wsxiaoys/terminal/color"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"k8s.io/client-go/kubernetes"
)

//Deployinator - vars for this thing
type Deployinator struct {
}

func helmDeploy(app *apphandler.App) {
	switch app.K8sApp.Kind {
	case "deployment":
		ah := apphandler.DeploymentHandler{App: app}
		ah.ManageHelmApp()
	case "daemonset":
		ah := apphandler.DaemonsetHandler{App: app}
		ah.ManageHelmApp()
	case "statefulset":
		ah := apphandler.StatefulsetHandler{App: app}
		ah.ManageHelmApp()
	default:
		ah := apphandler.NullHandler{App: app}
		ah.ManageApp()
	}
}

func main() {
	var (
		app           = kingpin.New("deployinator", "Deploy stuff to k8s cluster")
		incluster     = app.Flag("incluster", "Use kuebenetes in cluster config").Default("false").Bool()
		context       = app.Flag("context", "Kube context").Default("local").String()
		deployerType  = app.Flag("deployertype", "Type of deployer: default helm").Default("helm").String()
		appName       = app.Flag("appname", "Name of app to deploy").String()
		namespaceName = app.Flag("namespace", "kubernetes namespace to use").String()
		clusterConfig = app.Flag("clusterconfig", "Path to cluster config file").String()
		helmRepo      = app.Flag("helmrepo", "Name of helm repo").String()
		//helmURL       = app.Flag("helmurl", "URL of helm repo").String()
		helmValues  = app.Flag("helmvalues", "Path to helm values file").String()
		helmVersion = app.Flag("helmversion", "Version to install from helm repo").String()
		onetime     = app.Flag("onetime", "Only run one time").Default("false").Bool()
		envfile     = app.Flag("envfile", "Path to a file containing stuff").Envar("ENVFILE").String()
		clientset   *kubernetes.Clientset
		envVars     envfilehandler.Envfile
	)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *envfile != "" {
		envVars.LoadEnvfile(*envfile)
	} else {
		envVars.LoadFromFlags(*appName, *namespaceName, *helmRepo)
	}

	if *clusterConfig != "" {
		// handle continuous deploy from config
		for {
			cc := clusterconfig.ClusterConfig{}
			cc.ParseClusterConfig(*clusterConfig)
			for _, namespace := range cc.Deployment.Namespaces {
				fmt.Printf("  %s\n", namespace.Name)
				for _, service := range namespace.Services {
					name := strings.Replace(service.Chart, fmt.Sprintf("%s/", *helmRepo), "", -1)
					if service.DeployerType == "" || service.DeployerType == "helm" {
						envVars.LoadFromFlags(name, namespace.Name, *helmRepo)
						k8sApp := k8sbuddy.K8sApp{}
						helmInfo := helmbuddy.HelmInfo{}
						app := apphandler.App{K8sApp: &k8sApp, HelmInfo: &helmInfo, DeployerType: "helm"}
						k8sApp.GetAppInfo(envVars.Slug, envVars.Domain, clientset)
						helmInfo.GetHelmInfo(envVars.Slug, envVars.Domain, envVars.HelmRepo, *helmValues, *context, *helmVersion)
						helmDeploy(&app)
					}
					// fmt.Printf("Slug: %s\n", envVars.Slug)
					// fmt.Printf("    chart: %s\n", service.Chart)
					// if service.Version != "" {
					// 	fmt.Printf("    version: %s\n", service.Version)
					// }
					// if service.DeployerType != "" {
					// 	fmt.Printf("    deployertype: %s\n", service.DeployerType)
					// } else {
					// 	fmt.Printf("    deployertype: %s\n", "helm")
					// }
				}
			}
			// multideploy(config)
			if *onetime {
				break
			}
			color.Printf("@cDone\n")
			time.Sleep(30 * time.Second)
		}
	} else if envVars.Slug != "" {
		clientset = k8sbuddy.Connect(*incluster, *context)
		// choose deployer type
		switch *deployerType {
		case "helm":
			k8sApp := k8sbuddy.K8sApp{}
			helmInfo := helmbuddy.HelmInfo{}
			app := apphandler.App{K8sApp: &k8sApp, HelmInfo: &helmInfo, DeployerType: "helm"}

			k8sApp.GetAppInfo(envVars.Slug, envVars.Domain, clientset)
			helmInfo.GetHelmInfo(envVars.Slug, envVars.Domain, envVars.HelmRepo, *helmValues, *context, *helmVersion)
			helmDeploy(&app)
		case "newawesomedeployer":
			fmt.Println("newawesomedeployer")
		default:
			fmt.Println("unknown deployer")
		}
		color.Printf("@cDone\n")
	} else {
		color.Println("@yEither appname or clusterconfig must be set")
		os.Exit(1)
	}
}
