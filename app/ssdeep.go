package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/glaslos/ssdeep"
)

var (
	// VERSION is set by the makefile
	VERSION = "v0.0.0"
	// BUILDDATE is set by the makefile
	BUILDDATE = ""
)

var (
	forceHash   bool
	showVersion bool
)

func main() {
	fmt.Printf("ssdeep,%s--blocksize:hash:hash,filename\n", VERSION)

	flag.BoolVar(&forceHash, "force", false, "Force hash on error or invalid input length")
	flag.BoolVar(&showVersion, "version", false, "Print version")
	flag.Parse()

	if showVersion {
		fmt.Printf("%s %s\n", VERSION, BUILDDATE)
		return
	}
	ssdeep.Force = forceHash

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Please provide a file path: ./ssdeep /tmp/file")
		os.Exit(1)
	}

	h1, err := ssdeep.FuzzyFilename(args[0])
	if err != nil && !ssdeep.Force {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(args) == 2 {
		var h2 string
		h2, err = ssdeep.FuzzyFilename(args[1])
		if err != nil && !ssdeep.Force {
			fmt.Println(err)
			os.Exit(1)
		}

		var score int
		score, err = ssdeep.Distance(h1, h2)
		if score != 0 {
			fmt.Printf("%s matches %s (%d)\n", args[0], args[1], score)
		} else if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("The files doesn't match")
		}
	} else {
		if err != nil {
			fmt.Printf("%s,\"%s\"\n%s\n", h1, args[0], err)
		} else {
			fmt.Printf("%s,\"%s\"\n", h1, args[0])
		}
	}
}
