package helmbuddy

import (
	"fmt"
	"regexp"

	"github.com/wsxiaoys/terminal/color"
)

// HelmInfo stuff need to do helm things
type HelmInfo struct {
	ValuesFile     string
	Version        string
	Chart          string
	Repo           string
	AppName        string
	Namespace      string
	KubeContext    string
	ReleaseName    string
	ReleaseVersion string
	ReleaseExists  bool
}

// getRelease return something release
func (hi *HelmInfo) getRelease() {
	regex := fmt.Sprintf("^%s(|-[a-z0-9]{5})$", hi.AppName)
	r := regexp.MustCompile(regex)
	color.Printf("@cSearching for helm release ...")
	releases, _ := ListReleases(hi.Namespace, hi.KubeContext)
	for _, release := range releases.Releases {
		match := r.FindString(release.Name)
		if match != "" {
			hi.ReleaseName = release.Name
			hi.ReleaseVersion = release.AppVersion
			hi.ReleaseExists = true
			color.Printf("found %s\n", hi.ReleaseName)
			return
		}
	}
	hi.ReleaseName = fmt.Sprintf("%s-%s", hi.AppName, RandStringBytes(5))
	hi.ReleaseExists = false
	color.Printf("not found\n")
}

// GetHelmInfo load up a helminfo
func (hi *HelmInfo) GetHelmInfo(appname string, namespace string, repo string, values string, context string) {
	hi.ValuesFile = values
	hi.Repo = repo
	hi.AppName = appname
	hi.Namespace = namespace
	hi.KubeContext = context
	hi.Chart = fmt.Sprintf("%s/%s", repo, appname)
	hi.getRelease()
}
