package ssdeep

import (
	"testing"

	cssdeep "github.com/dutchcoders/gossdeep"
)

func TestHashDistanceSame(t *testing.T) {
	d := Distance(&h1, &h1)
	if d != 100 {
		t.Errorf("Invalid hash distance: %d", d)
	}
}

func TestHashDistance(t *testing.T) {
	d := Distance(&h1, &h2)
	if d != 35 {
		t.Errorf("Invalid hash distance: %d", d)
	}
}

func TestNilHash1(t *testing.T) {
	d := Distance(nil, &h2)
	if d != 0 {
		t.Errorf("hash1 is nil: %d", d)
	}
}

func TestNilHash2(t *testing.T) {
	d := Distance(&h1, nil)
	if d != 0 {
		t.Errorf("hash2 is nil: %d", d)
	}
}

func TestNilHashes(t *testing.T) {
	d := Distance(nil, nil)
	if d != 0 {
		t.Errorf("hash1 and hash2 are nil: %d", d)
	}
}

func TestHashDistanceNoName(t *testing.T) {
	d := Distance(&h3, &h4)
	if d != 97 {
		t.Errorf("Invalid hash distance: %d", d)
	}
}

func BenchmarkDistance(b *testing.B) {
	var h1 = `7DSC8olnoL1v/uawvbQD7XlZUFYzYyMb615NktYHF7dREN/JNnQrmhnUPI+/n2Y7`
	var h2 = `7DSC8olnoL1v/uawvbQD7XlZUFYzYyMb615NktYHF7dREN/JNnQrmhnUPI+/ngrr`
	for i := 0; i < b.N; i++ {
		distance(h1, h2)
	}
}

func BenchmarkDistanceC(b *testing.B) {
	var h1 = `7DSC8olnoL1v/uawvbQD7XlZUFYzYyMb615NktYHF7dREN/JNnQrmhnUPI+/n2Y7`
	var h2 = `7DSC8olnoL1v/uawvbQD7XlZUFYzYyMb615NktYHF7dREN/JNnQrmhnUPI+/ngrr`
	for i := 0; i < b.N; i++ {
		cssdeep.Compare(h1, h2)
	}
}
