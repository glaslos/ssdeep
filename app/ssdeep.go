package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/glaslos/ssdeep"
)

func readFile(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(f)
	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, r)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("Please provide a file path: ./ssdeep /tmp/file")
		return
	}
	sdeep := ssdeep.NewSSDEEP()
	b, err := readFile(flag.Args()[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	h1, err := sdeep.FuzzyByte(b)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(flag.Args()) == 2 {
		sdeep := ssdeep.NewSSDEEP()
		b, err := readFile(flag.Args()[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		h2, err := sdeep.FuzzyByte(b)
		if err != nil {
			fmt.Println(err)
			return
		}
		score := ssdeep.Distance(h1, h2)
		fmt.Println(score)
	}
}
