package ssdeep

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	rollingWindow uint32 = 7
	blockMin             = 3
	spamSumLength        = 64
	minFileSize          = 4096
	hashPrime     uint32 = 0x01000193
	hashIinit     uint32 = 0x28021967
	b64String            = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

var b64 = []byte(b64String)

type rollingState struct {
	window []byte
	h1     uint32
	h2     uint32
	h3     uint32
	n      uint32
}

// FuzzyHash struct for comparison
type FuzzyHash struct {
	blockSize    int
	hashString1  string
	hashString2  string
	fileLocation string
}

func (h FuzzyHash) String() string {
	if h.fileLocation == "" {
		return fmt.Sprintf("%d:%s:%s", h.blockSize, h.hashString1, h.hashString2)
	}
	return fmt.Sprintf("%d:%s:%s,\"%s\"", h.blockSize, h.hashString1, h.hashString2, h.fileLocation)
}

// SSDEEP state struct
type SSDEEP struct {
	rollingState rollingState
	blockSize    int
	hashString1  string
	hashString2  string
	blockHash1   uint32
	blockHash2   uint32
}

// NewSSDEEP creates a new SSDEEP hash
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

// sumHash based on FNV hash
func sumHash(c byte, h uint32) uint32 {
	return (h * hashPrime) ^ uint32(c)
}

// rollHash based on Adler checksum
func (sdeep *SSDEEP) rollHash(c byte) uint32 {
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

func (sdeep *SSDEEP) getBlockSize(n int) {
	blockSize := blockMin
	for blockSize*spamSumLength < n {
		blockSize = blockSize * 2
	}
	sdeep.blockSize = blockSize
}

func getFileSize(f *os.File) (int, error) {
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return int(fi.Size()), nil
}

func (sdeep *SSDEEP) processByte(b byte) {
	sdeep.blockHash1 = sumHash(b, sdeep.blockHash1)
	sdeep.blockHash2 = sumHash(b, sdeep.blockHash2)
	rh := int(sdeep.rollHash(b))
	if rh%sdeep.blockSize == (sdeep.blockSize - 1) {
		if len(sdeep.hashString1) < spamSumLength-1 {
			sdeep.hashString1 += string(b64[sdeep.blockHash1%64])
			sdeep.blockHash1 = hashIinit
			sdeep.newRollingState()
		}
		if rh%(sdeep.blockSize*2) == ((sdeep.blockSize * 2) - 1) {
			if len(sdeep.hashString2) < spamSumLength/2-1 {
				sdeep.hashString2 += string(b64[sdeep.blockHash2%64])
				sdeep.blockHash2 = hashIinit
				sdeep.newRollingState()
			}
		}
	}
}

type fuzzyReader interface {
	ReadByte() (byte, error)
	Read([]byte) (int, error)
	Seek(offset int64, whence int) (int64, error)
}

func (sdeep *SSDEEP) process(r fuzzyReader) {
	sdeep.newRollingState()
	b, err := r.ReadByte()
	for err == nil {
		sdeep.processByte(b)
		b, err = r.ReadByte()
	}
	// Finalize the hash string with the remaining data
	sdeep.hashString1 += string(b64[sdeep.blockHash1%64])
	sdeep.hashString2 += string(b64[sdeep.blockHash2%64])
}

// FuzzyReader hash of a provided reader
func (sdeep *SSDEEP) FuzzyReader(r fuzzyReader, fileLocation string) (*FuzzyHash, error) {
	// This is not optimal as you have to read the whole file into memory
	if sdeep.blockSize == 0 {
		buf := &bytes.Buffer{}
		n, err := io.Copy(buf, r)
		if err != nil {
			return nil, err
		}
		if n < minFileSize {
			return nil, errors.New("Did not process files large enough to produce meaningful results")
		}
		sdeep.getBlockSize(int(n))
	}
	sdeep.process(r)
	for len(sdeep.hashString1) < spamSumLength/2 {
		if sdeep.blockSize < blockMin {
			return nil, errors.New("Unable to establish a sufficient block size")
		}
		sdeep.blockSize = sdeep.blockSize / 2
		r.Seek(0, 0)
		sdeep.process(r)
	}
	return &FuzzyHash{
		blockSize:    sdeep.blockSize,
		hashString1:  sdeep.hashString1,
		hashString2:  sdeep.hashString2,
		fileLocation: fileLocation,
	}, nil
}

// FuzzyByte hash of a provided byte array
func (sdeep *SSDEEP) FuzzyByte(blob []byte) (*FuzzyHash, error) {
	n := len(blob)
	if n < minFileSize {
		return nil, errors.New("Did not process files large enough to produce meaningful results")
	}
	sdeep.getBlockSize(n)
	r := bytes.NewReader(blob)
	return sdeep.FuzzyReader(r, "")
}
