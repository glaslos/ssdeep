package ssdeep

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
)

const (
	rollingWindow uint32 = 7
	blockMin             = 3
	spamSumLength        = 64
	minFileSize          = 4096
	hashPrime     uint32 = 0x01000193
	hashInit      uint32 = 0x28021967
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

func (rs rollingState) rollSum() uint32 {
	return rs.h1 + rs.h2 + rs.h3
}

type ssdeepState struct {
	rollingState rollingState
	blockSize    int
	hashString1  string
	hashString2  string
	blockHash1   uint32
	blockHash2   uint32
}

func newSsdeepState() ssdeepState {
	return ssdeepState{
		blockHash1: hashInit,
		blockHash2: hashInit,
		rollingState: rollingState{
			window: make([]byte, rollingWindow),
		},
	}
}

func (state *ssdeepState) newRollingState() {
	state.rollingState = rollingState{}
	state.rollingState.window = make([]byte, rollingWindow)
}

// sumHash based on FNV hash
func sumHash(c byte, h uint32) uint32 {
	return (h * hashPrime) ^ uint32(c)
}

// rollHash based on Adler checksum
func (state *ssdeepState) rollHash(c byte) {
	rs := &state.rollingState
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
}

// getBlockSize calculates the block size based on file size
func (state *ssdeepState) getBlockSize(n int) {
	blockSize := blockMin
	for blockSize*spamSumLength < n {
		blockSize = blockSize * 2
	}
	state.blockSize = blockSize
}

// getFileSize returns the files size
func getFileSize(f *os.File) (int, error) {
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return int(fi.Size()), nil
}

func (state *ssdeepState) processByte(b byte) {
	state.blockHash1 = sumHash(b, state.blockHash1)
	state.blockHash2 = sumHash(b, state.blockHash2)
	state.rollHash(b)
	rh := int(state.rollingState.rollSum())
	if rh%state.blockSize == (state.blockSize - 1) {
		if len(state.hashString1) < spamSumLength-1 {
			state.hashString1 += string(b64[state.blockHash1%64])
			state.blockHash1 = hashInit
		}
		if rh%(state.blockSize*2) == ((state.blockSize * 2) - 1) {
			if len(state.hashString2) < spamSumLength/2-1 {
				state.hashString2 += string(b64[state.blockHash2%64])
				state.blockHash2 = hashInit
			}
		}
	}
}

type fuzzyReader interface {
	io.Seeker
	io.Reader
}

func (state *ssdeepState) process(r *bufio.Reader) {
	state.newRollingState()
	b, err := r.ReadByte()
	for err == nil {
		state.processByte(b)
		b, err = r.ReadByte()
	}
}

func (state *ssdeepState) fuzzyReader(f fuzzyReader, n int) (string, error) {
	if n < minFileSize {
		return "", errors.New("Did not process files large enough to produce meaningful results")
	}

	state.getBlockSize(n)
	for {
		f.Seek(0, 0)
		r := bufio.NewReader(f)
		state.process(r)
		if state.blockSize < blockMin {
			return "", errors.New("Unable to establish a sufficient block size")
		}
		if len(state.hashString1) < spamSumLength/2 {
			state.blockSize = state.blockSize / 2
			state.blockHash1 = hashInit
			state.blockHash2 = hashInit
			state.hashString1 = ""
			state.hashString2 = ""
		} else {
			rh := state.rollingState.rollSum()
			if rh != 0 {
				// Finalize the hash string with the remaining data
				state.hashString1 += string(b64[state.blockHash1%64])
				state.hashString2 += string(b64[state.blockHash2%64])
			}
			break
		}
	}
	return fmt.Sprintf("%d:%s:%s", state.blockSize, state.hashString1, state.hashString2), nil
}

func FuzzyFilename(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return FuzzyFile(f)
}

func FuzzyFile(f *os.File) (string, error) {
	currentPosition, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return "", err
	}

	f.Seek(0, io.SeekStart)
	state := newSsdeepState()
	n, err := getFileSize(f)
	if err != nil {
		return "", err
	}

	result, err := state.fuzzyReader(f, n)
	if err != nil {
		return "", err
	}

	f.Seek(currentPosition, io.SeekStart)
	return result, nil
}

func FuzzyBytes(buffer []byte) (string, error) {
	state := newSsdeepState()
	n := len(buffer)
	br := bytes.NewReader(buffer)

	result, err := state.fuzzyReader(br, n)
	if err != nil {
		return "", err
	}

	return result, nil
}


// Hash interface related struct and functions
type digest struct {
	buffer []byte
}

func (d *digest) Write(p []byte) (n int, err error) {
	n = len(p)
	d.buffer = append(d.buffer, p...)
	return
}

func (d *digest) Sum(b []byte) []byte {
	result, err := FuzzyBytes(d.buffer)
	if err != nil {
		return nil
	}
	return append(b, []byte(result)...)
}

func (d *digest) Reset() {
	d.buffer = make([]byte, 0, minFileSize)
}

func (d *digest) Size() int {
	return len(d.buffer)
}

func (d *digest) BlockSize() int {
	return minFileSize
}

func New() hash.Hash {
	digest := new(digest)
	digest.Reset()
	return digest
}
