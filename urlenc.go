package main

import (
	"bufio"
	"os"
)

func main() {
	var err error
	outputFile := os.Stdout
	inputFile := os.Stdin
	if fOutputPath != "" {
		outputFile, err = os.Create(fOutputPath)
		check(err, nil)
		defer outputFile.Close()
	}
	if fInputPath != "" {
		inputFile, err = os.Open(fInputPath)
		check(err, nil)
		defer inputFile.Close()
	}
	r := bufio.NewReaderSize(inputFile, fBufSize)
	w := bufio.NewWriterSize(outputFile, fBufSize)
	run(r, w)
}
