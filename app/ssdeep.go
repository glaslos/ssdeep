package main

import (
	"flag"
	"fmt"

	"github.com/glaslos/ssdeep"
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("Please provide a file path: ./ssdeep /tmp/file")
		return
	}
	sdeep := ssdeep.NewSSDEEP()
	sdeep.Fuzzy(flag.Args()[0])
}
