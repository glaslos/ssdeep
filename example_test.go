package ssdeep_test

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/glaslos/ssdeep"
)

func ExampleFuzzyFile() {
	f, err := os.Open("file.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h, err := ssdeep.FuzzyFile(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(h)
}

func ExampleFuzzyBytes() {
	buffer := make([]byte, 4097)
	rand.Read(buffer)
	h, err := ssdeep.FuzzyBytes(buffer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(h)
}

func ExampleFuzzyReader() {
	buffer := make([]byte, 4097)
	rand.Read(buffer)
	h, err := ssdeep.FuzzyReader(bytes.NewReader(buffer))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(h)
}

func ExampleFuzzyFilename() {
	h, err := ssdeep.FuzzyFilename("file.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(h)
}
