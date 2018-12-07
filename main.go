package main

import (
	"fmt"
	"os"
	"time"

	"weavelab.xyz/ops-deployinator/deployers"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		app          = kingpin.New("deployinator", "Deploy stuff to k8s cluster")
		onetime      = app.Flag("onetime", "Do only once, don't loop").Default("false").Bool()
		dryrun       = app.Flag("dryrun", "Print stuff, don't actually do anything").Default("false").Bool()
		incluster    = app.Flag("incluster", "Use kuebenetes in cluster config").Default("false").Bool()
		context      = app.Flag("context", "Kube context").Default("local").String()
		deployerType = app.Flag("deployertype", "Type of deployer: default helm").Default("helm").String()
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
