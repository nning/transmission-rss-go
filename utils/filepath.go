package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// IsExistsPath check path exist
func IsExistsPath(p string) bool {
	_, err := os.Stat(p)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// MkdirP like mkdir -p
func MkdirP(dirPath string) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.MkdirAll(dirPath, 0777)
	}
}

func Touch(filename string, defaultValue string) (err error) {
	MkdirP(filepath.Dir(filename))

	err = ioutil.WriteFile(filename, []byte(defaultValue), 0666)
	if err != nil {
		return
	}

	return nil
}

func TouchIfNotExist(filename string, defaultValue string) (err error) {
	if IsExistsPath(filename) {
		return nil
	}

	return Touch(filename, defaultValue)
}
