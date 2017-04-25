package ssdeep

import "math"

// Distance between two strings
func Distance(hash1, hash2 *FuzzyHash) (score int) {
	if hash1 == nil || hash2 == nil {
		return 0
	}
	// We can only compare equal or *2 block sizes
	if hash1.blockSize != hash2.blockSize && hash1.blockSize != hash2.blockSize*2 && hash2.blockSize != hash1.blockSize*2 {
		return
	}
	if hash1.blockSize == hash2.blockSize && hash1.hashString1 == hash2.hashString1 {
		return 100
	}
	if hash1.blockSize == hash2.blockSize {
		d1 := scoreDistance(hash1.hashString1, hash2.hashString1, hash1.blockSize)
		d2 := scoreDistance(hash1.hashString2, hash2.hashString2, hash1.blockSize*2)
		score = int(math.Max(float64(d1), float64(d2)))
	} else if hash1.blockSize == hash2.blockSize*2 {
		score = scoreDistance(hash1.hashString1, hash2.hashString2, hash1.blockSize)
	} else {
		score = scoreDistance(hash1.hashString2, hash2.hashString1, hash2.blockSize)
	}
	return
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
