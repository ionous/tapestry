package jess

import "git.sr.ht/~ionous/tapestry/support/match"

// the input state consists of a series of hashes
// and chunks of the original string.
//
// fix? rather than pre-processing the string
// a "reader" might make more sense.
// especially as seeing some pieces of the string
// ( ex. quoted text ) need to be processed again later on.
// first pass would be change makeSpan into a reader ( MakeSpans uses that reader )
// and then expose that reader here.
type InputState struct {
	ws    []match.Word
	index int
}

func MakeInput(ws []match.Word) InputState {
	return InputState{ws, 0}
}

func (in InputState) Len() int {
	return len(in.ws)
}

func (in InputState) Offset() int {
	return in.index
}

func (in InputState) Words() []match.Word {
	return in.ws
}

// return an input state that is the passed number of words after this one.
func (in InputState) Skip(skip int) InputState {
	return InputState{in.ws[skip:], in.index + skip}
}

// return an input state that is the passed number of words after this one.
func (in InputState) Cut(width int) string {
	return match.Span(in.ws[:width]).String()
}

func (in InputState) MatchWord(choices ...uint64) (width int) {
	if len(in.ws) > 0 {
		w := in.ws[0]
		for _, opt := range choices {
			if w.Equals(opt) {
				width = 1
				break
			}
		}
	}
	return
}
