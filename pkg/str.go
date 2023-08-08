package pkg

import (
	"bytes"
	"fmt"
	"strings"
)

func All(s string, ch rune) bool {
	for _, c := range s {
		if c != ch {
			return false
		}
	}
	return true
}

func First(s string, ch rune) bool {
	if len(s) == 0 {
		return false
	}
	r := bytes.Runes([]byte{s[0]})[0]
	return r == ch
}

func StartsWith(s string, start string) bool {
  if len(start) > len(s) {
    return false
  }
  cond := s[0:len(start)] == start
  if cond {
    fmt.Println(s)
  }
  return cond
}

func CountPrefixReturnRest(s string, ch rune) (int, string) {
	cnt := 0
	for i, c := range s {
		if c != ch {
			return cnt, s[i:]
		}
		c += 1
	}
	return cnt, ""
}

func HeadingName(s string) string {
  end := 0
  for i, c := range s {
    if c == '#' {
      continue
    }
    end = i
    break
  }

  return strings.TrimSpace(s[end:])
}
