package jess

import (
	"git.sr.ht/~ionous/tapestry/support/match"
)

// the input state consists of a series of hashes
// and chunks of the original string.
// fix? change to a "reader" ( pull ) rather than pre-process?
type InputState []match.TokenValue

func (in InputState) Len() int {
	return len(in)
}

func (in InputState) Words() []match.TokenValue {
	return in
}

func (in InputState) GetNext(t match.Token) (ret any, okay bool) {
	if len(in) > 0 && in[0].Token == t {
		ret, okay = in[0].Value, true
	}
	return
}

// return an input state that is the passed number of words after this one.
func (in InputState) Skip(skip int) InputState {
	return in[skip:]
}

// return the specified number of words from the input as a slice of words
func (in InputState) Cut(width int) []match.TokenValue {
	return in[:width]
}

// see if the input begins with any of the passed choices
// always returns 1.
func (in InputState) MatchWord(choices ...uint64) (width int) {
	if len(in) > 0 {
		w := in[0]
		for _, opt := range choices {
			if w.Equals(opt) {
				width = 1
				break
			}
		}
	}
	return
}
