package grok

import "strings"

type Grokker interface {
	// if the passed words starts with a determiner,
	// return the number of words in  that match.
	FindArticle(Span) (Article, error)

	// if the passed words starts with a kind,
	// return the number of words in  that match.
	FindKind(Span) (Match, error)

	// if the passed words starts with a trait,
	// return the number of words in  that match.
	FindTrait(Span) (Match, error)

	// if the passed words starts with a macro,
	// return information about that match
	FindMacro(Span) (Macro, error)
}

type Article struct {
	Match Match
	Count int
}

func (a Article) String() (ret string) {
	if a.Match != nil {
		ret = a.Match.String()
	}
	return
}

type Macro struct {
	Name     string
	Match    Match
	Type     MacroType
	Reversed bool
}

type Results struct {
	Sources []Noun
	Targets []Noun // usually just one, except for nxm relations
	Macro   Macro
}

type Noun struct {
	Article Article
	Name    Span
	Exact   bool // when the phrase contains "called", we shouldn't fold the noun into other similarly named nouns.
	Traits  []Match
	Kinds   []Match // it's possible, if rare, to apply multiple kinds
	// ex. The container called the coffin is a closed openable thing.
}

// Match - generic interface so Grokker implementations can track their own backchannel data.
// ex. a row id for kinds.
type Match interface {
	String() string
	NumWords() int // should match strings.Fields(Str)
}

// returns the size of a match;
// 0 if the match is nil.
func MatchLen(m Match) (ret int) {
	if m != nil {
		ret = m.NumWords()
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
		if w.equals(keywords.is) || w.equals(keywords.are) {
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
var keywords = struct {
	and, are, called, comma, has, is uint64
}{
	and:    Hash("and"),
	are:    Hash("are"),
	called: Hash("called"),
	comma:  Hash(","),
	has:    Hash("has"),
	is:     Hash("is"),
}
