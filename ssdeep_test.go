package ssdeep

import (
	"fmt"
	"log"
	"os/exec"
	"testing"
)

func TestTimeConsuming(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
}

func TestRollingHash(t *testing.T) {
	sdeep := NewSSDEEP()
	if rollHash(&sdeep, byte('A')) != 585 {
		t.Error("Rolling hash not matching")
	}
}

func TestFindBlock(t *testing.T) {
	Fuzzy("/tmp/dat")
	out, err := exec.Command("ssdeep", "/tmp/dat").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}

func BenchmarkRollingHash(b *testing.B) {
	sdeep := NewSSDEEP()
	for i := 0; i < b.N; i++ {
		rollHash(&sdeep, byte(i))
	}
}
