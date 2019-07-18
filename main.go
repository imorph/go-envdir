package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

// exitOnErr111 exits on err
func exitOnErr111(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(111)
	}
}

func main() {
	// custom help message
	flag.Usage = func() {
		printHelp()
	}

	// checking parameters and files
	flag.Parse()
	inParams := flag.Args()
	exitOnErr111(checkInParams(inParams))
	envDir := inParams[0]
	envFiles, err := ioutil.ReadDir(envDir)
	exitOnErr111(err)
	exitOnErr111(checkEnvFiles(envFiles))

	// populating env
	currEnv := os.Environ()
	for _, envFile := range envFiles {
		envVarName := envFile.Name()
		if envFile.Size() == 0 {
			if varInEnv(currEnv, envVarName) {
				os.Unsetenv(envVarName)
			}
		} else {
			f, err := os.Open(envDir + "/" + envFile.Name())
			if err != nil {
				f.Close()
				exitOnErr111(err)
			}
			r := bufio.NewReader(f)
			envVarValue, err := readSingleLine(r)
			if err != nil {
				f.Close()
				exitOnErr111(err)
			}
			os.Setenv(envVarName, cleanValue(envVarValue))
			f.Close()
		}
	}

	// running child command
	childCommand := inParams[1]
	childCommandParams := inParams[2:]
	cmd := exec.Command(childCommand, childCommandParams...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Println(err)
	}
	// exit with childs command exit code
	os.Exit(cmd.ProcessState.ExitCode())
}
