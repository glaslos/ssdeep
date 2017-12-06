package ssdeep

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"
)

func readFile(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(f)
	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, r)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func concatCopyPreAllocate(slices [][]byte) []byte {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	tmp := make([]byte, totalLen)
	var i int
	for _, s := range slices {
		i += copy(tmp[i:], s)
	}
	return tmp
}

func TestRollingHash(t *testing.T) {
	s := SSDEEP{
		rollingState: rollingState{
			window: make([]byte, rollingWindow),
		},
	}
	if s.rollHash(byte('A')) != 585 {
		t.Error("Rolling hash not matching")
	}
}

func TestCompareHashFile(t *testing.T) {
	s := NewSSDEEP()
	b, err := readFile("LICENSE")
	if err != nil {
		t.Error(err)
	}
	b = concatCopyPreAllocate([][]byte{b, b})
	libhash, err := s.FuzzyByte(b)
	if err != nil {
		t.Error(err)
	}
	expectedResult := "96:PuNQHTo6pYrYJWrYJ6N3w53hpYTdhuNQHTo6pYrYJWrYJ6N3w53hpYTP:+QHTrpYrsWrs6N3g3LaGQHTrpYrsWrsa"
	if libhash.String() != expectedResult {
		t.Errorf(
			"Hash mismatch: %s vs %s", libhash.String(), expectedResult,
		)
	}
}

func TestFuzzyReaderError(t *testing.T) {
	s := NewSSDEEP()
	b, err := readFile("ssdeep.go")
	if err != nil {
		t.Error(err)
	}
	r := bytes.NewReader(b)
	if _, err := s.FuzzyReader(r, "ssdeep.go"); err == nil {
		t.Error("Expecting error for missing block size")
	}
}

func TestFuzzyReader(t *testing.T) {
	s := NewSSDEEP()
	b, err := readFile("ssdeep.go")
	if err != nil {
		t.Error(err)
	}
	s.GetBlockSize(len(b))
	r := bytes.NewReader(b)
	h, err := s.FuzzyReader(r, "ssdeep.go")
	if err != nil {
		t.Error(err)
	}
	if r := h.String(); r == "" {
		t.Error("No hash string returned")
	}
}

func TestFuzzyFile(t *testing.T) {
	s := NewSSDEEP()
	f, err := os.Open("ssdeep.go")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	_, err = s.FuzzyFile(f, "ssdeep.go")
	if err != nil {
		t.Error(err)
	}
}

func TestEmptyByte(t *testing.T) {
	s := NewSSDEEP()
	_, err := s.FuzzyByte([]byte{})
	if err == nil {
		t.Error("Expecting error for empty file")
	}
}

func BenchmarkRollingHash(b *testing.B) {
	s := NewSSDEEP()
	for i := 0; i < b.N; i++ {
		s.rollHash(byte(i))
	}
}

func BenchmarkSumHash(b *testing.B) {
	testHash := hashInit
	data := []byte("Hereyougojustsomedatatomakeyouhappy")
	for i := 0; i < b.N; i++ {
		testHash = sumHash(data[rand.Intn(len(data))], testHash)
	}
}

func BenchmarkBlockSize(b *testing.B) {
	s := NewSSDEEP()
	for i := 0; i < b.N; i++ {
		s.GetBlockSize(207160)
	}
}
