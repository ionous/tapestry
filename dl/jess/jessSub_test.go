package jess_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/support/match"
)

func TestSubAssignments(t *testing.T) {
	test := func(str string) {
		var m jess.SubAssignment
		in := jess.MakeInput(match.PanicSpan(str))
		if !m.Match(&in) {
			t.Fatal("couldnt match", str)
		} else if _, e := m.Deduce(); e != nil {
			t.Fatal("couldnt deduce assignment", e)
		}
	}
	test(assignText)
	test(assignExe)
}

var assignText = `:
  FromText: "text"
# plain text`

var assignExe = `:
  - Say: "hello"
# plain text`
