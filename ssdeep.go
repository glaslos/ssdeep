package ssdeep

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
)

var (
	ErrFileTooSmall      = errors.New("did not process files large enough to produce meaningful results")
	ErrBlockSizeTooSmall = errors.New("unable to establish a sufficient block size")
	ErrZeroBlockSize     = errors.New("reached zero block size, unable to compute hash")
	ErrFileTooBig        = errors.New("input file length exceeds max processable length")
)

type Hash interface {
	io.Writer
	Sum(b []byte) []byte
}

const (
	rollingWindow     = 7
	blockMin          = 3
	spamSumLength     = 64
	minFileSize       = 4096
	hashPrime         = 0x93
	hashInit          = 0x27
	b64String         = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	halfspamSumLength = spamSumLength / 2
	numBlockhashes    = 31
	maxTotalSize      = blockMin << (numBlockhashes - 1) * spamSumLength
)

var (
	b64 = []byte(b64String)
	// Force calculates the hash on invalid input
	Force = false
)

type rollingState struct {
	window [rollingWindow + 1]byte
	h1     uint32
	h2     uint32
	h3     uint32
	n      uint32
}

func (rs *rollingState) rollSum() uint32 {
	return rs.h1 + rs.h2 + rs.h3
}
func (rs *rollingState) rollReset() {
	rs.h1 = 0
	rs.h2 = 0
	rs.h3 = 0
	rs.n = 0
	for i := 0; i < len(rs.window); i++ {
		rs.window[i] = 0
	}
}

type ssdeepState struct {
	rollingState rollingState
	iStart, iEnd int
	totalSize    uint64
	bsizeMask    uint32
	blocks       [numBlockhashes]blockHashState
}

type blockHashState struct {
	hashString             []byte
	blockSize              uint32
	blockHash1, blockHash2 byte
	tail1, tail2           byte
}

func newSSDEEPState() ssdeepState {
	s := ssdeepState{
		iEnd: 1,
	}
	for i := range s.blocks {
		s.blocks[i].blockSize = blockMin << i
		s.blocks[i].blockHash1 = hashInit
		s.blocks[i].blockHash2 = hashInit
	}
	return s
}

// sumHash based on FNV hash
func sumHash(c byte, h byte) byte {
	return ((h * hashPrime) ^ c) % 64
}

// rollHash based on Adler checksum
func (rs *rollingState) rollHash(c byte) {
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

func (state *ssdeepState) processByte(b byte) {
	for i := state.iStart; i < state.iEnd; i++ {
		state.blocks[i].blockHash1 = sumHash(b, state.blocks[i].blockHash1)
		state.blocks[i].blockHash2 = sumHash(b, state.blocks[i].blockHash2)
	}

	state.rollingState.rollHash(b)
	rh := state.rollingState.rollSum()
	if rh == math.MaxUint32 {
		return
	}
	// rh % 2**N > 0 will match all rh % (3 * 2**N) > 0 but can use fast bitmask
	if ((rh+1)/blockMin)&state.bsizeMask > 0 {
		return
	}
	if (rh+1)%blockMin > 0 {
		return
	}

	for i := state.iStart; i < state.iEnd; i++ {
		block := &state.blocks[i]
		if rh%block.blockSize == (block.blockSize - 1) {
			if len(block.hashString) == 0 {
				old := &state.blocks[state.iEnd-1]
				if state.iEnd <= numBlockhashes-1 {
					newb := &state.blocks[state.iEnd]
					newb.blockHash1 = old.blockHash1
					newb.blockHash2 = old.blockHash2
					state.iEnd++
				}
			}
			block.tail1 = block.blockHash1
			block.tail2 = block.blockHash2
			if len(block.hashString) < spamSumLength-1 {
				block.hashString = append(block.hashString, block.tail1)
				block.tail1 = 0
				block.blockHash1 = hashInit
				if len(block.hashString) < halfspamSumLength {
					block.blockHash2 = hashInit
					block.tail2 = 0
				}
			} else if state.isStartBlockFull() {
				state.iStart++
				state.bsizeMask = (state.bsizeMask << 1) + 1
			}
		}
	}
}

func (state *ssdeepState) isStartBlockFull() bool {
	return state.totalSize > uint64(state.blocks[state.iStart].blockSize*spamSumLength) &&
		len(state.blocks[state.iStart+1].hashString) >= halfspamSumLength
}

// Reader is the minimum interface that ssdeep needs in order to calculate the fuzzy hash.
// Reader groups io.Seeker and io.Reader.
type Reader interface {
	io.Seeker
	io.Reader
}

func (state *ssdeepState) Write(r []byte) (n int, err error) {
	state.totalSize += uint64(len(r))
	for _, b := range r {
		state.processByte(b)
	}
	return len(r), nil
}

// FuzzyReader computes the fuzzy hash of a Reader interface with a given input size.
// It is the caller's responsibility to append the filename, if any, to result after computation.
// Returns an error when ssdeep could not be computed on the Reader.
func FuzzyReader(f io.Reader) (string, error) {
	state := newSSDEEPState()
	if _, err := io.Copy(&state, f); err != nil {
		return "", err
	}
	digest, err := state.digest()
	return digest, err
}

func (state *ssdeepState) digest() (string, error) {
	if !Force && state.totalSize <= minFileSize {
		return "", ErrFileTooSmall
	}
	if state.totalSize > maxTotalSize {
		return "", ErrFileTooBig
	}

	var i = state.iStart
	for ; uint64(uint32(blockMin)<<i*spamSumLength) < state.totalSize; i++ {
	}

	if i >= state.iEnd {
		i = state.iEnd - 1
	}
	for i > state.iStart && len(state.blocks[i].hashString) < halfspamSumLength {
		i--
	}
	var bl1 = state.blocks[i]
	var bl2 = state.blocks[i+1]
	if i >= state.iEnd-1 {
		bl2 = state.blocks[i]
		bl2.hashString = append([]byte{}, bl1.hashString...)
	}
	var rh = state.rollingState.rollSum()

	if len(bl2.hashString) > halfspamSumLength-1 {
		bl2.hashString = bl2.hashString[:halfspamSumLength-1]
	}

	if rh != 0 {
		bl1.hashString = append(bl1.hashString, bl1.blockHash1)
		bl2.hashString = append(bl2.hashString, bl2.blockHash2)
	} else {
		if len(bl1.hashString) == spamSumLength-1 && bl1.tail1 != 0 {
			bl1.hashString = append(bl1.hashString, bl1.tail1)
		}
		if bl2.tail2 != 0 {
			bl2.hashString = append(bl2.hashString, bl2.tail2)
		}
	}
	for i := range bl1.hashString {
		bl1.hashString[i] = b64[bl1.hashString[i]]
	}
	for i := range bl2.hashString {
		bl2.hashString[i] = b64[bl2.hashString[i]]
	}

	return fmt.Sprintf("%d:%s:%s", bl1.blockSize, bl1.hashString, bl2.hashString), nil
}

// FuzzyFilename computes the fuzzy hash of a file.
// FuzzyFilename will opens, reads, and hashes the contents of the file 'filename'.
// It is the caller's responsibility to append the filename to the result after computation.
// Returns an error when the file doesn't exist or ssdeep could not be computed on the file.
func FuzzyFilename(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return FuzzyFile(f)
}

// FuzzyFile computes the fuzzy hash of a file using os.File pointer.
// FuzzyFile will computes the fuzzy hash of the contents of the open file, starting at the beginning of the file.
// When finished, the file pointer is returned to its original position.
// If an error occurs, the file pointer's value is undefined.
// It is the callers's responsibility to append the filename to the result after computation.
// Returns an error when ssdeep could not be computed on the file.
func FuzzyFile(f *os.File) (string, error) {
	currentPosition, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return "", err
	}
	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	out, err := FuzzyReader(f)
	if err != nil {
		return out, err
	}
	_, err = f.Seek(currentPosition, io.SeekStart)
	return out, err
}

// FuzzyBytes computes the fuzzy hash of a slice of byte.
// It is the caller's responsibility to append the filename, if any, to result after computation.
// Returns an error when ssdeep could not be computed on the buffer.
func FuzzyBytes(buffer []byte) (string, error) {
	br := bytes.NewReader(buffer)

	result, err := FuzzyReader(br)
	if err != nil {
		return "", err
	}

	return result, nil
}
