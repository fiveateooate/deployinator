package clusterconfig

import (
	"fmt"
	"io/ioutil"

	"github.com/fiveateooate/deployinator/common"
	yaml "gopkg.in/yaml.v2"
)

// Service type holds info on a single service
type Service struct {
	Chart        string `yaml:"chart"`
	Version      string `yaml:"version"`
	DeployerType string `yaml:"deployertype"`
}

// Namespace type info on service in a single namespace
type Namespace struct {
	Name     string    `yaml:"name"`
	Services []Service `yaml:"services"`
}

// Deployment - list of namespaces and services to deplopy to them
type Deployment struct {
	Namespaces []Namespace `yaml:"namespaces"`
}

//ClusterConfig info on services that should be deployed in a cluster
type ClusterConfig struct {
	Cenv                string     `yaml:"cenv"`
	Cid                 string     `yaml:"cid"`
	ClusterDomain       string     `yaml:"clusterDomain"`
	CurrentVersionsFile string     `yaml:"current_versions_file"`
	Deployment          Deployment `yaml:"deployments"`
}

// ParseClusterConfig parse a file and return a map
func (cc *ClusterConfig) ParseClusterConfig(path string) {
	if !sharedfuncs.FileExists(path) {
		fmt.Printf("cc file %s does not exist\n", path)
		return
	}
	data, _ := ioutil.ReadFile(path)
	if err := yaml.Unmarshal([]byte(data), cc); err != nil {
		fmt.Println(err)
	}
}
