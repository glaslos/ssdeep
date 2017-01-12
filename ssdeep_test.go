package ssdeep

import "testing"

func TestTimeConsuming(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
}

func TestFNVHash(t *testing.T) {
	if hash("ssdeep") != 3010289213 {
		t.Error("FNV hash not matching")
	}
}

func TestRollingHash(t *testing.T) {
	rs := newRollingState()
	if rollHash(&rs, byte('A')) != 585 {
		t.Error("Rolling hash not matching")
	}
	rollHash(&rs, byte('B'))
	rollHash(&rs, byte('C'))
}

func TestFindBlocks(t *testing.T) {
	getBlocks()
}

func BenchmarkRollingHash(b *testing.B) {
	rs := newRollingState()
	for i := 0; i < b.N; i++ {
		rollHash(&rs, byte(i))
	}
}
