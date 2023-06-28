package grok

type Grokker interface {
	// if the passed words starts with a determiner,
	// return the number of words in  that match.
	FindDeterminer([]Word) int

	// if the passed words starts with a kind,
	// return the number of words in  that match.
	FindKind([]Word) int

	// if the passed words starts with a trait,
	// return the number of words in  that match.
	FindTrait([]Word) int

	// if the passed words starts with a macro,
	// return information about that match
	FindMacro([]Word) (ret MacroInfo, okay bool)
}

type MacroInfo struct {
	Type  MacroType
	Str   string
	Width int // should match strings.Fields(Str)
}

type Results struct {
	Sources []Noun
	Targets []Noun // usually just one, except for nxm relations
	Macro   MacroInfo
}

type Noun struct {
	Det    []Word
	Name   []Word
	Traits [][]Word
	Kinds  [][]Word // it's possible, if rare, to apply multiple kinds
	// ex. The container called the coffin is a closed openable thing.
}

func Grok(p string) (ret Results, err error) {
	out := &Results{}
	if words, e := hashWords(p); e != nil {
		err = e
	} else {
		// scan for "is/are" or a macro verb, which ever comes first;
		// the order can reverse subjects and objects.
		for i, w := range words {
			if w.equals(keywords.is) || w.equals(keywords.are) {
				err = beingPhrase(out, words[:i], words[i+1:])
				break
			} else {
				if macro, ok := known.FindMacro(words[i:]); ok {
					out.Macro = macro
					err = macroPhrase(out, words[i+macro.Width:])
					break
				}
			}
		}
	}
	if err == nil {
		ret = *out
	}
	return
}

// make customizable?
var keywords = struct {
	and, are, called, comma, has, is uint64
}{
	and:    plainHash("and"),
	are:    plainHash("are"),
	called: plainHash("called"),
	comma:  plainHash(","),
	has:    plainHash("has"),
	is:     plainHash("is"),
}
