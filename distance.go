// Copyright (c) 2015, Arbo von Monkiewitsch All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssdeep

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

// HashDistance between two strings
func HashDistance(str1, str2 string) (int, error) {
	if str1 == "" || str2 == "" {
		return 0, errors.New("Invaild input. Need two len(hashes) > 0")
	}
	str1Split := strings.Split(str1, ":")
	if len(str1Split) < 3 {
		return 0, errors.New("Invalid first hash, need format len:hash:hash")
	}
	bsize1, hash11, hash12 := str1Split[0], str1Split[1], strings.Split(str1Split[2], ",")[0]
	blockSize1, err := strconv.Atoi(bsize1)
	if err != nil {
		return 0, err
	}
	str2Split := strings.Split(str2, ":")
	if len(str2Split) < 3 {
		return 0, errors.New("Invalid second hash, need format len:hash:hash")
	}
	bsize2, hash21, hash22 := str2Split[0], str2Split[1], strings.Split(str2Split[2], ",")[0]
	blockSize2, err := strconv.Atoi(bsize2)
	if err != nil {
		return 0, err
	}
	// We can only compare equal or *2 block sizes
	if blockSize1 != blockSize2 && blockSize1 != blockSize2*2 && blockSize2 != blockSize1*2 {
		return 0, errors.New("Apples != Grapes")
	}
	// TODO: remove char repetitions in hashes here as they skew the results
	// TODO: compare char by char to exit fast
	if blockSize1 == blockSize2 && hash11 == hash21 {
		return 100, nil
	}
	var score int
	if blockSize1 == blockSize2 {
		d1 := scoreDistance(hash11, hash21, blockSize1)
		d2 := scoreDistance(hash12, hash22, blockSize1*2)
		score = int(math.Max(float64(d1), float64(d2)))
	} else if blockSize1 == blockSize2*2 {
		score = scoreDistance(hash11, hash22, blockSize1)
	} else {
		score = scoreDistance(hash12, hash21, blockSize2)
	}
	return score, nil
}

func scoreDistance(h1, h2 string, blockSize int) int {
	d := distance(h1, h2)
	d = (d * spamSumLength) / (len(h1) + len(h2))
	d = (100 * d) / spamSumLength
	d = 100 - d
	/* TODO: Figure out this black magic...
	matchSize := float64(blockSize) / float64(blockMin) * math.Min(float64(len(h1)), float64(len(h2)))
	if d > int(matchSize) {
		d = int(matchSize)
	}
	*/
	return d
}

func distance(str1, str2 string) int {
	var cost, lastdiag, olddiag int
	s1 := []rune(str1)
	s2 := []rune(str2)

	lenS1 := len(s1)
	lenS2 := len(s2)

	column := make([]int, lenS1+1)

	for y := 1; y <= lenS1; y++ {
		column[y] = y
	}

	for x := 1; x <= lenS2; x++ {
		column[0] = x
		lastdiag = x - 1
		for y := 1; y <= lenS1; y++ {
			olddiag = column[y]
			cost = 0
			if s1[y-1] != s2[x-1] {
				cost = 1
			}
			column[y] = min(
				column[y]+1,
				column[y-1]+1,
				lastdiag+cost)
			lastdiag = olddiag
		}
	}
	return column[lenS1]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	} else {
		if b < c {
			return b
		}
	}
	return c
}
