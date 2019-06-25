package tfvarser

import "os"

func makeDirIfNotExists(path string) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModeDir)
	}
}
