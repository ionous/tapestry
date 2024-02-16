package grok

import "strings"

type Grokker interface {
	// if the passed words starts with a determiner,
	// return the number of words that matched.
	FindArticle(Span) (Article, error)

	// if the passed words starts with a kind,
	// return the number of words that matched.
	FindKind(Span) (Matched, error)

	// if the passed words starts with a trait,
	// return the number of words that matched.
	FindTrait(Span) (Matched, error)

	// if the passed words starts with a macro,
	// return information about that match.
	FindMacro(Span) (Macro, error)
}

type Article struct {
	Matched
	Count int // for counted nouns: "seven (apples)"
}

func (a Article) Len() int {
	return MatchedLen(a.Matched)
}

func (a Article) String() (ret string) {
	if a.Matched != nil {
		ret = a.Matched.String()
	}
	return
}

type Macro struct {
	Matched
	Name     string
	Type     MacroType
	Reversed bool
}

func (m Macro) Len() int {
	return MatchedLen(m.Matched)
}

type Results struct {
	Primary   []Name
	Secondary []Name // usually just one, except for nxm relations
	Macro     Macro
}

type Name struct {
	Article Article
	Span    Span
	Exact   bool // when the phrase contains "called", we shouldn't fold the noun into other similarly named nouns.
	Traits  []Matched
	Kinds   []Matched // it's possible, if rare, to apply multiple kinds
	// ex. The container called the coffin is a closed openable thing.
}

func (n *Name) String() string {
	return n.Span.String()
}

// Matched - generic interface so Grokker implementations can track their own backchannel data.
// ex. a row id for kinds.
type Matched interface {
	String() string
	NumWords() int // should match strings.Fields(Str)
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

// Span - implements Match for a chain of individual words.
type Span []Word

func (s Span) String() string {
	var b strings.Builder
	for i, w := range s {
		if i > 0 {
			b.WriteRune(' ')
		}
		b.WriteString(w.String())
	}
	return b.String()
}

func (s Span) NumWords() int {
	return len(s)
}

func HasPrefix(s, prefix []Word) (okay bool) {
	// a prefix must be the same as or shorter than us
	if len(prefix) <= len(s) {
		okay = true // provisionally
		for i, a := range prefix {
			if a.Hash() != s[i].Hash() {
				okay = false
				break
			}
		}
	}
	return
}

// make customizable?
var Keyword = struct {
	And, Are, Called, Comma, Is uint64
}{
	And:    Hash("and"),
	Are:    Hash("are"),
	Called: Hash("called"),
	Comma:  Hash(","),
	Is:     Hash("is"),
}
