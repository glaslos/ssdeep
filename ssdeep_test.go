package ssdeep

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrity(t *testing.T) {
	rand.Seed(1)

	fh, err := os.Open("ssdeep_results.json")
	require.NoError(t, err)

	originalResults := make(map[string]string)
	err = json.NewDecoder(fh).Decode(&originalResults)
	require.NoError(t, err)

	for i := 4097; i < 10*1024*1024; i += 4096 * 10 {
		t.Run(fmt.Sprintf("Bytes in size of %d", i), func(t *testing.T) {
			size := i
			if size == 4097 {
				i--
			}
			blob := make([]byte, size)
			_, err = rand.Read(blob)
			require.NoError(t, err)
			result, err := FuzzyBytes(blob)
			require.NoError(t, err)
			require.Equal(t, originalResults[fmt.Sprint(size)], result)
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
	fh, err := os.Open("LICENSE")
	require.NoError(t, err)
	b, err := io.ReadAll(fh)
	require.NoError(t, err)

	b = concatCopyPreAllocate([][]byte{b, b})
	s := New()

	_, err = io.Copy(s, bytes.NewReader(b))
	require.NoError(t, err)

	expectedResult := "96:PuNQHTo6pYrYJWrYJ6N3w53hpYTdhuNQHTo6pYrYJWrYJ6N3w53hpYTP:+QHTrpYrsWrs6N3g3LaGQHTrpYrsWrsa"
	prepend := []byte("prepend")

	sum := s.Sum(prepend)

	require.Equal(t, string(append(prepend, expectedResult...)), string(sum))
}

func TestFuzzyBytesOutputsTheRightResult(t *testing.T) {
	fh, err := os.Open("LICENSE")
	require.NoError(t, err)
	b, err := io.ReadAll(fh)
	require.NoError(t, err)

	b = concatCopyPreAllocate([][]byte{b, b})
	hashResult, err := FuzzyBytes(b)
	require.NoError(t, err)

	expectedResult := "96:PuNQHTo6pYrYJWrYJ6N3w53hpYTdhuNQHTo6pYrYJWrYJ6N3w53hpYTP:+QHTrpYrsWrs6N3g3LaGQHTrpYrsWrsa"
	require.Equal(t, expectedResult, hashResult)
}

func TestFuzzyFileOutputsTheRightResult(t *testing.T) {
	f, err := os.Open("ssdeep_results.json")
	require.NoError(t, err)
	defer f.Close()

	hashResult, err := FuzzyFile(f)
	require.NoError(t, err)

	expectedResult := "1536:74peLhFipssVfuInITTTZzMoW0379xy3u:VVFosEfudTj579k3u"
	require.Equal(t, expectedResult, hashResult)

}

func TestFuzzyFileOutputsAnErrorForSmallFiles(t *testing.T) {
	f, err := os.Open("LICENSE")
	require.NoError(t, err)
	defer f.Close()

	_, err = FuzzyFile(f)
	require.Error(t, err)
}

func TestFuzzyFilenameOutputsTheRightResult(t *testing.T) {
	hashResult, err := FuzzyFilename("ssdeep_results.json")
	require.NoError(t, err)

	expectedResult := "1536:74peLhFipssVfuInITTTZzMoW0379xy3u:VVFosEfudTj579k3u"
	require.Equal(t, expectedResult, hashResult)
}

func TestFuzzyFilenameOutputsErrorWhenFileNotExists(t *testing.T) {
	_, err := FuzzyFilename("foo.bar")
	require.Error(t, err)
}

func TestFuzzyBytesWithLenLessThanMinimumOutputsAnError(t *testing.T) {
	_, err := FuzzyBytes([]byte{})
	require.Error(t, err)
}

func TestFuzzyBytesWithOutputsAnError(t *testing.T) {
	_, err := FuzzyBytes(make([]byte, 4096))
	require.Error(t, err)
}

func BenchmarkRollingHash(b *testing.B) {
	s := newSSDEEPState()
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
	s := newSSDEEPState()
	for i := 0; i < b.N; i++ {
		s.processByte(byte(i))
	}
}

func TestDigestLargeFile(t *testing.T) {
	state := ssdeepState{
		rollingState: rollingState{
			window: [8]uint8{97, 97, 97, 97, 97, 97, 97, 0},
			h1:     679,
			h2:     2716,
			h3:     2216757313,
			n:      6,
		},
		iStart:    0,
		iEnd:      2,
		totalSize: 4500000000,
		bsizeMask: 0,
		blocks: [31]blockHashState{
			{hashString: []uint8{45, 35}, blockSize: 3, blockHash1: 53, blockHash2: 53, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 6, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 12, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 24, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 48, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 96, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 192, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 384, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 768, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 1536, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 3072, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 6144, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 12288, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 24576, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 49152, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 98304, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 196608, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 393216, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 786432, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 1572864, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 3145728, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 6291456, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 12582912, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 25165824, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 50331648, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 100663296, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 201326592, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 402653184, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 805306368, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 1610612736, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
			{hashString: nil, blockSize: 3221225472, blockHash1: 39, blockHash2: 39, tail1: 0, tail2: 0},
		},
	}
	digest, err := state.digest()
	require.NoError(t, err)
	require.Equal(t, "3:tj1:n", digest)
}
