package helmbuddy

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
)

// 	helm "k8s.io/helm/pkg/helm"

// *****************************************************************************
//
// maybe someone could get the actual client working
//  for now just going to jump to shell
//
// WIP: make client work
// *****************************************************************************

//
// var (
// 	client = helm.NewClient()
// )
//
// // ListReleases return list of releases
// func ListReleases() {
// 	resp, _ := client.ListReleases()
// 	fmt.Println(resp)
// }

// HelmRelease hold info about a helm release
type HelmRelease struct {
	Name       string `json:"Name"`
	Revision   int    `json:"Revision"`
	Updated    string `json:"Updated"`
	Status     string `json:"Status"`
	Chart      string `json:"Chart"`
	AppVersion string `json:"AppVersion"`
	Namespace  string `json:"Namespace"`
}

// HelmListOutput struct to hold unmarshaled json
type HelmListOutput struct {
	Next     string        `json:"Next"`
	Releases []HelmRelease `json:"Releases"`
}

// ListReleases return struct of releases
func ListReleases(namespace string, kubeContext string) (HelmListOutput, error) {
	var (
		cmdOut []byte
		err    error
		output HelmListOutput
	)
	cmd := "helm"
	args := []string{"--kube-context", kubeContext, "--namespace", namespace, "--output", "json", "list"}
	if cmdOut, err = exec.Command(cmd, args...).Output(); err != nil {
		return output, err
	}
	json.Unmarshal(cmdOut, &output)
	return output, nil
}

// GetRelease return something release
func GetRelease(appname string, namespace string, kubeContext string) HelmRelease {
	var retval HelmRelease
	r, _ := regexp.Compile(fmt.Sprintf("^%s.*", appname))
	releases, _ := ListReleases(namespace, kubeContext)
	for _, release := range releases.Releases {
		match := r.MatchString(release.Name)
		if match {
			retval = release
			break
		}
	}
	return retval
}
