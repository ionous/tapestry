package jess

import (
	"git.sr.ht/~ionous/tapestry/support/match"
)

// helper for line matching
type InputState struct {
	p        *Paragraph // helpful for debugging
	pronouns pronounSource
	// fix? change to a "reader" ( pull ) rather than pre-process?
	words []match.TokenValue // the current line from the paragraph
}

func (in InputState) Len() int {
	return len(in.words)
}

func (in InputState) Words() []match.TokenValue {
	return in.words
}

func (in InputState) GetNext(t match.Token) (ret match.TokenValue, okay bool) {
	if len(in.words) > 0 && in.words[0].Token == t {
		ret, okay = in.words[0], true
	}
	return
}

// return an input state that is the passed number of words after this one.
func (in InputState) Skip(skip int) InputState {
	return InputState{
		p:        in.p,
		pronouns: in.pronouns,
		words:    in.words[skip:],
	}
}

// return the specified number of words from the input as a slice of words
func (in InputState) Cut(width int) Matched {
	return in.words[:width]
}

// see if the input begins with any of the passed choices
// always returns 1.
func (in InputState) MatchWord(choices ...uint64) (width int) {
	if len(in.words) > 0 {
		w := in.words[0]
		for _, opt := range choices {
			if w.Equals(opt) {
				width = 1
				break
			}
		}
	}
	return
}
