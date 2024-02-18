package jess

import "git.sr.ht/~ionous/tapestry/support/match"

type localResults struct {
	Primary   []resultName
	Secondary []resultName // usually just one, except for nxm relations
	Macro     Macro
}

type resultName struct {
	Article articleResult
	Span    match.Span
	Exact   bool // when the phrase contains "called", we shouldn't fold the noun into other similarly named nouns.
	Traits  []Matched
	Kinds   []Matched // it's possible, if rare, to apply multiple kinds
	// ex. The container called the coffin is a closed openable thing.
}

func (n *resultName) String() string {
	return n.Span.String()
}

// FIX: just use article?
type articleResult struct {
	Matched
	Count int // for counted nouns: "seven (apples)"
}

func (a articleResult) NumWords() int {
	return MatchedLen(a.Matched)
}

func (a articleResult) String() (ret string) {
	if a.Matched != nil {
		ret = a.Matched.String()
	}
	return
}

// returns the size of a match;
// 0 if the match is nil.
func MatchedLen(m Matched) (ret int) {
	if m != nil {
		ret = m.NumWords()
	}
	return
}

// returns the string of a match;
// empty if the match is nil.
func MatchedString(m Matched) (ret string) {
	if m != nil {
		ret = m.String()
	}
	return
}
