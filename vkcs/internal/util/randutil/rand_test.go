package randutil

import "testing"

func TestRandomName(t *testing.T) {
	rndName := RandomName(5)
	if len(rndName) != 5 {
		t.Fatalf("Got wrong result length: %d, expected: 5", len(rndName))
	}
}
