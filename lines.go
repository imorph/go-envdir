package main

import (
	"bufio"
	"strings"
)

// readSingleLine tries read one line from provided reader
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

func cleanValue(inValue string) string {
	if inValue == "" {
		return ""
	}
	v := strings.TrimRight(inValue, ` 	`)
	return strings.ReplaceAll(v, `\0`, "\n")
}
