package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fiveateooate/deployinator/apphandler"
	"github.com/fiveateooate/deployinator/k8sclient"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func setKubeConfig(kubeconfig *string) *string {
	if *kubeconfig == "" {
		temp := fmt.Sprintf("%s/.kube/config", os.Getenv("HOME"))
		kubeconfig = &temp
	}
	return kubeconfig
}

func main() {
	var (
		app           = kingpin.New("deployinator", "Deploy stuff to k8s cluster")
		onetime       = app.Flag("onetime", "Do only once, don't loop").Default("false").Bool()
		dryrun        = app.Flag("dryrun", "Print stuff, don't actually do anything").Default("false").Bool()
		incluster     = app.Flag("incluster", "Use kuebenetes in cluster config").Default("false").Bool()
		context       = app.Flag("context", "Kube context").Default("local").String()
		deployerType  = app.Flag("deployertype", "Type of deployer: default helm").Default("helm").String()
		appName       = app.Flag("appname", "Name of app to deploy").String()
		namespaceName = app.Flag("namespace", "kubernetes namespace to use").String()
		clusterConfig = app.Flag("clusterconfig", "Path to cluster config file").String()
	//	kubeconfig   = app.Flag("kubeconfig", "path to kube config").String()
	)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *appName == "" && *clusterConfig == "" {
		fmt.Println("Either appname or clusterconfig must be set")
		os.Exit(1)
	}

	for {
		if !*dryrun {
			fmt.Printf("Using Context: %s\n", *context)
			// choose deployer type
			// switch *deployerType {
			// case "helm":
			// 	deployer := helmdeployer.Deployer()
			// case "newawesomedeployer":
			// 	deployer := newawesomedeployer.Deployer()
			// }
			fmt.Printf("Deployer: %s\n", *deployerType)
			if *incluster {
				fmt.Println("Using incluster config")
			} else {
				fmt.Println("Using external kubeconfig")
				clientset, err := k8sclient.ExternalClient(*context)
				if err != nil {
					fmt.Printf("Failed to connect to k8s: %s\n", err)
					os.Exit(1)
				}
				if *appName != "" {
					apphandler.ManageApp(*appName, *namespaceName, *context, clientset)
				}
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
