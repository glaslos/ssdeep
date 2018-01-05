package main

import (
	"flag"
	"fmt"
	"github.com/glaslos/ssdeep"
	"os"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Please provide a file path: ./ssdeep /tmp/file")
		os.Exit(1)
	}

	h1, err := ssdeep.FuzzyFilename(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(args) == 2 {
		h2, err := ssdeep.FuzzyFilename(args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		score, _ := ssdeep.Distance(h1, h2)
		if score != 0 {
			fmt.Printf("%s matches %s (%d)\n", args[0], args[1], score)
		} else {
			fmt.Println("The files doesn't match")
		}
	} else {
		fmt.Printf("%s,\"%s\"\n", h1, args[0])
	}
}
