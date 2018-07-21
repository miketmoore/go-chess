package chess_test

import (
	"testing"

	"github.com/miketmoore/chess"
)

func TestGetNextFile(t *testing.T) {
	tests := []struct {
		in   string
		out  string
		pass bool
	}{
		{"invalidinput", "", false},
		{"a", "b", true},
		{"b", "c", true},
		{"c", "d", true},
		{"d", "e", true},
		{"e", "f", true},
		{"f", "g", true},
		{"g", "h", true},
		{"h", "", false},
	}
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			next, ok := chess.GetNextFile(test.in)
			assertOk(t, ok == test.pass)
			assertOk(t, next == test.out)
		})

	}

}

func TestGetPreviousFile(t *testing.T) {
	tests := []struct {
		in   string
		out  string
		pass bool
	}{
		{"", "", false},
		{"h", "g", true},
		{"g", "f", true},
		{"f", "e", true},
		{"e", "d", true},
		{"d", "c", true},
		{"c", "b", true},
		{"b", "a", true},
		{"a", "", false},
	}
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			next, ok := chess.GetPreviousFile(test.in)
			assertOk(t, ok == test.pass)
			assertOk(t, next == test.out)
		})

	}

}

func assertOk(t *testing.T, b bool) {
	if b == false {
		t.Fatal("not ok")
	}
}
