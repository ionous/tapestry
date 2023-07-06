package grok

type Grokker interface {
	// if the passed words starts with a determiner,
	// return the number of words in  that match.
	FindArticle([]Word) Match

	// if the passed words starts with a kind,
	// return the number of words in  that match.
	FindKind([]Word) Match

	// if the passed words starts with a trait,
	// return the number of words in  that match.
	FindTrait([]Word) Match

	// if the passed words starts with a macro,
	// return information about that match
	FindMacro([]Word) (MacroInfo, bool)
}

type MacroInfo struct {
	Name     string
	Match    Match
	Type     MacroType
	Reversed bool
}

type Results struct {
	Sources []Noun
	Targets []Noun // usually just one, except for nxm relations
	Macro   MacroInfo
}

type Noun struct {
	Det    Match
	Name   Span
	Traits []Match
	Kinds  []Match // it's possible, if rare, to apply multiple kinds
	// ex. The container called the coffin is a closed openable thing.
}

// Match - generic interface so Grokker implementations can track their own backchannel data.
// ex. a row id for kinds.
type Match interface {
	String() string
	NumWords() int // should match strings.Fields(Str)
}

// Span - implements Match for a chain of individual words.
type Span []Word

func (s Span) String() string {
	return WordsWithSep(s, ' ')
}

func (s Span) NumWords() int {
	return len(s)
}

func Grok(known Grokker, p string) (ret Results, err error) {
	if words, e := MakeSpan(p); e != nil {
		err = e
	} else {
		// scan for "is/are" or a macro verb, which ever comes first;
		// the order can reverse subjects and objects.
		for i, w := range words {
			if w.equals(keywords.is) || w.equals(keywords.are) {
				ret, err = beingPhrase(known, words[:i], words[i+1:])
				break
			} else {
				if macro, ok := known.FindMacro(words[i:]); ok {
					ret, err = macroPhrase(known, macro, words[i+macro.Match.NumWords():])
					break
				}
			}
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
