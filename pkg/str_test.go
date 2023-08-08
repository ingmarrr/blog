package pkg

import (
  "testing"
)


func TestStartsWith(t *testing.T) {
  s := "@Dend hello"
  start := "@Dend"

  if !StartsWith(s, start) {
    t.Error("Should be true")
  }
  
}

