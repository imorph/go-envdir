package main

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// checkInParams checks input parameters validity
func checkInParams(inParams []string) error {

	if len(inParams) < 2 {
		errString := "Not enough arguments: " + strconv.Itoa(len(inParams)) + " need (at least): 2"
		return errors.New(errString)
	}

	envDir := inParams[0]
	dirInfo, err := os.Stat(envDir)
	if os.IsNotExist(err) {
		errString := "Input directory: " + envDir + " does not exist, please enter valid directory path."
		return errors.New(errString)
	}
	if err != nil {
		return err
	}
	if !dirInfo.IsDir() {
		errString := "Input directory: " + envDir + " is not directory, please enter valid directory path."
		return errors.New(errString)
	}

	childCommand := inParams[1]
	_, err = exec.LookPath(childCommand)
	if err != nil {
		return err
	}
	return nil
}

// checkEnvFiles checks env files validity
func checkEnvFiles(envFiles []os.FileInfo) error {
	for _, file := range envFiles {
		if strings.Contains(file.Name(), "=") {
			errString := "File name: " + file.Name() + ` contains symbol: "=", this is not valid name for Env file, please rename it and try again`
			return errors.New(errString)
		}
	}
	return nil
}

// varInEnv checks that string with varName prefix within env slice
func varInEnv(env []string, varName string) bool {
	for _, envVar := range env {
		if strings.HasPrefix(envVar, varName) {
			return true
		}
	}
	return false
}
