package ssdeep

import "testing"

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
