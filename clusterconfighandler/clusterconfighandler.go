package clusterconfighandler

import (
	"fmt"

	"github.com/fiveateooate/deployinator/common"
	"github.com/fiveateooate/deployinator/model"
)

// ParseClusterConfig parse a file and return a map
func ParseClusterConfig(path string) model.ClusterConfig {
	var (
		clusterConfig model.ClusterConfig
	)
	if sharedfuncs.FileExists(path) {
		fmt.Println(path)
	}
	return clusterConfig
}
