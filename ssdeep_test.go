package ssdeep

import (
	"fmt"
	"log"
	"os/exec"
	"testing"
)

func TestRollingHash(t *testing.T) {
	sdeep := newSSDEEP()
	if sdeep.rollHash(byte('A')) != 585 {
		t.Error("Rolling hash not matching")
	}
}

func TestFindBlock(t *testing.T) {
	sdeep := newSSDEEP()
	sdeep.Fuzzy("/tmp/data")
	out, err := exec.Command("ssdeep", "/tmp/data").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}

func BenchmarkRollingHash(b *testing.B) {
	sdeep := newSSDEEP()
	for i := 0; i < b.N; i++ {
		sdeep.rollHash(byte(i))
	}
}
