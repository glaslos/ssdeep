package ssdeep

import (
	"fmt"
	"log"
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
	out, err := exec.Command("ssdeep", "/tmp/data").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}

func BenchmarkRollingHash(b *testing.B) {
	sdeep := NewSSDEEP()
	for i := 0; i < b.N; i++ {
		sdeep.rollHash(byte(i))
	}
}

func BenchmarkEditDistance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		levEditDistance("aaa123123123123aaaa", "aab123123123aaaa")
	}
}
