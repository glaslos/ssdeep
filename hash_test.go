package ssdeep

import (
	"hash"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashinterface(t *testing.T) {
	h := New()
	var _ hash.Hash = h
	t.Log(h.BlockSize())
	h.Reset()

	fh, err := os.Open("ssdeep_results.json")
	require.NoError(t, err)
	b, err := io.ReadAll(fh)
	require.NoError(t, err)

	n, err := h.Write(b)
	require.NoError(t, err)

	t.Log(n)
	t.Log(h.Size())

	hashResult := h.Sum(nil)

	expectedResult := "1536:74peLhFipssVfuInITTTZzMoW0379xy3u:VVFosEfudTj579k3u"
	require.Equal(t, expectedResult, string(hashResult))

	t.Log(hashResult)
	t.Logf("%x", hashResult[:])
}

func TestHashWrite(t *testing.T) {
	// hash using the hash.Hash interface methods
	fh, err := os.Open("ssdeep_results.json")
	require.NoError(t, err)
	b, err := io.ReadAll(fh)
	require.NoError(t, err)

	h1 := New()
	h1.Write([]byte("1234"))
	h1.Write(b)
	t.Logf("h1: %x", h1.Sum(nil))

	// hash from read
	h2, err := FuzzyBytes(append([]byte("1234"), b...))
	require.NoError(t, err)
	t.Logf("h2: %s", h2)

	// compare hashes
	diff := distance(string(h1.Sum(nil)), h2)
	if diff != 0 {
		t.Errorf("hashes differ by: %d", diff)
	}
}
