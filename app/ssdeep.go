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
	if len(flag.Args()) == 1 {
		sdeep := ssdeep.NewSSDEEP()
		fmt.Println(sdeep.Fuzzy(flag.Args()[0]))
	}
	if len(flag.Args()) == 2 {
		sdeep := ssdeep.NewSSDEEP()
		h1 := sdeep.Fuzzy(flag.Args()[0])
		sdeep = ssdeep.NewSSDEEP()
		h2 := sdeep.Fuzzy(flag.Args()[1])
		fmt.Printf("%s\n%s\n", h1, h2)
		score := ssdeep.Distance(h1, h2)
		fmt.Println(score)
	}
}
