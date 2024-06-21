package vlog

import (
	"log"
	"os"
)

var verbose = false

func Printf(format string, v ...any) {
	if !verbose {
		return
	}

	log.Printf(format, v...)
}

func Println(v ...any) {
	if !verbose {
		return
	}
	log.Println(v...)
}

func init() {
	if len(os.Args) >= 4 {
		verbose = os.Args[3] == "-v"
	}
}
