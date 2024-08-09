package jess

import (
	"git.sr.ht/~ionous/tapestry/support/match"
)

// The english language phrase being matched.
type InputState struct {
	// p           *Paragraph
	// phraseIndex int
	words []match.TokenValue // the remainder of the phrase being parsed
}

// func (in InputState) Source() (ret compact.Source) {
// 	if p, ws := in.p, in.words; len(ws) > 0 {
// 		lineOfs := ws[0].Pos.Y // this is the source file line ( not phrase )
// 		ret = compact.Source{
// 			File:    p.File,
// 			Line:    lineOfs,
// 			Comment: "a plain-text paragraph",
// 		}
// 	}
// 	return
// }

func (in InputState) Len() int {
	return len(in.words)
}

func (in InputState) DebugString() string {
	return match.DebugStringify(in.words)
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
		// p:           in.p,
		// phraseIndex: in.phraseIndex,
		words: in.words[skip:],
	}
}

// return an input state of the passed width; dropping the trailing ones.
func (in InputState) Slice(width int) InputState {
	return InputState{
		// p:           in.p,
		// phraseIndex: in.phraseIndex,
		words: in.words[:width],
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
