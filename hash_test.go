package ssdeep

import (
	"fmt"
	"hash"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashinterface(t *testing.T) {
	h := New()
	var _ hash.Hash = h
	t.Log(h.BlockSize())
	h.Reset()

	b, err := ioutil.ReadFile("ssdeep_results.json")
	assert.NoError(t, err)

	n, err := h.Write(b)
	require.NoError(t, err)

	t.Log(n)
	t.Log(h.Size())

	hashResult := h.Sum(nil)

	expectedResult := "1536:74peLhFipssVfuInITTTZzMoW0379xy3u:VVFosEfudTj579k3u"
	require.Equal(t, expectedResult, string(hashResult))

	t.Log(hashResult)
	t.Log(fmt.Sprintf("%x", hashResult[:]))
}

func TestHashWrite(t *testing.T) {
	// hash using the hash.Hash interface methods
	b, err := ioutil.ReadFile("ssdeep_results.json")
	assert.NoError(t, err)

	h1 := New()
	h1.Write([]byte("1234"))
	h1.Write(b)
	t.Log(fmt.Sprintf("h1: %x", h1.Sum(nil)))

	// hash from read
	h2, err := FuzzyBytes(append([]byte("1234"), b...))
	require.NoError(t, err)
	t.Log(fmt.Sprintf("h2: %s", h2))

	// compare hashes
	diff := distance(string(h1.Sum(nil)), h2)
	if diff != 0 {
		t.Errorf("hashes differ by: %d", diff)
	}
}
