package model

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
