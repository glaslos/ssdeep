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
		h1, err := sdeep.Fuzzy(flag.Args()[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		sdeep = ssdeep.NewSSDEEP()
		h2, err := sdeep.Fuzzy(flag.Args()[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%s\n%s\n", h1, h2)
		score := ssdeep.Distance(h1, h2)
		fmt.Println(score)
	}
}
