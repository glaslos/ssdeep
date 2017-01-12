package ssdeep

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"math"
	"os"
)

var rollingWindow uint32 = 7
var blockMin uint32 = 3
var spamSumLength uint32 = 64

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

type rollingState struct {
	window []byte
	h1     uint32
	h2     uint32
	h3     uint32
	n      uint32
}

func newRollingState() rollingState {
	return rollingState{
		window: make([]byte, rollingWindow),
	}
}

func rollHash(rs *rollingState, c byte) uint32 {
	rs.h2 -= rs.h1
	rs.h2 += rollingWindow * uint32(c)
	rs.h1 += uint32(c)
	rs.h1 -= uint32(rs.window[rs.n])
	rs.window[rs.n] = c
	rs.n++
	if rs.n == rollingWindow {
		rs.n = 0
	}
	rs.h3 = rs.h3 << 5
	rs.h3 ^= uint32(c)
	return rs.h1 + rs.h2 + rs.h3
}

func getBlockSize(n uint32) int {
	blockInt := int(blockMin) * int(math.Exp2(math.Floor(math.Log2(float64(n/(spamSumLength*blockMin))))))
	/*blockInt := blockMin
	for blockInt*spamSumLength < n {
		blockInt = blockInt * 2
	}*/
	return int(blockInt)
}

func getFileSize(f *os.File) uint32 {
	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}
	return uint32(fi.Size())
}

func getBlocks() {
	f, err := os.Open("/tmp/dat")
	if err != nil {
		panic(err)
	}

	n := getFileSize(f)
	blockInt := getBlockSize(n)
	fmt.Printf("block size: %d, file size: %d, bs/n: %d\n", blockInt, n, int(n)/blockInt)

	rs := newRollingState()
	r := bufio.NewReader(f)
	b, err := r.ReadByte()
	nb := 0
	for err == nil {
		h := int(rollHash(&rs, b))
		//fmt.Printf("state: %d, trigger %d\n", h%blockInt, blockInt-1)
		if h%blockInt == (blockInt - 1) {
			nb++
		}
		b, err = r.ReadByte()
	}
	fmt.Println(nb)
}
