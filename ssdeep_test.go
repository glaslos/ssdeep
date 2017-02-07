package ssdeep

import (
	"log"
	"math/rand"
	"os/exec"
	"testing"
)

func TestRollingHash(t *testing.T) {
	sdeep := NewSSDEEP()
	if sdeep.rollHash(byte('A')) != 585 {
		t.Error("Rolling hash not matching")
	}
}

func TestFindBlock(t *testing.T) {
	sdeep := NewSSDEEP()
	sdeep.Fuzzy("/tmp/data")
	_, err := exec.Command("ssdeep", "/tmp/data").Output()
	if err != nil {
		log.Fatal(err)
	}
}

func BenchmarkRollingHash(b *testing.B) {
	sdeep := NewSSDEEP()
	for i := 0; i < b.N; i++ {
		sdeep.rollHash(byte(i))
	}
}

func BenchmarkDistance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		distance("aaa123123123123aaaa", "aab123123123aaaa")
	}
}

func BenchmarkSumHash(b *testing.B) {
	testHash := hashIinit
	data := []byte("Hereyougojustsomedatatomakeyouhappy")
	for i := 0; i < b.N; i++ {
		testHash = sumHash(data[rand.Intn(len(data))], testHash)
	}
}
