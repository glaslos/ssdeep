package ssdeep

import "hash"

var _ hash.Hash = &ssdeepState{}

func New() *ssdeepState {
	s := newSSDEEPState()
	return &s
}

func (state *ssdeepState) Sum(b []byte) []byte {
	digest, _ := state.digest()
	return append(b, digest...)
}

func (state *ssdeepState) BlockSize() int {
	return blockMin
}

func (state *ssdeepState) Size() int {
	return spamSumLength
}

func (state *ssdeepState) Reset() {
	*state = newSSDEEPState()
}
