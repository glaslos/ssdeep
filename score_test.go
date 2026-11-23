package ssdeep

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var h1 = "192:MUPMinqP6+wNQ7Q40L/iB3n2rIBrP0GZKF4jsef+0FVQLSwbLbj41iH8nFVYv980:x0CllivQiFmt"

var h2 = "192:JkjRcePWsNVQza3ntZStn5VfsoXMhRD9+xJMinqF6+wNQ7Q40L/i737rPVt:JkjlQyIrx+kll2"

var h3 = "196608:pDSC8olnoL1v/uawvbQD7XlZUFYzYyMb615NktYHF7dREN/JNnQrmhnUPI+/n2Yr:5DHoJXv7XOq7Mb2TwYHXREN/3QrmktPd"

var h4 = "196608:7DSC8olnoL1v/uawvbQD7XlZUFYzYyMb615NktYHF7dREN/JNnQrmhnUPI+/n2Y7:3DHoJXv7XOq7Mb2TwYHXREN/3QrmktPt"

var h5 = "24:YDVLfsT1ds/1H9Wpgq7n4XMijV6h4Z3QCw4qat:YD51H9CiMuV6uACwVat"

var h6 = "24:YDVLfyvDj+C+opg8DV0Mdle6hPZ3QCw4qat:YDMvDj+C+kBOM+6HACwVat"

func assertDistanceEqual(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Fatalf("Distance mismatch: %d (expected)\n"+
			"                != %d (actual)", expected, actual)
	}
}

func TestHashDistanceSame(t *testing.T) {
	d, err := Distance(h1, h1)
	require.NoError(t, err)
	assertDistanceEqual(t, 100, d)
}

func TestHashDistance1(t *testing.T) {
	d, err := Distance(h1, h2)
	require.NoError(t, err)
	assertDistanceEqual(t, 35, d)
}

func TestHashDistance2(t *testing.T) {
	d, err := Distance(h3, h4)
	require.NoError(t, err)
	assertDistanceEqual(t, 97, d)
}

func TestHashDistance3(t *testing.T) {
	d, err := Distance(h5, h6)
	require.NoError(t, err)
	assertDistanceEqual(t, 54, d)
}

// The expected scores come from "ssdeep -a -k" of ssdeep 2.14.1.
func TestDistanceMatchesSsdeep(t *testing.T) {
	tests := []struct {
		name  string
		hash1 string
		hash2 string
		want  int
	}{
		{"cap at block size 3", "3:ABCDEFGHIJKLMNOPQRST:abcdefghij", "3:ABCDEFGHIJKLMNOPQRSX:zyxwvutsrq", 20},
		{"cap at block size 24", "24:ABCDEFGHIJ:abcdefghij", "24:ABCDEFGHIX:zyxwvutsrq", 80},
		{"no cap at block size 96", "96:ABCDEFGHIJKLMNOPQRSTUVWXYZabcdef:ghijklmnop", "96:ABCDEFGHIJKLMNOPQRSTUVWXYZabcdeX:zyxwvutsrq", 97},
		{"cap at double block size", "6:ABCDEFGHIJKLMNOPQRST:abcdefghijklmnopqrst", "12:abcdefghijklmnopqrsX:ZYXWVUTSRQ", 80},
		{"cap at half block size", "12:abcdefghijklmnopqrsX:ZYXWVUTSRQ", "6:ABCDEFGHIJKLMNOPQRST:abcdefghijklmnopqrst", 80},
		{"runs of more than three characters", "192:XYZXYZAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABCDEFGHIJKLMNOP:qq", "192:XYZXYZAAAABCDEFGHIJKLMNOP:rr", 100},
		{"identical after the sequence elimination", "3:AAAAAABCDEFGHIJ:abc", "3:AAABCDEFGHIJ:abc", 100},
		{"same first part, other second part", "3:ABCDEFGHIJ:abc", "3:ABCDEFGHIJ:xyz", 10},
		{"same short first part, other second part", "48:ABC:abc", "48:ABC:xyz", 0},
		{"identical short signatures", "3:ABC:abc", "3:ABC:abc", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := Distance(tt.hash1, tt.hash2)
			require.NoError(t, err)
			assertDistanceEqual(t, tt.want, d)
		})
	}
}

func TestEmptyHash1(t *testing.T) {
	d, err := Distance("", h2)
	require.Error(t, err)
	assertDistanceEqual(t, 0, d)
}

func TestEmptyHash2(t *testing.T) {
	d, err := Distance(h1, "")
	require.Error(t, err)
	if d != 0 {
		t.Errorf("hash2 is nil: %d", d)
	}
}

func TestInvalidHash1(t *testing.T) {
	d, err := Distance("192:asdasd", h1)
	require.Error(t, err)

	if d != 0 {
		t.Errorf("hash1 and hash2 are nil: %d", d)
	}
}

func TestInvalidHash2(t *testing.T) {
	d, err := Distance(h1, "asd:asdasd:aaaa")
	require.Error(t, err)
	if d != 0 {
		t.Errorf("hash1 and hash2 are nil: %d", d)
	}
}

// Crash seed: body length 65 triggers index out of range [58] with length 58.
func TestHasCommonSubstringOOB(t *testing.T) {
	s1 := "3:" + strings.Repeat("0123456789", 7)[:65] + ":"
	s2 := "3:0123456:"
	// Must not panic.
	_, err := Distance(s1, s2)
	if err != nil {
		t.Logf("Distance returned error (acceptable): %v", err)
	}
}

// FuzzSSDeepDistanceDirect feeds arbitrary hash strings directly to Distance.
// The library must not panic; an error return is acceptable.
func FuzzSSDeepDistanceDirect(f *testing.F) {
	f.Add("3:"+strings.Repeat("0123456789", 7)[:65]+":", "3:0123456:")
	f.Add("3:abc:abc", "3:abc:abc")
	f.Add("", "")
	f.Add("not-a-hash", "also-not-a-hash")

	f.Fuzz(func(t *testing.T, s1, s2 string) {
		// Must not panic.
		_, _ = Distance(s1, s2)
	})
}

func BenchmarkDistance(b *testing.B) {
	var h1 = `7DSC8olnoL1v/uawvbQD7XlZUFYzYyMb615NktYHF7dREN/JNnQrmhnUPI+/n2Y7`
	var h2 = `7DSC8olnoL1v/uawvbQD7XlZUFYzYyMb615NktYHF7dREN/JNnQrmhnUPI+/ngrr`
	for i := 0; i < b.N; i++ {
		distance(h1, h2)
	}
}
