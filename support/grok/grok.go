package grok

import "strings"

type Grokker interface {
	// if the passed words starts with a determiner,
	// return the number of words that matched.
	FindArticle(Span) (Article, error)

	// if the passed words starts with a kind,
	// return the number of words that matched.
	FindKind(Span) (Match, error)

	// if the passed words starts with a trait,
	// return the number of words that matched.
	FindTrait(Span) (Match, error)

	// if the passed words starts with a macro,
	// return information about that match.
	FindMacro(Span) (Macro, error)
}

type Article struct {
	Match
	Count int // for counted nouns: "seven (apples)"
}

func (a Article) Len() int {
	return MatchedLen(a.Match)
}

func (a Article) String() (ret string) {
	if a.Match != nil {
		ret = a.Match.String()
	}
	return
}

type Macro struct {
	Match
	Name     string
	Type     MacroType
	Reversed bool
}

func (m Macro) Len() int {
	return MatchedLen(m.Match)
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
	Traits  []Match
	Kinds   []Match // it's possible, if rare, to apply multiple kinds
	// ex. The container called the coffin is a closed openable thing.
}

func (n *Name) String() string {
	return n.Span.String()
}

// Match - generic interface so Grokker implementations can track their own backchannel data.
// ex. a row id for kinds.
type Match interface {
	String() string
	NumWords() int // should match strings.Fields(Str)
}

// returns the size of a match;
// 0 if the match is nil.
func MatchedLen(m Match) (ret int) {
	if m != nil {
		ret = m.NumWords()
	}
	return
}

// returns the string of a match;
// empty if the match is nil.
func MatchedString(m Match) (ret string) {
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

// expects at most a single sentence.
func Grok(known Grokker, p string) (ret Results, err error) {
	if words, e := MakeSpan(p); e != nil {
		err = e
	} else {
		ret, err = GrokSpan(known, words)
	}
	return
}

func GrokSpan(known Grokker, words Span) (ret Results, err error) {
	// scan for "is/are" or a macro verb, which ever comes first;
	// the order can reverse subjects and objects.
	for i, w := range words {
		if w.equals(Keyword.Is) || w.equals(Keyword.Are) {
			ret, err = beingPhrase(known, words[:i], words[i+1:])
			break
		} else if macro, e := known.FindMacro(words[i:]); e != nil {
			err = e
		} else if len(macro.Name) > 0 {
			ret, err = macroPhrase(known, macro, words[i+macro.Match.NumWords():])
			break
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
