package ssdeep

import "math"

func distance(sample1, sample2 string) int {
	var cost0 [spamSumLength + 1]int
	var cost1 [spamSumLength + 1]int

	for i := 0; i <= len(sample2); i++ {
		cost0[i] = i
	}

	for i := 0; i < len(sample1); i++ {
		cost1[0] = i + 1
		for j := 0; j < len(sample2); j++ {
			// Insert
			var costI = cost0[j+1] + 1
			// Delete
			var costD = cost1[j] + 1
			// Replace
			var costR int
			if sample1[i] == sample2[j] {
				costR = cost0[j]
			} else {
				costR = cost0[j] + 2
			}
			cost1[j+1] = int(math.Min(math.Min(float64(costI), float64(costD)), float64(costR)))
		}
		cost0, cost1 = cost1, cost0
	}
	return cost0[len(sample2)]
}
