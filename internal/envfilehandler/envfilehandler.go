package envfilehandler

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

//Envfile - struct for yaml unmarshaling of envfile
type Envfile struct {
	Schema       string `yaml:"schema"`
	Slug         string `yaml:"slug"`
	FriendlyName string `yaml:"name"`
	Owner        string `yaml:"owner"`
	Repo         string `yaml:"repo"`
	Domain       string `yaml:"domain"`
	Slack        string `yaml:"slack"`
	HelmRepo     string `yaml:"helmrepo"`
}

//LoadEnvfile load up the struct from file
func (ef *Envfile) LoadEnvfile(path string) {
	data, _ := ioutil.ReadFile(path)
	if err := yaml.Unmarshal([]byte(data), ef); err != nil {
		fmt.Println(err)
	}
}

//LoadFromFlags - load up struct from flags
func (ef *Envfile) LoadFromFlags(name string, namespaceName string, helmrepo string) {
	ef.Slug = name
	ef.Domain = namespaceName
	ef.HelmRepo = helmrepo
}
