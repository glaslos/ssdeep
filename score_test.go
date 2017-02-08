package ssdeep

import "testing"

var h1 = `196608:7DSC8olnoL1v/uawvbQD7XlZUFYzYyMb615NktYHF7dREN/JNnQrmhnUPI+/n2Y7:3DHoJXv7XOq7Mb2TwYHXREN/3QrmktPt,"/tmp/ssdeep/data"`
var h2 = `196608:pDSC8olnoL1v/uawvbQD7XlZUFYzYyMb615NktYHF7dREN/JNnQrmhnUPI+/n2Yr:5DHoJXv7XOq7Mb2TwYHXREN/3QrmktPd,"/tmp/ssdeep/data2"`
var h3 = `196607:pDSC8olnoL1v/uawvbQD7XlZUFYzYyMb615NktYHF7dREN/JNnQrmhnUPI+/n2Yr:5DHoJXv7XOq7Mb2TwYHXREN/3QrmktPd,"/tmp/ssdeep/data2"`
var h4 = `196608:7DSC8olnoL1v/uawvbQD7XlZUFYzYyMb615NktYHF7dREN/JNnQrmhnUPI+/n2Y7:3DHoJXv7XOq7Mb2TwYHXREN/3QrmktPt`
var h5 = `196608:7DSC8olnoL1v/uawvbQD7XlZUFYzYyMb615NktYHF7dREN/JNnQrmhnUPI+/n2Y7`

func TestHashDistanceSame(t *testing.T) {
	d, err := HashDistance(h1, h1)
	if d != 100 {
		t.Errorf("Invalid edit distance: %d", d)
	}
	if err != nil {
		t.Error(err)
	}
}

func TestHashDistance(t *testing.T) {
	d, err := HashDistance(h1, h2)
	if d != 97 {
		t.Errorf("Invalid edit distance: %d", d)
	}
	if err != nil {
		t.Error(err)
	}
}

func TestHashDistanceNoName(t *testing.T) {
	d, err := HashDistance(h2, h4)
	if d != 97 {
		t.Errorf("Invalid edit distance: %d", d)
	}
	if err != nil {
		t.Error(err)
	}
}

func TestHashDistanceFormatFail(t *testing.T) {
	_, err := HashDistance("", h2)
	if err == nil {
		t.Errorf("%s", "Expected error, got nil")
	}
}

func TestHashDistanceFormatFail1(t *testing.T) {
	_, err := HashDistance(h5, h2)
	if err == nil {
		t.Errorf("%s", "Expected error, got nil")
	}
}

func TestHashDistanceFormatFail2(t *testing.T) {
	_, err := HashDistance(h2, h5)
	if err == nil {
		t.Errorf("%s", "Expected error, got nil")
	}
}

func TestHashDistanceFailLength(t *testing.T) {
	_, err := HashDistance(h1, h3)
	if err == nil {
		t.Errorf("%s", "Expected error, got nil")
	}
}

func TestHashDistanceLengthFail1(t *testing.T) {
	_, err := HashDistance("a1:hash:hash", h2)
	if err == nil {
		t.Errorf("%s", "Expected error, got nil")
	}
}

func TestHashDistanceLengthFail2(t *testing.T) {
	_, err := HashDistance(h2, "a1:hash:hash")
	if err == nil {
		t.Errorf("%s", "Expected error, got nil")
	}
}

func BenchmarkDistance(b *testing.B) {
	var h1 = `7DSC8olnoL1v/uawvbQD7XlZUFYzYyMb615NktYHF7dREN/JNnQrmhnUPI+/n2Y7`
	var h2 = `7DSC8olnoL1v/uawvbQD7XlZUFYzYyMb615NktYHF7dREN/JNnQrmhnUPI+/ngrr`
	for i := 0; i < b.N; i++ {
		distance(h1, h2)
	}
}
