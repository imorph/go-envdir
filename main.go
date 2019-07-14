package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	au "github.com/logrusorgru/aurora"
)

func printHelp() {
	fmt.Println("")
	fmt.Println(au.BrightMagenta("go-envdir"), " - runs another program with environment modified according to files in a specified")
	fmt.Println("directory.")
	fmt.Println("")
	fmt.Println(au.BrightGreen(au.Bold("USAGE:")))
	fmt.Println("go-envdir /env/dir command")
	fmt.Println("")
	fmt.Println(au.BrightGreen(au.Bold("DESCRIPTION:")))
	fmt.Println(`"/env/dir" is a single argument, "command" consists of one or more arguments.

go-envdir sets various environment variables as specified by files in the directory  named 
"/env/dir". It then runs "command".

If "/env/dir" contains a file named s whose first line is t, envdir removes an environment 
variable named s if one exists, and then adds an environment variable named s with  
value t. The name  s  must  not  contain =. Spaces and tabs at the end of t are removed.
Nulls in t are changed to newlines in the environment variable.

If the file s is completely empty (0 bytes long), envdir removes an  environment  variable
named s if one exists, without adding a new variable.`)
	fmt.Println("")
}

func printUsage() {
	fmt.Println(au.BrightGreen(au.Bold("USAGE:")), "go-envdir /env/dir command")
}

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

func checkEnvFiles(envFiles []os.FileInfo) error {
	for _, file := range envFiles {
		if strings.Contains(file.Name(), "=") {
			errString := "File name: " + file.Name() + ` contains symbol: "=", this is not valid name for Env file, please rename it and try again`
			return errors.New(errString)
		}
	}
	return nil
}

func readSingleLine(reader *bufio.Reader) (string, error) {

	isPrefix := true
	var err error
	var rawLine, outLine []byte

	for isPrefix && err == nil {
		rawLine, isPrefix, err = reader.ReadLine()
		outLine = append(outLine, rawLine...)
	}
	return string(outLine), err
}

func varInEnv(env []string, varName string) bool {
	for _, envVar := range env {
		if strings.HasPrefix(envVar, varName) {
			return true
		}
	}
	return false
}

func main() {
	flag.Parse()
	inParams := flag.Args()
	err := checkInParams(inParams)
	if err != nil {
		log.Println(err)
		printUsage()
		os.Exit(111)
	}
	envDir := inParams[0]
	childCommand := inParams[1]
	childCommandParams := inParams[2:]
	// fmt.Println("will execute:", childCommand, "with params:", childCommandParams, "and environment from dir:", envDir)
	envFiles, err := ioutil.ReadDir(envDir)
	if err != nil {
		log.Println(err)
		os.Exit(111)
	}
	err = checkEnvFiles(envFiles)
	if err != nil {
		log.Println(err)
		os.Exit(111)
	}
	//var currEnv []string
	currEnv := os.Environ()
	for _, envFile := range envFiles {
		envVarName := envFile.Name()
		if envFile.Size() == 0 {
			// fmt.Println("here we delete envvar from env if it exist")
			if varInEnv(currEnv, envVarName) {
				os.Unsetenv(envVarName)
			}
		}
		if envFile.Size() != 0 {
			f, err := os.Open(envDir + "/" + envFile.Name())
			if err != nil {
				log.Println(err)
				f.Close()
				os.Exit(111)
			}
			r := bufio.NewReader(f)
			envVarValue, err := readSingleLine(r)
			if err != nil {
				log.Println(err)
				f.Close()
				os.Exit(111)
			}
			os.Setenv(envVarName, envVarValue)
			f.Close()
		}
	}
	cmd := exec.Command(childCommand, childCommandParams...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Println(err)
	}
	exitCode := cmd.ProcessState.ExitCode()
	os.Exit(exitCode)
}
