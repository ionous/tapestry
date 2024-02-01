package gestalt

import "git.sr.ht/~ionous/tapestry/support/grok"

type Interpreter interface {
	// Scan for results.
	// note: by design, cursor may be out of range when scan is called.
	Match(Query, []InputState) []InputState
}

type InputState struct {
	words []grok.Word
	pos   int
	res   []grok.Match
}

func (in InputState) Results() []grok.Match {
	return in.res
}

func MakeInputState(words []grok.Word) InputState {
	return InputState{words: words}
}

// return an input state that is the passed number of words after this one.
func (in InputState) Next(skip int) InputState {
	return InputState{
		words: in.words,
		pos:   in.pos + skip,
		// clone the array so that results dont accidentally share memory
		res: append([]grok.Match{}, in.res...),
	}
}

func (in InputState) Words() []grok.Word {
	return in.words[in.pos:]
}

func (in *InputState) AddResult(m grok.Match) {
	in.res = append(in.res, m)
}
