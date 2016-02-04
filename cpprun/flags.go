package main

import (
	"fmt"
	"os"
)

const usage = `Usage of cpprun:
    cpprun <source> [ -l libA libB ... ]
    
source
    The C++ file to compile, link and run.
    
-l
    The -l flag indicates a list of libraries to link with is following.
    Libraries are seperated by spaces. There can be a number of 1 to n libs.
`

type flags struct {
	source    string
	libraries []string
}

func printUsage(err string) {
	fmt.Println("[ERROR]", err)
	fmt.Print(usage)
	os.Exit(1)
}

func parseFlags(args []string) *flags {
	flags := new(flags)

	libFlagSet := false
	for _, arg := range args {
		if arg == "-l" {
			// set libflag
			if !libFlagSet {
				libFlagSet = true
				continue
			} else {
				printUsage("args: -l flag defined more than once")
			}
		}

		if !libFlagSet {
			// set source
			if len(flags.source) == 0 {
				flags.source = arg
			} else {
				printUsage("args: source defined more than once")
			}
		} else {
			// add to libs
			flags.libraries = append(flags.libraries, arg)
		}
	}

	if len(flags.source) == 0 {
		printUsage("args: source not defined")
	}
	if libFlagSet && len(flags.libraries) == 0 {
		printUsage("args: -l flag without libraries defined")
	}

	return flags
}
