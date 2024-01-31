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
	res   []Result
}

func (in InputState) Results() []Result {
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
		res: append([]Result{}, in.res...),
	}
}

func (in InputState) Words() []grok.Word {
	return in.words[in.pos:]
}

func (in *InputState) AddResult(t ResultType, m grok.Match) {
	in.res = append(in.res, Result{t, m})
}

type Verb struct {
	VerbType string
}

type Count struct {
	One  []Interpreter
	Many []Interpreter
}

type Done struct {
}

func (*Verb) Match(q Query, cs []InputState) (ret []InputState) {
	panic("not implemented")
}

func (*Count) Match(q Query, cs []InputState) (ret []InputState) {
	panic("not implemented")
}

// matches a specified string and advances the input
type Words struct {
	Str   string
	cache grok.Span
	err   error
}

// fix: cache?
func (op *Words) Match(q Query, cs []InputState) (ret []InputState) {
	if match, e := op.makeSpan(); e == nil {
		for _, in := range cs {
			if grok.HasPrefix(in.Words(), match) {
				rest := in.Next(len(match))
				ret = append(ret, rest)
			}
		}
	}
	return
}

func (op *Words) makeSpan() (grok.Span, error) {
	if op.cache == nil && op.err == nil {
		op.cache, op.err = grok.MakeSpan(op.Str)
	}
	return op.cache, op.err
}

// tries each branch; the result is the union of all separate successes.
type Branch struct {
	Options []Interpreter
}

func (op *Branch) Match(q Query, cs []InputState) (ret []InputState) {
	for _, m := range op.Options {
		if next := m.Match(q, cs); len(next) > 0 {
			ret = append(ret, next...)
		}
	}
	return
}

// matches each of its sub-expressions in turn
// filter the results with each new
type Sequence struct {
	Series []Interpreter
}

func (op *Sequence) Match(q Query, cs []InputState) []InputState {
	next := cs
	for _, m := range op.Series {
		if next = m.Match(q, next); len(next) == 0 {
			break
		}
	}
	return next
}

type Is struct{}

func (*Is) Match(q Query, cs []InputState) (ret []InputState) {
	for _, in := range cs {
		if ws := in.Words(); len(ws) > 0 {
			if w := ws[0]; w.Equals(grok.Keyword.Is) || w.Equals(grok.Keyword.Are) {
				ret = append(ret, in.Next(1))
			}
		}
	}
	return
}

func (*Done) Match([]InputState) (ret []InputState) {
	panic("not implemented")
}

// type Grokker interface {
// 	// if the passed words starts with a determiner,
// 	// return the number of words that matched.
// 	FindArticle(Span) (Article, error)

// 	// if the passed words starts with a kind,
// 	// return the number of words that matched.
// 	FindKind(Span) (Match, error)

// 	// if the passed words starts with a trait,
// 	// return the number of words that matched.
// 	FindTrait(Span) (Match, error)

// 	// if the passed words starts with a macro,
// 	// return information about that match.
// 	FindMacro(Span) (Macro, error)
// }
