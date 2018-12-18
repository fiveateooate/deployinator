package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fiveateooate/deployinator/apphandler"
	"github.com/fiveateooate/deployinator/helmbuddy"
	"github.com/fiveateooate/deployinator/k8sclient"
	"github.com/fiveateooate/deployinator/model"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"k8s.io/client-go/kubernetes"
)

func connectK8s(incluster bool, context string) *kubernetes.Clientset {
	var (
		clientset *kubernetes.Clientset
		err       error
	)
	if incluster {
		fmt.Println("Using incluster config")
	} else {
		fmt.Println("Using external kubeconfig")
		clientset, err = k8sclient.ExternalClient(context)
		if err != nil {
			fmt.Printf("Failed to connect to k8s: %s\n", err)
			os.Exit(1)
		}
	}
	return clientset
}

func setKubeConfig(kubeconfig *string) *string {
	if *kubeconfig == "" {
		temp := fmt.Sprintf("%s/.kube/config", os.Getenv("HOME"))
		kubeconfig = &temp
	}
	return kubeconfig
}

func main() {
	var (
		app            = kingpin.New("deployinator", "Deploy stuff to k8s cluster")
		onetime        = app.Flag("onetime", "Do only once, don't loop").Default("false").Bool()
		dryrun         = app.Flag("dryrun", "Print stuff, don't actually do anything").Default("false").Bool()
		incluster      = app.Flag("incluster", "Use kuebenetes in cluster config").Default("false").Bool()
		context        = app.Flag("context", "Kube context").Default("local").String()
		deployerType   = app.Flag("deployertype", "Type of deployer: default helm").Default("helm").String()
		appName        = app.Flag("appname", "Name of app to deploy").String()
		namespaceName  = app.Flag("namespace", "kubernetes namespace to use").String()
		clusterConfig  = app.Flag("clusterconfig", "Path to cluster config file").String()
		helmRepo       = app.Flag("helmrepo", "Name of helm repo").Default("weavelabxyz").String()
		helmURL        = app.Flag("helmurl", "URL of helm repo").Default("https://adsfadsf").String()
		helmValuesFile = app.Flag("helmvaluesfile", "Path to values file").String()
		deploymentType = app.Flag("deploymenttype", "Deployment Type [deployment, daemonset, statefulset]").Default("deployment").String()
		clientset      *kubernetes.Clientset
	)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *appName == "" && *clusterConfig == "" {
		fmt.Println("Either appname or clusterconfig must be set")
		os.Exit(1)
	}
	for {
		if !*dryrun {
			clientset = connectK8s(*incluster, *context)
			// choose deployer type
			switch *deployerType {
			case "helm":
				if err := helmbuddy.CheckHelmSetup(*helmRepo, *helmURL); err != nil {
					fmt.Printf("Helm setup incomplete: %s\n", err)
					os.Exit(1)
				}
				chart := fmt.Sprintf("%s/%s", *helmRepo, *appName)
				helmInfo := model.HelmInfo{ValuesFile: *helmValuesFile, Repo: *helmRepo, AppName: *appName, Namespace: *namespaceName, KubeContext: *context, Chart: chart}
				if *appName != "" {
					switch *deploymentType {
					case "deployment":
						apphandler := apphandler.DeploymentHandler{}
						apphandler.ManageHelmApp(helmInfo, clientset)
					case "daemonset":
						apphandler := apphandler.DaemonsetHandler{}
						apphandler.ManageHelmApp(helmInfo, clientset)
					case "statefulset":
						apphandler := apphandler.StatefulsetHandler{}
						apphandler.ManageHelmApp(helmInfo, clientset)
					}
				}
			case "newawesomedeployer":
				// deployer := newawesomedeployer.Deployer()
			}
		} else {
			fmt.Println("Would have done xyz")
		}
		fmt.Printf("Done\n")
		if *onetime {
			break
		}
		time.Sleep(30 * time.Second)
	}
}
