package tfvars

import (
	"fmt"
	"os"
)

func makeDirIfNotExists(path string) error {
	fmt.Println("creating dir", path)
	var err error
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Other than os.IsNotExist error, return
		return os.MkdirAll(path, os.ModePerm)
	}
	return err
}
