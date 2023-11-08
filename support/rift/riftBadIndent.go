package rift

import "git.sr.ht/~ionous/tapestry/support/charm"

type badIndent struct {
	have, want int // number of spaces
}

func (badIndent) Error() string {
	return "bad indent"
}

func BadIndent(have, want int) charm.State {
	return charm.Error(badIndent{have, want})
}
