package clusterconfig

import (
	"fmt"
	"io/ioutil"

	"github.com/fiveateooate/deployinator/common"
	yaml "gopkg.in/yaml.v2"
)

// Service type holds info on a single service
type Service struct {
	Chart   string `yaml:"chart"`
	Version string `yaml:"version"`
}

// Namespace type info on service in a single namespace
type Namespace struct {
	Services []Service
}

//ClusterConfig info on services that should be deployed in a cluster
type ClusterConfig struct {
	clusterName struct {
		ClusterDomain       string `yaml:"clusterDomain"`
		CurrentVersionsFile string `yaml:"current_versions_file"`
		Services            struct {
			Namespaces []Namespace
		}
	}
}

// ParseClusterConfig parse a file and return a map
func (cc *ClusterConfig) ParseClusterConfig(path string) {
	if !sharedfuncs.FileExists(path) {
		fmt.Printf("cc file %s does not exist\n", path)
		return
	}
	data, _ := ioutil.ReadFile(path)
	fmt.Println(string(data))
	yaml.Unmarshal([]byte(data), cc)
	fmt.Println(cc)
}
