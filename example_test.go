package ssdeep_test

import (
	"github.com/glaslos/ssdeep"
	"math/rand"
	"fmt"
	"os"
	"log"
)

func ExampleNew() {
	h := ssdeep.New()

	// Create data larger than 4096 bytes
	data := make([]byte, 4097)
	rand.Read(data)

	h.Write(data)
	fmt.Printf("%s", h.Sum(nil))
}

func ExampleFuzzyFilename() {
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
