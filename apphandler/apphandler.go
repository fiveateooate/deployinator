package apphandler

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/fiveateooate/deployinator/helmbuddy"
	"github.com/fiveateooate/deployinator/model"
	"github.com/wsxiaoys/terminal/color"
	"k8s.io/client-go/kubernetes"
)

// AppHandler commect to get rid of lint
type AppHandler interface {
	ManageHelmApp(helmInfo model.HelmInfo, clientset *kubernetes.Clientset)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz01234566789"

// RandStringBytes copied from https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func selectVersion(chart string) string {
	var (
		selectedVersion int
	)
	pkgs := helmbuddy.GetPkgs(chart)
	if len(pkgs) == 0 {
		fmt.Printf("No helm packages found, please check repo")
		os.Exit(0)
	}
	color.Println("@cSelect version:")
	for idx, version := range pkgs {
		fmt.Printf("  %d. %s\n", idx+1, version)
	}
	fmt.Print("select: ")
	fmt.Scanln(&selectedVersion)
	return pkgs[selectedVersion-1]
}

func checkVersion(k8sVersion string, helmVersion string, version string) bool {
	if k8sVersion == helmVersion && k8sVersion == version {
		return false
	}
	return true
}
