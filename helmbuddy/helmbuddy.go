package helmbuddy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/fiveateooate/deployinator/model"
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

func runCmd(cmd string, args []string) ([]byte, error) {
	var (
		cmdOut []byte
		err    error
	)
	if cmdOut, err = exec.Command(cmd, args...).Output(); err != nil {
		return cmdOut, err
	}
	return cmdOut, nil
}

func processLine(text string) string {
	headerRE := regexp.MustCompile("^NAME.*")
	if headerRE.MatchString(text) {
		return ""
	}
	tmp := strings.Fields(text)
	return tmp[1]
}

func checkHelmRepo(helmrepo string) bool {
	var (
		cmdName = "helm"
		cmdArgs = []string{"repo", "list", helmrepo}
		cmdOut  []byte
		err     error
	)
	if cmdOut, err = runCmd(cmdName, cmdArgs); err != nil {
		fmt.Printf("error: %s", err)
		return false
	}
	for _, line := range strings.Split(string(cmdOut), "\n") {
		tmp := strings.Fields(line)
		if len(tmp) > 0 && tmp[0] == helmrepo {
			return true
		}
	}
	return false
}

func addHelmRepo(helmrepo string, helmURL string) {
	var (
		cmdName = "helm"
		cmdArgs = []string{"repo", "add", helmrepo, helmURL}
		cmdOut  []byte
		err     error
	)
	if cmdOut, err = runCmd(cmdName, cmdArgs); err != nil {
		fmt.Printf("error: %s", err)
	} else {
		fmt.Printf("%s", string(cmdOut))
	}
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
	if cmdOut, err = runCmd(cmd, args); err != nil {
		return output, err
	}
	json.Unmarshal(cmdOut, &output)
	return output, nil
}

// GetRelease return something release
func GetRelease(helmInfo *model.HelmInfo) {
	r, _ := regexp.Compile(fmt.Sprintf("^%s.*", helmInfo.AppName))
	releases, _ := ListReleases(helmInfo.Namespace, helmInfo.KubeContext)
	for _, release := range releases.Releases {
		match := r.MatchString(release.Name)
		if match {
			helmInfo.ReleaseName = release.Name
			helmInfo.ReleaseVersion = release.AppVersion
			break
		}
	}
}

// GetPkgs return list of packages known to helm
func GetPkgs(chart string) []string {
	var (
		cmdReader io.Reader
		err       error
		scanner   *bufio.Scanner
		cmdName   = "helm"
		cmdArgs   = []string{"search", "-l", chart}
		versions  []string
	)
	fmt.Printf("Searching for %s\n", chart)
	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err = cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}
	scanner = bufio.NewScanner(cmdReader)
	go func(versions *[]string) {
		for scanner.Scan() {
			if out := processLine(scanner.Text()); out != "" {
				*versions = append(*versions, out)
			}
		}
	}(&versions)
	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		os.Exit(1)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		os.Exit(1)
	}
	return versions
}

// RepoUpdate update a helm repo
func RepoUpdate(helmInfo model.HelmInfo) {
	var (
		cmdOut  []byte
		err     error
		cmdName = "helm"
		cmdArgs = []string{"repo", "update", helmInfo.Repo}
	)
	if helmInfo.ValuesFile != "" {
		cmdArgs = append(cmdArgs, "-f")
		cmdArgs = append(cmdArgs, helmInfo.ValuesFile)
	}
	if cmdOut, err = runCmd(cmdName, cmdArgs); err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(cmdOut))
}

// HelmUpgrade do a hlem upgrade
func HelmUpgrade(helmInfo model.HelmInfo, version string) bool {
	var (
		cmdOut  []byte
		err     error
		cmdName = "helm"
		cmdArgs = []string{"--kube-context", helmInfo.KubeContext, "--namespace", helmInfo.Namespace, "upgrade", helmInfo.ReleaseName, "--version", version, helmInfo.Chart}
	)
	if helmInfo.ValuesFile != "" {
		cmdArgs = append(cmdArgs, "-f")
		cmdArgs = append(cmdArgs, helmInfo.ValuesFile)
	}
	if cmdOut, err = runCmd(cmdName, cmdArgs); err != nil {
		return false
	}
	fmt.Println(string(cmdOut))
	return true
}

// HelmInstall install something with helm
func HelmInstall(helmInfo model.HelmInfo, version string) bool {
	var (
		cmdOut  []byte
		err     error
		cmdName = "helm"
		cmdArgs = []string{"--kube-context", helmInfo.KubeContext, "--namespace", helmInfo.Namespace, "install", "--version", version, helmInfo.Chart}
	)
	if helmInfo.ValuesFile != "" {
		cmdArgs = append(cmdArgs, "-f")
		cmdArgs = append(cmdArgs, helmInfo.ValuesFile)
	}
	if cmdOut, err = runCmd(cmdName, cmdArgs); err != nil {
		return false
	}
	fmt.Println(string(cmdOut))
	return true
}

// CheckHelmSetup check that helm is setup and configured
func CheckHelmSetup(helmRepo string, helmURL string) error {
	_, err := exec.LookPath("helm")
	if err != nil {
		fmt.Printf("Helm not found: %s\n", err)
		return err
	}
	if !checkHelmRepo(helmRepo) {
		addHelmRepo(helmRepo, helmURL)
	}
	return nil
}
