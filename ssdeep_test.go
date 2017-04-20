package ssdeep

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"testing"
)

func TestRollingHash(t *testing.T) {
	sdeep := NewSSDEEP()
	if sdeep.rollHash(byte('A')) != 585 {
		t.Error("Rolling hash not matching")
	}
}

func TestCompareHashFile(t *testing.T) {
	sdeep := NewSSDEEP()
	libhash := sdeep.Fuzzy("/tmp/data")
	out, err := exec.Command("ssdeep", "/tmp/data").Output()
	if err != nil {
		log.Fatal(err)
	}
	if strings.Split(string(out[:]), "\n")[1] != libhash {
		t.Error("Hash mismatch")
	}
}

func TestCompareHashBytes(t *testing.T) {
	blob, err := ioutil.ReadFile("/tmp/data")
	if err != nil {
		t.Error(err)
	}
	sdeep := NewSSDEEP()
	libhash := sdeep.FuzzyByte(blob)
	out, err := exec.Command("ssdeep", "/tmp/data").Output()
	if err != nil {
		t.Error(err)
	}
	if strings.Split(string(out[:]), "\n")[1] != libhash+",\"/tmp/data\"" {
		t.Error("Hash mismatch")
	}
}

func BenchmarkRollingHash(b *testing.B) {
	sdeep := NewSSDEEP()
	for i := 0; i < b.N; i++ {
		sdeep.rollHash(byte(i))
	}
}

func BenchmarkSumHash(b *testing.B) {
	testHash := hashIinit
	data := []byte("Hereyougojustsomedatatomakeyouhappy")
	for i := 0; i < b.N; i++ {
		testHash = sumHash(data[rand.Intn(len(data))], testHash)
	}
}
