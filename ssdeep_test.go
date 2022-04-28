package ssdeep

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertHashEqual(t *testing.T, expected, actual string) {
	require.Equal(t, expected, actual, "hash mismatch")
}

func TestIntegrity(t *testing.T) {
	rand.Seed(1)
	resultsFile, err := ioutil.ReadFile("ssdeep_results.json")
	assert.NoError(t, err)

	originalResults := make(map[string]string)
	err = json.Unmarshal(resultsFile, &originalResults)
	assert.NoError(t, err)

	for i := 4097; i < 10*1024*1024; i += 4096 * 10 {
		t.Run(fmt.Sprintf("Bytes in size of %d", i), func(t *testing.T) {
			size := i
			if size == 4097 {
				i--
			}
			blob := make([]byte, size)
			_, err = rand.Read(blob)
			assert.NoError(t, err)
			result, err := FuzzyBytes(blob)
			assert.NoError(t, err)
			assertHashEqual(t, originalResults[fmt.Sprint(size)], result)
		})
	}
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
	s := ssdeepState{
		rState: rollingState{
			window: make([]byte, rollingWindow),
		},
	}
	s.rollHash(byte('A'))
	rh := s.rState.rollSum()
	if rh != 585 {
		t.Fatal("Rolling hash not matching")
	}
}

func TestFuzzyBytesOutputsTheRightResult(t *testing.T) {
	b, err := ioutil.ReadFile("LICENSE")
	assert.NoError(t, err)

	b = concatCopyPreAllocate([][]byte{b, b})
	hashResult, err := FuzzyBytes(b)
	assert.NoError(t, err)

	expectedResult := "96:PuNQHTo6pYrYJWrYJ6N3w53hpYTdhuNQHTo6pYrYJWrYJ6N3w53hpYTP:+QHTrpYrsWrs6N3g3LaGQHTrpYrsWrsa"
	assertHashEqual(t, expectedResult, hashResult)
}

func TestFuzzyFileOutputsTheRightResult(t *testing.T) {
	f, err := os.Open("ssdeep_results.json")
	assert.NoError(t, err)
	defer f.Close()

	hashResult, err := FuzzyFile(f)
	assert.NoError(t, err)

	expectedResult := "1536:74peLhFipssVfuInITTTZzMoW0379xy3u:VVFosEfudTj579k3u"
	assertHashEqual(t, expectedResult, hashResult)

}

func TestFuzzyFileOutputsAnErrorForSmallFiles(t *testing.T) {
	f, err := os.Open("LICENSE")
	assert.NoError(t, err)
	defer f.Close()

	_, err = FuzzyFile(f)
	assert.Error(t, err)
}

func TestFuzzyFilenameOutputsTheRightResult(t *testing.T) {
	hashResult, err := FuzzyFilename("ssdeep_results.json")
	assert.NoError(t, err)

	expectedResult := "1536:74peLhFipssVfuInITTTZzMoW0379xy3u:VVFosEfudTj579k3u"
	assertHashEqual(t, expectedResult, hashResult)
}

func TestFuzzyFilenameOutputsErrorWhenFileNotExists(t *testing.T) {
	_, err := FuzzyFilename("foo.bar")
	assert.Error(t, err)
}

func TestFuzzyBytesWithLenLessThanMinimumOutputsAnError(t *testing.T) {
	_, err := FuzzyBytes([]byte{})
	assert.Error(t, err)
}

func TestFuzzyBytesWithOutputsAnError(t *testing.T) {
	_, err := FuzzyBytes(make([]byte, 4096))
	assert.Error(t, err)
}

func BenchmarkRollingHash(b *testing.B) {
	s := newSsdeepState()
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
	s := newSsdeepState()
	for i := 0; i < b.N; i++ {
		s.setBlockSize(207160)
	}
}

func BenchmarkProcessByte(b *testing.B) {
	s := newSsdeepState()
	s.blockSize = 42
	s.newRollingState()
	for i := 0; i < b.N; i++ {
		s.processByte(byte(i))
	}
}
