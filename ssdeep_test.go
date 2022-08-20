package ssdeep

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertHashEqual(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Fatalf("Hash mismatch: %s (expected)\n"+
			"            != %s (actual)", expected, actual)
	}
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
	s := rollingState{}
	s.rollHash(byte('A'))
	rh := s.rollSum()
	require.Equal(t, rh, uint32(585), "Rolling hash not matching")
}

func TestFuzzyHashOutputsTheRightResult(t *testing.T) {
	b, err := ioutil.ReadFile("LICENSE")
	require.NoError(t, err)

	b = concatCopyPreAllocate([][]byte{b, b})
	s := New()

	_, err = io.Copy(s, bytes.NewReader(b))
	require.NoError(t, err)

	expectedResult := "96:PuNQHTo6pYrYJWrYJ6N3w53hpYTdhuNQHTo6pYrYJWrYJ6N3w53hpYTP:+QHTrpYrsWrs6N3g3LaGQHTrpYrsWrsa"
	prepend := []byte("prepend")

	sum := s.Sum(prepend)

	assertHashEqual(t, string(append(prepend, expectedResult...)), string(sum))
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
		s.rollingState.rollHash(byte(i))
	}
}

func BenchmarkSumHash(b *testing.B) {
	var testHash byte = hashInit
	data := []byte("Hereyougojustsomedatatomakeyouhappy")
	for i := 0; i < b.N; i++ {
		testHash = sumHash(data[rand.Intn(len(data))], testHash)
	}
}

func BenchmarkProcessByte(b *testing.B) {
	s := newSsdeepState()
	for i := 0; i < b.N; i++ {
		s.processByte(byte(i))
	}
}
