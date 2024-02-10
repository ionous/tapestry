package jess

import "git.sr.ht/~ionous/tapestry/support/grok"

type InputState []grok.Word

func (in InputState) Words() []grok.Word {
	return in
}

// return an input state that is the passed number of words after this one.
func (in InputState) Skip(skip int) InputState {
	return in[skip:]
}

// return an input state that is the passed number of words after this one.
func (in InputState) Cut(width int) grok.Span {
	return grok.Span(in[:width])
}

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
