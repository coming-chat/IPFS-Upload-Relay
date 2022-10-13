package utils

import (
	"os"
	"testing"
)

func TestCalcFileHash(t *testing.T) {
	file := "../tmp/10GB.bin"
	f, _ := os.Open(file)
	defer f.Close()
	hash := CalcFileHash(f)
	t.Log(hash)
}
