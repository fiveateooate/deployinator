package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"weavelab.xyz/ops-deployinator/deployers"
	"weavelab.xyz/ops-deployinator/k8sclient"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func setKubeConfig(kubeconfig *string) *string {
	if *kubeconfig == "" {
		temp := fmt.Sprintf("%s/.kube/config", os.Getenv("HOME"))
		kubeconfig = &temp
	}
	return kubeconfig
}

func prompt() {
	fmt.Printf("-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}

func main() {
	var (
		app          = kingpin.New("deployinator", "Deploy stuff to k8s cluster")
		onetime      = app.Flag("onetime", "Do only once, don't loop").Default("false").Bool()
		dryrun       = app.Flag("dryrun", "Print stuff, don't actually do anything").Default("false").Bool()
		incluster    = app.Flag("incluster", "Use kuebenetes in cluster config").Default("false").Bool()
		context      = app.Flag("context", "Kube context").Default("local").String()
		deployerType = app.Flag("deployertype", "Type of deployer: default helm").Default("helm").String()
	//	kubeconfig   = app.Flag("kubeconfig", "path to kube config").String()
	)
	kingpin.MustParse(app.Parse(os.Args[1:]))
	for {
		if !*dryrun {
			fmt.Println("Actually doing stuff")
			fmt.Printf("Context: %s\n", *context)
			if *deployerType == "helm" {
				deployer := helmdeployer.Deployer()
				fmt.Println(deployer)
			}
			if *incluster {
				fmt.Println("Using incluster config")
			} else {
				fmt.Println("Using external config")
				clientset, err := k8sclient.ExternalClient(*context)
				if err != nil {
					fmt.Printf("Failed to connect to k8s: %s\n", err)
					os.Exit(1)
				}
				fmt.Println(clientset)
				deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
				prompt()
				fmt.Printf("Listing deployments in namespace %q:\n", apiv1.NamespaceDefault)
				list, err := deploymentsClient.List(metav1.ListOptions{})
				if err != nil {
					panic(err)
				}
				for _, d := range list.Items {
					fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
				}
			}
		} else {
			fmt.Println("Would have done xyz")
		}
		if *onetime {
			break
		}
		time.Sleep(30 * time.Second)
	}
}
