package ssdeep

import "hash"

var _ hash.Hash = &ssdeepState{}

// New instance of the SSDEEP hash
func New() *ssdeepState {
	s := newSSDEEPState()
	return &s
}

// Sum appends the current hash state to b.
func (state *ssdeepState) Sum(b []byte) []byte {
	digest, _ := state.digest()
	return append(b, digest...)
}

// BlockSize returns the acceptable minimum amount of data
func (state *ssdeepState) BlockSize() int {
	return blockMin
}

// Size of the hash to be returned
func (state *ssdeepState) Size() int {
	return spamSumLength
}

// Reset the hash to initial state
func (state *ssdeepState) Reset() {
	*state = newSSDEEPState()
}
