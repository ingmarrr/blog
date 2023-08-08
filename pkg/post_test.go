package pkg

import (
	"bytes"
	"testing"
)

func TestExtractToc(t *testing.T) {
	tst := "# Hello"
	rn := '#'
	tstR := bytes.Runes([]byte{tst[0]})[0]
	if rn != tstR {
		t.Error("Should be the same thing")
	}
}
