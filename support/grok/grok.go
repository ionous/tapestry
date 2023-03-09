package grok

// fix: should be customizable
var known = struct {
	determiners, macros, kinds, traits spans
}{
	determiners: panicSpans([]string{
		"the", "a", "an", "some", "our",
		// ex. kettle of fish
		"a kettle of",
	}),
	macros: panicSpans([]string{
		// right now assumes the first set are kinds of
		// could it use the fact there's only one set of definitions.
		// is any thing else like that?
		// tbd: need more thought.
		"kind of",   // for "a closed kind of container"
		"kinds of",  // for "are closed containers"
		"a kind of", // for "a kind of container"
		// other macros
		"on",
		"in",
	}),
	kinds: panicSpans([]string{
		"thing", "things",
		"container", "containers",
		"supporter", "supporters",
	}),
	traits: panicSpans([]string{
		"closed",
		"open",
		"openable",
		"transparent",
		"fixed in place",
	}),
}

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

// fix: would be better to have a push interface so we can just add things as we go
// this is easier for development though
type Results struct {
	sources []Noun
	targets []Noun // usually just one, except for nxm relations
	macro   []Word
}

type Noun struct {
	det    []Word
	name   []Word
	traits [][]Word
	kinds  [][]Word // it's possible, if rare, to apply multiple kinds
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
				if at, skip := known.macros.findPrefix(words[i:]); skip > 0 {
					out.macro = known.macros[at]
					err = macroPhrase(out, words[i+skip:])
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
