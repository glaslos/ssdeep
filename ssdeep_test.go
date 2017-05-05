package ssdeep

import (
	"bufio"
	"io/ioutil"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"testing"
)

var h1 = FuzzyHash{
	blockSize:   192,
	hashString1: "MUPMinqP6+wNQ7Q40L/iB3n2rIBrP0GZKF4jsef+0FVQLSwbLbj41iH8nFVYv980",
	hashString2: "x0CllivQiFmt",
}
var h2 = FuzzyHash{
	blockSize:   192,
	hashString1: "JkjRcePWsNVQza3ntZStn5VfsoXMhRD9+xJMinqF6+wNQ7Q40L/i737rPVt",
	hashString2: "JkjlQyIrx+kll2",
}
var h3 = FuzzyHash{
	blockSize:   196608,
	hashString1: "pDSC8olnoL1v/uawvbQD7XlZUFYzYyMb615NktYHF7dREN/JNnQrmhnUPI+/n2Yr",
	hashString2: "5DHoJXv7XOq7Mb2TwYHXREN/3QrmktPd",
}
var h4 = FuzzyHash{
	blockSize:   196608,
	hashString1: "7DSC8olnoL1v/uawvbQD7XlZUFYzYyMb615NktYHF7dREN/JNnQrmhnUPI+/n2Y7",
	hashString2: "3DHoJXv7XOq7Mb2TwYHXREN/3QrmktPt",
}
var h5 = FuzzyHash{
	blockSize:   196608,
	hashString1: "7DSC8olnoL1v/uawvbQD7XlZUFYzYyMb615NktYHF7dREN/JNnQrmhnUPI+/n2Y7",
	hashString2: "",
}

func TestRollingHash(t *testing.T) {
	sdeep := SSDEEP{
		rollingState: rollingState{
			window: make([]byte, rollingWindow),
		},
	}
	if sdeep.rollHash(byte('A')) != 585 {
		t.Error("Rolling hash not matching")
	}
}

func TestCompareHashFile(t *testing.T) {
	sdeep := NewSSDEEP()
	libhash, err := sdeep.Fuzzy("/tmp/data")
	if err != nil {
		t.Error(err)
	}
	out, err := exec.Command("ssdeep", "/tmp/data").Output()
	if err != nil {
		log.Fatal(err)
	}
	if strings.Split(string(out[:]), "\n")[1] != libhash.String() {
		t.Error("Hash mismatch")
	}
}

func TestEmptyByte(t *testing.T) {
	sdeep := NewSSDEEP()
	_, err := sdeep.FuzzyByte([]byte{})
	if err == nil {
		t.Error("Expecting error for empty file")
	}
}

func TestEmptyFile(t *testing.T) {
	f, err := ioutil.TempFile("/tmp", "ssdeep-test")
	if err != nil {
		log.Fatal(err)
	}
	sdeep := NewSSDEEP()
	r := bufio.NewReader(f)
	_, err = sdeep.FuzzyReader(r, f.Name())
	if err == nil {
		t.Error("Expecting error for empty file")
	}
}

func TestCompareHashBytes(t *testing.T) {
	blob, err := ioutil.ReadFile("/tmp/data")
	if err != nil {
		t.Error(err)
	}
	sdeep := NewSSDEEP()
	libhash, err := sdeep.FuzzyByte(blob)
	if err != nil {
		t.Error(err)
	}
	out, err := exec.Command("ssdeep", "/tmp/data").Output()
	if err != nil {
		t.Error(err)
	}
	if strings.Split(string(out[:]), "\n")[1] != libhash.String()+",\"/tmp/data\"" {
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

func BenchmarkBlockSize(b *testing.B) {
	sdeep := NewSSDEEP()
	for i := 0; i < b.N; i++ {
		sdeep.getBlockSize(207160)
	}
}
