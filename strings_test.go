package xlib

import (
	"strings"
	"testing"
)

func TestStrJoin(t *testing.T) {
	cases := [...]struct {
		res  string
		args []string
	}{
		{"", nil},
		{"AAA", []string{"AAA"}},
		{"AAA, BBB", []string{"AAA", "BBB"}},
		{"AAA, BBB, CCC", []string{"AAA", "BBB", "CCC"}},
		{"AAA, BBB, CCC, DDD", []string{"AAA", "BBB", "CCC", "DDD"}},
		{"AAA, BBB, CCC, DDD, EEE", []string{"AAA", "BBB", "CCC", "DDD", "EEE"}},
	}

	for i, c := range cases {
		if r := JoinStrings(", ", c.args...); r != c.res {
			t.Errorf("[%d] unexpected result: %q instead of %q", i, r, c.res)
			return
		}
	}
}

func TestWrites(t *testing.T) {
	var buff strings.Builder

	if err := WriteString(&buff, "Hello!"); err != nil {
		t.Error(err)
		return
	}

	if err := WriteByte(&buff, ' '); err != nil {
		t.Error(err)
		return
	}

	if err := WriteRune(&buff, '😀'); err != nil {
		t.Error(err)
		return
	}

	if s := buff.String(); s != "Hello! 😀" {
		t.Errorf("unexpected string: %q instead of \"Hello! 😀\"", s)
		return
	}
}
