package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func getConfig() map[string]string {
	conffilepath := fmt.Sprintf("%s/conf.json", getExecFilepath())
	confJSON, err := ioutil.ReadFile(conffilepath)

	if err != nil {
		panic(err)
	}

	var conf map[string]string
	json.Unmarshal(confJSON, &conf)
	return conf
}

func getExecFilepath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func sliceContainsUint64(slice []uint64, num uint64) bool {
	for _, sliceNum := range slice {
		if sliceNum == num {
			return true
		}
	}
	return false
}
