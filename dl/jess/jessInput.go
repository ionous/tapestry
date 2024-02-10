package jess

import "git.sr.ht/~ionous/tapestry/support/grok"

type InputState struct {
	words   []grok.Word
	pos     int
	Matched []grok.Matched
}

func MakeInputState(words []grok.Word) InputState {
	return InputState{words: words}
}

// return an input state that is the passed number of words after this one.
func (in InputState) Skip(skip int) InputState {
	return InputState{
		words: in.words,
		pos:   in.pos + skip,
		// clone the array so that results dont accidentally share memory
		Matched: append([]grok.Matched{}, in.Matched...),
	}
}

// return an input state that is the passed number of words after this one.
func (in InputState) Cut(width int) grok.Span {
	sub := in.words[in.pos : in.pos+width]
	return sub
}

func (in InputState) Words() []grok.Word {
	return in.words[in.pos:]
}

func (in InputState) MatchWord(choices ...uint64) (width int) {
	if ws := in.Words(); len(ws) > 0 {
		w := ws[0]
		for _, opt := range choices {
			if w.Equals(opt) {
				width = 1
				break
			}
		}
	}
	return
}
