package tfvars

import (
	"os"
)

func makeDirIfNotExists(path string) error {
	var err error
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Other than os.IsNotExist error, return
		return os.MkdirAll(path, os.ModePerm)
	}
	return err
}
