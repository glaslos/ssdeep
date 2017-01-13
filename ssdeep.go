package ssdeep

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

var rollingWindow uint32 = 7
var blockMin uint32 = 3
var spamSumLength uint32 = 64
var b64String = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
var b64 = []byte(b64String)
var hashPrime uint32 = 0x01000193
var hashIinit uint32 = 0x28021967

type rollingState struct {
	window []byte
	h1     uint32
	h2     uint32
	h3     uint32
	n      uint32
}

type SSDEEP struct {
	rollingState rollingState
	hashString1  string
	hashString2  string
	blockHash1   uint32
	blockHash2   uint32
}

func NewSSDEEP() SSDEEP {
	return SSDEEP{
		blockHash1: hashIinit,
		blockHash2: hashIinit,
		rollingState: rollingState{
			window: make([]byte, rollingWindow),
		},
	}
}

func (sdeep *SSDEEP) newRollingState() {
	sdeep.rollingState = rollingState{}
	sdeep.rollingState.window = make([]byte, rollingWindow)
}

func sumHash(c byte, h uint32) uint32 {
	return (h * hashPrime) ^ uint32(c)
}

func rollHash(sdeep *SSDEEP, c byte) uint32 {
	rs := &sdeep.rollingState
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
	return int(blockInt)
}

func getFileSize(f *os.File) uint32 {
	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}
	return uint32(fi.Size())
}

func getBlock(f *os.File, blockInt int, sdeep *SSDEEP) {
	r := bufio.NewReader(f)
	sdeep.newRollingState()
	b, err := r.ReadByte()
	for err == nil {
		sdeep.blockHash1 = sumHash(b, sdeep.blockHash1)
		sdeep.blockHash2 = sumHash(b, sdeep.blockHash2)
		rh := int(rollHash(sdeep, b))
		if rh%blockInt == (blockInt - 1) {
			sdeep.hashString1 += string(b64[sdeep.blockHash1%64])
			sdeep.blockHash1 = hashIinit
			sdeep.newRollingState()
		}
		if rh%(blockInt*2) == ((blockInt * 2) - 1) {
			sdeep.hashString2 += string(b64[sdeep.blockHash2%64])
			sdeep.blockHash2 = hashIinit
		}
		b, err = r.ReadByte()
	}
	sdeep.hashString1 += string(b64[sdeep.blockHash1%64])
	sdeep.hashString2 += string(b64[sdeep.blockHash2%64])
}

func Fuzzy(fileLocation string) {
	sdeep := NewSSDEEP()
	f, err := os.Open(fileLocation)
	if err != nil {
		panic(err)
	}
	n := getFileSize(f)
	blockInt := getBlockSize(n)
	fmt.Printf("block size: %d, file size: %d, bs/n: %d\n", blockInt, n, int(n)/blockInt)
	getBlock(f, blockInt, &sdeep)
	fmt.Printf("%d:%s:%s,\"%s\"\n", blockInt, sdeep.hashString1, sdeep.hashString2, fileLocation)
}
