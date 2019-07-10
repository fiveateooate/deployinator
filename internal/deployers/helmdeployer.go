package deployers

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	pb "github.com/fiveateooate/deployinator/deployproto"
	sharedfuncs "github.com/fiveateooate/deployinator/internal/common"
	"github.com/wsxiaoys/terminal/color"
)

// HelmDeployer stuff need to do helm things
type HelmDeployer struct {
	Version        string
	Chart          string
	ChartPath      string
	TmpDir         string
	Repo           string
	AppName        string
	Namespace      string
	HelmHost       string
	ReleaseName    string
	ReleaseVersion string
	ValuesFile     string
	DeployResponse string
	ReleaseExists  bool
}

// HelmRelease hold info about a helm release
type HelmRelease struct {
	Name          string `json:"Name"`
	Revision      int    `json:"Revision"`
	Updated       string `json:"Updated"`
	Status        string `json:"Status"`
	Chart         string `json:"Chart"`
	AppVersion    string `json:"AppVersion"`
	Namespace     string `json:"Namespace"`
	ReleaseExists bool   `json:"ReleaseExists"`
}

// HelmListOutput struct to hold unmarshaled json
type HelmListOutput struct {
	Next     string        `json:"Next"`
	Releases []HelmRelease `json:"Releases"`
}

// process output of helm list
func processLine(text string) string {
	headerRE := regexp.MustCompile("^NAME.*")
	if headerRE.MatchString(text) {
		return ""
	}
	tmp := strings.Fields(text)
	return tmp[1]
}

// helmUpgrade do a helm upgrade
func (hi *HelmDeployer) helmUpgrade() bool {
	var (
		cmdOut  []byte
		err     error
		cmdName = "helm"
		cmdArgs = []string{"--host", hi.HelmHost, "--namespace", hi.Namespace, "upgrade", hi.ReleaseName, "--version", hi.Version, hi.ChartPath}
	)
	if hi.ValuesFile != "" {
		cmdArgs = append(cmdArgs, "-f")
		cmdArgs = append(cmdArgs, hi.ValuesFile)
	}
	if cmdOut, err = sharedfuncs.RunCmd(cmdName, cmdArgs); err != nil {
		hi.DeployResponse += color.Sprintf("@rRelease upgrade failed: %s\n%s\n", err, cmdArgs)
		return false
	}
	hi.DeployResponse += color.Sprintf(string(cmdOut))
	return true
}

// helmInstall install something with helm
func (hi *HelmDeployer) helmInstall() bool {
	var (
		cmdOut  []byte
		err     error
		cmdName = "helm"
		cmdArgs = []string{"--host", hi.HelmHost, "--namespace", hi.Namespace, "install", "--version", hi.Version, "--name", hi.ReleaseName, hi.ChartPath}
	)
	if hi.ValuesFile != "" {
		cmdArgs = append(cmdArgs, "-f")
		cmdArgs = append(cmdArgs, hi.ValuesFile)
	}
	if cmdOut, err = sharedfuncs.RunCmd(cmdName, cmdArgs); err != nil {
		hi.DeployResponse += color.Sprintf("@rInstall failed: %s\n%s\n", err, cmdArgs)
		return false
	}
	hi.DeployResponse += color.Sprintf(string(cmdOut))
	return true
}

// findValuesFile - Use CENV and CID env vars to pick a values-xxx.yaml file
//              defaults to just plain old values
//              order of preference -> values-CENV:CID.yaml values-CENV.yaml values.yaml
func (hi *HelmDeployer) findValuesFile(cenv string, cid string) {
	tmpdir := fmt.Sprintf("%s/%s/", hi.TmpDir, hi.AppName)
	hi.DeployResponse += color.Sprintf("@cLooking for values file ...\n")
	tmp := color.Sprintf("@yno values file found\n")
	possibleValuesFiles := []string{tmpdir + "values-" + cenv + ":" + cid + ".yaml", tmpdir + "values-" + cenv + ".yaml", tmpdir + "values.yaml"}
	for _, filename := range possibleValuesFiles {
		if sharedfuncs.FileExists(filename) {
			hi.ValuesFile = filename
			tmp = color.Sprintf("using %s\n", strings.Replace(filename, tmpdir, "", -1))
			break
		}
	}
	hi.DeployResponse += tmp
}

// getRelease return something release
func (hi *HelmDeployer) getRelease() {
	regex := fmt.Sprintf("^%s(|-[a-z0-9]{5})$", hi.AppName)
	r := regexp.MustCompile(regex)
	// preset some values in case release is not found
	tmp := color.Sprintf("not found\n")
	hi.ReleaseName = fmt.Sprintf("%s-%s", hi.AppName, sharedfuncs.RandString(5))
	hi.ReleaseExists = false
	hi.DeployResponse += color.Sprintf("@cSearching for helm release ...\n")
	releases, _ := hi.listReleases()
	for _, release := range releases.Releases {
		match := r.FindString(release.Name)
		if match != "" {
			hi.ReleaseName = release.Name
			hi.ReleaseVersion = release.AppVersion
			hi.ReleaseExists = true
			tmp = color.Sprintf("found %s\n", hi.ReleaseName)
		}
	}
	hi.DeployResponse += tmp
}

//FetchChart - fetch a chart from helm repo and untar it
func (hi *HelmDeployer) fetchChart() error {
	var (
		cmdName = "helm"
		cmdArgs = []string{"fetch", "--untar", "--version", hi.Version, hi.Chart, "--untardir", hi.TmpDir}
		cmdOut  []byte
		err     error
	)
	hi.DeployResponse += color.Sprintf("@cFetching chart ...\n")
	if cmdOut, err = sharedfuncs.RunCmd(cmdName, cmdArgs); err != nil {
		hi.DeployResponse += color.Sprintf("@rfailed - %v", cmdOut)
		return err
	}
	hi.DeployResponse += color.Sprintf("%s\n", hi.Chart)
	return nil
}

// listReleases return struct of releases
func (hi *HelmDeployer) listReleases() (HelmListOutput, error) {
	var (
		cmdOut []byte
		err    error
		output HelmListOutput
	)
	cmd := "helm"
	args := []string{"--host", hi.HelmHost, "--namespace", hi.Namespace, "--output", "json", "list"}
	if cmdOut, err = sharedfuncs.RunCmd(cmd, args); err != nil {
		return output, err
	}
	json.Unmarshal(cmdOut, &output)
	return output, nil
}

// RepoUpdate update a helm repo
func (hi *HelmDeployer) repoUpdate() {
	var (
		cmdOut  []byte
		err     error
		cmdName = "helm"
		cmdArgs = []string{"repo", "update", hi.Repo}
	)
	hi.DeployResponse += color.Sprintf("@cUpdating helm repositories ...\n")
	if cmdOut, err = sharedfuncs.RunCmd(cmdName, cmdArgs); err != nil {
		fmt.Println(err)
	}
	hi.DeployResponse += color.Sprintf("%s\n", cmdOut)
}

// NewHelmDeployer a new helmdeployer with some defaults set
func NewHelmDeployer(appname string, namespace string, version string, repo string) *HelmDeployer {
	hi := new(HelmDeployer)
	hi.Repo = repo
	hi.AppName = appname
	hi.Namespace = namespace
	hi.HelmHost = "tiller-deploy:44134"
	hi.Version = version
	hi.Chart = fmt.Sprintf("%s/%s", hi.Repo, hi.AppName)
	hi.TmpDir = "/tmp"
	hi.ValuesFile = ""
	hi.DeployResponse = ""
	hi.ChartPath = fmt.Sprintf("%s/%s", hi.TmpDir, hi.AppName)
	return hi
}

// HelmDeploy - deploy a service or whatever
func (hi *HelmDeployer) HelmDeploy(msg *pb.DeployMessage) error {
	hi.getRelease()
	if hi.ReleaseExists && (msg.Version == hi.ReleaseVersion) {
		hi.DeployResponse += color.Sprintf("Version %s already deployed\n", msg.Version)
		return nil
	}
	hi.repoUpdate()
	if err := hi.fetchChart(); err != nil {
		return err
	}
	hi.findValuesFile(msg.Cenv, msg.Cid)
	hi.DeployResponse += color.Sprintf("@cDeploying %s ...\n", hi.AppName)
	if hi.ReleaseExists {
		hi.helmUpgrade()
	} else {
		hi.helmInstall()
	}
	hi.DeployResponse += color.Sprintf("@gSuccess\n")
	return nil
}
