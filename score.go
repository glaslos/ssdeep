package ssdeep

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

var (
	// ErrEmptyHash is returned when no hash string is provided for scoring.
	ErrEmptyHash = errors.New("empty string")

	// ErrInvalidFormat is returned when a hash string is malformed.
	ErrInvalidFormat = errors.New("invalid ssdeep format")
)

// Distance computes the match score between two fuzzy hash signatures.
// Returns a value from zero to 100 indicating the match score of the two signatures.
// A match score of zero indicates the signatures did not match.
// Returns an error when one of the inputs are not valid signatures.
func Distance(hash1, hash2 string) (int, error) {
	var score int
	hash1BlockSize, hash1String1, hash1String2, err := splitSsdeep(hash1)
	if err != nil {
		return score, err
	}
	hash2BlockSize, hash2String1, hash2String2, err := splitSsdeep(hash2)
	if err != nil {
		return score, err
	}

	hash1String1 = eliminateSequences(hash1String1)
	hash1String2 = eliminateSequences(hash1String2)
	hash2String1 = eliminateSequences(hash2String1)
	hash2String2 = eliminateSequences(hash2String2)

	if hash1BlockSize == hash2BlockSize && hash1String1 == hash2String1 && hash1String2 == hash2String2 {
		return 100, nil
	}

	// We can only compare equal or *2 block sizes
	if hash1BlockSize != hash2BlockSize && hash1BlockSize != hash2BlockSize*2 && hash2BlockSize != hash1BlockSize*2 {
		return score, err
	}

	if hash1BlockSize == hash2BlockSize {
		d1 := scoreDistance(hash1String1, hash2String1, hash1BlockSize)
		d2 := scoreDistance(hash1String2, hash2String2, hash1BlockSize*2)
		score = int(math.Max(float64(d1), float64(d2)))
	} else if hash1BlockSize == hash2BlockSize*2 {
		score = scoreDistance(hash1String1, hash2String2, hash1BlockSize)
	} else {
		score = scoreDistance(hash1String2, hash2String1, hash2BlockSize)
	}
	return score, nil
}

func splitSsdeep(hash string) (int, string, string, error) {
	if hash == "" {
		return 0, "", "", ErrEmptyHash
	}

	parts := strings.Split(hash, ":")
	if len(parts) != 3 {
		return 0, "", "", ErrInvalidFormat
	}

	blockSize, err := strconv.Atoi(parts[0])
	if err != nil {
		return blockSize, "", "", fmt.Errorf("%s: %w", ErrInvalidFormat.Error(), err)
	}

	return blockSize, parts[1], parts[2], nil
}

// eliminateSequences keeps at most three consecutive copies of a character,
// as ssdeep does before it compares two hash parts.
func eliminateSequences(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if i < 3 || s[i] != s[i-1] || s[i] != s[i-2] || s[i] != s[i-3] {
			out = append(out, s[i])
		}
	}
	return string(out)
}

func scoreDistance(h1, h2 string, blockSize int) int {
	if !hasCommonSubstring(h1, h2) {
		return 0
	}
	d := distance(h1, h2)
	d = (d * spamSumLength) / (len(h1) + len(h2))
	d = (100 * d) / spamSumLength
	d = 100 - d

	// A short hash of a small block size matches too easily. Below a block
	// size of 45, the score cannot be higher than the block size in units of
	// blockMin, times the length of the shorter hash (as in ssdeep 2.14.1).
	// ssdeep reads the block size as an unsigned number, so a negative block
	// size gets no cap.
	if blockSize < 0 || blockSize >= (99+rollingWindow)/rollingWindow*blockMin {
		return d
	}
	shorter := len(h1)
	if len(h2) < shorter {
		shorter = len(h2)
	}
	if matchSize := blockSize / blockMin * shorter; d > matchSize {
		d = matchSize
	}
	return d
}
func hasCommonSubstring(s1, s2 string) bool {
	i := 0
	j := 0
	s1Len := len(s1)
	s2Len := len(s2)
	hashes := make([]uint32, (spamSumLength - (rollingWindow - 1)))
	if s1Len > spamSumLength || s2Len > spamSumLength {
		return false
	}
	if s1Len < rollingWindow || s2Len < rollingWindow {
		return false
	}
	state := &rollingState{}
	for i = 0; i < rollingWindow-1; i++ {
		state.rollHash(s1[i])
	}
	for i = rollingWindow - 1; i < s1Len; i++ {
		state.rollHash(s1[i])
		hashes[i-(rollingWindow-1)] = state.rollSum()
	}
	s1Len -= (rollingWindow - 1)
	state.rollReset()
	for j = 0; j < rollingWindow-1; j++ {
		state.rollHash(s2[j])
	}
	for j = 0; j < s2Len-(rollingWindow-1); j++ {
		state.rollHash(s2[j+(rollingWindow-1)])
		var h = state.rollSum()
		for i = 0; i < s1Len; i++ {
			if hashes[i] == h && s1[i:i+rollingWindow] == s2[j:j+rollingWindow] {
				return true
			}
		}
	}
	return false
}
