package sharedfuncs

import "os"

// FileExists - checks if a file exists and returns bool
func FileExists(path string) bool {
	exists := false
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		exists = true
	}
	return exists
}
