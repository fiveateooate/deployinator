package artifactsbuddy

import "fmt"

// GetPkgs return list of versions for specified app
func GetPkgs(appName string) {
	fmt.Printf("Fetching latest packages for: %s\n", appName)
}
