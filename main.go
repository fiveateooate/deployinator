package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fiveateooate/deployinator/apphandler"
	"github.com/fiveateooate/deployinator/clusterconfig"
	"github.com/fiveateooate/deployinator/helmbuddy"
	"github.com/fiveateooate/deployinator/k8sbuddy"
	"github.com/wsxiaoys/terminal/color"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"k8s.io/client-go/kubernetes"
)

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
		helmRepo      = app.Flag("helmrepo", "Name of helm repo").Default("weavelabxyz").String()
		// helmURL       = app.Flag("helmurl", "URL of helm repo").Default("https://adsfadsf").String()
		helmValues = app.Flag("helmvalues", "Path to helm values file").String()
		onetime    = app.Flag("onetime", "Only run one time").Default("false").Bool()
		clientset  *kubernetes.Clientset
	)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *appName != "" {
		clientset = k8sbuddy.Connect(*incluster, *context)
		// choose deployer type
		switch *deployerType {
		case "helm":
			k8sApp := k8sbuddy.K8sApp{}
			helmInfo := helmbuddy.HelmInfo{}
			app := apphandler.App{K8sApp: &k8sApp, HelmInfo: &helmInfo, DeployerType: "helm"}

			k8sApp.GetAppInfo(*appName, *namespaceName, clientset)
			helmInfo.GetHelmInfo(*appName, *namespaceName, *helmRepo, *helmValues, *context)
			helmDeploy(&app)
		case "newawesomedeployer":
			fmt.Println("newawesomedeployer")
		default:
			fmt.Println("unknown deployer")
		}
		color.Printf("@cDone\n")
	} else if *clusterConfig != "" {
		cc := clusterconfig.ClusterConfig{}
		cc.ParseClusterConfig(*clusterConfig)
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
}
