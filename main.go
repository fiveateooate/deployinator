package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fiveateooate/deployinator/apphandler"
	"github.com/fiveateooate/deployinator/helmbuddy"
	"github.com/fiveateooate/deployinator/k8sbuddy"
	"github.com/wsxiaoys/terminal/color"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"k8s.io/client-go/kubernetes"
)

// func setKubeConfig(kubeconfig *string) *string {
// 	if *kubeconfig == "" {
// 		temp := fmt.Sprintf("%s/.kube/config", os.Getenv("HOME"))
// 		kubeconfig = &temp
// 	}
// 	return kubeconfig
// }

func main() {
	var (
		app           = kingpin.New("deployinator", "Deploy stuff to k8s cluster")
		incluster     = app.Flag("incluster", "Use kuebenetes in cluster config").Default("false").Bool()
		context       = app.Flag("context", "Kube context").Default("local").String()
		deployerType  = app.Flag("deployertype", "Type of deployer: default helm").Default("helm").String()
		appName       = app.Flag("appname", "Name of app to deploy").String()
		namespaceName = app.Flag("namespace", "kubernetes namespace to use").String()
		clusterConfig = app.Flag("clusterconfig", "Path to cluster config file").String()
		helmRepo      = app.Flag("helmrepo", "Name of helm repo").Default("weavelabxyz").String()
		// helmURL       = app.Flag("helmurl", "URL of helm repo").Default("https://adsfadsf").String()
		helmValues = app.Flag("helmvalues", "Path to helm values file").String()
		onetime    = app.Flag("onetime", "Only run one time").Default("false").Bool()
		clientset  *kubernetes.Clientset
	)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *appName != "" {
		// where should k8s stuff go?
		clientset = k8sbuddy.Connect(*incluster, *context)
		// choose deployer type
		switch *deployerType {
		case "helm":
			k8sApp := k8sbuddy.K8sApp{}
			helmInfo := helmbuddy.HelmInfo{}
			app := apphandler.App{K8sApp: &k8sApp, HelmInfo: &helmInfo, DeployerType: "helm"}

			k8sApp.GetAppInfo(*appName, *namespaceName, clientset)
			helmInfo.GetHelmInfo(*appName, *namespaceName, *helmRepo, *helmValues, *context)
			switch k8sApp.Kind {
			case "deployment":
				ah := apphandler.DeploymentHandler{App: &app}
				ah.ManageHelmApp()
			case "daemonset":
				ah := apphandler.DaemonsetHandler{App: &app}
				ah.ManageHelmApp()
			case "statefulset":
				ah := apphandler.StatefulsetHandler{App: &app}
				ah.ManageHelmApp()
			default:
				ah := apphandler.NullHandler{App: &app}
				ah.ManageApp()
			}
		case "newawesomedeployer":
			fmt.Println("newawesomedeployer")
		default:
			fmt.Println("unknown deployer")
		}
		color.Printf("@cDone\n")
	} else if *clusterConfig != "" {
		// handle continuous deploy from config
		for {
			// parseconfig
			// multideploy(config)
			if *onetime {
				break
			}
			color.Printf("@cDone\n")
			time.Sleep(30 * time.Second)
		}
	} else {
		color.Println("@yEither appname or clusterconfig must be set")
		os.Exit(1)
	}
	// where should k8s stuff go?
	// clientset = connectK8s(*incluster, *context)
	// ah := apphandler.AbstractHandler{}
	// app := apphandler.App{}
	// k8sApp := k8sbuddy.K8sApp{}
	// // choose deployer type
	// switch *deployerType {
	// case "helm":
	// 	helmInfo := helmbuddy.HelmInfo{}
	// 	// if err := helmbuddy.CheckHelmSetup(*helmRepo, *helmURL); err != nil {
	// 	// 	color.Printf("@rHelm setup incomplete: %s\n", err)
	// 	// 	os.Exit(1)
	// 	// }
	// 	// if *helmValues != "" && !sharedfuncs.FileExists(*helmValues) {
	// 	// 	os.Exit(2)
	// 	// }
	// 	// app.DeployerType = "helm"
	// 	k8sApp.GetAppInfo(*appName, *namespaceName, clientset)
	// 	helmInfo.GetHelmInfo(*appName, *namespaceName, *helmRepo, *helmValues, *context)
	// 	helmbuddy.HelmHandler(helmInfo, k8sApp)
	// 	// app.K8sApp = &k8sApp
	// 	// app.HelmInfo = &helmInfo
	// 	// ah.App = &app
	// 	// // if *appName != "" {
	// 	// // 	app.HandleApp()
	// 	// switch *k8sapp.Kind {
	// 	// case "deployment":
	// 	// 	apphandler := apphandler.DeploymentHandler{}
	// 	// case "daemonset":
	// 	// 	apphandler := apphandler.DaemonsetHandler{}
	// 	// case "statefulset":
	// 	// 	apphandler := apphandler.StatefulsetHandler{}
	// 	// default:
	// 	// 	apphandler := apphandler.NullHander{}
	// 	// }
	// 	// apphandler.ManageHelmApp(helmInfo, clientset)
	// 	// // }
	// case "newawesomedeployer":
	// 	app.DeployerType = "newawesomedeployer"
	// default:
	// 	app.DeployerType = ""
	// }
	// ah.HandleApp()
	// color.Printf("@cDone\n")
	// time.Sleep(30 * time.Second)
}
