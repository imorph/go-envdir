package main

import (
	"fmt"

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
