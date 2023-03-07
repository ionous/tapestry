package grok

import (
	"github.com/ionous/errutil"
)

// fix: should be customizable
var traits = []string{
	// fix: kinds as traits is handy, but becomes difficult with plurals.
	// and there _are_ some places where its wrong to have a trait and required to have a kind
	"thing",
	"container",
	"supporter",
	"things",
	"containers",
	"supporters",
	// real traits
	"closed",
	"open",
	"openable",
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

// fix: should be customizable
var known = struct {
	determiners, traits spans
}{
	traits: makeSpans(traits),
	determiners: makeSpans([]string{
		"the", "a", "an", "some", "our",
		// ex. kettle of fish
		"a kettle of",
	}),
}

type match struct {
	word  word // original word and reference to the text
	index int  // position within
}

// fix: would be better to have a push interface so we can just add things as we go
// this is easier for development though
type results struct {
	subjects []noun
	targets  []noun // usually just one, except for nxm relations
	macro    span
}

// fix: should be customizable;
// fix: special handling for kinds/of?
// should things have aliases?
var macros = makeSpans([]string{
	"kind of",   // for "a closed kind of container"
	"kinds of",  // for "are closed containers"
	"a kind of", // for "a kind of container"
	"on",
})

func Grok(p string) (ret results, err error) {
	out := &results{}
	if words, e := hashWords(p); e != nil {
		err = e
	} else {
		matched := struct {
			macro, toBe int // one-index within "words" of the successful match
			whichMacro  int // zero index with the list of all macros
		}{}

		// scan left to right, looking for "is/are" or a macro phrase
		for i, w := range words {
			if w.equals(keywords.is) || w.equals(keywords.are) {
				// split here.
				if matched.toBe > 0 {
					err = makeWordError(w, "only one is/are expected")
					break
				}
				matched.toBe = i + 1
			}
			// should this try to keep matching no matter what? ( like to/be does )
			if matched.macro == 0 {
				restOfTheSentence := words[i:]
				if at := macros.findPrefix(restOfTheSentence); at >= 0 {
					out.macro = macros[at]
					matched.macro = i + 1
					matched.whichMacro = at
				}
			}
			if matched.toBe > 0 && matched.macro > 0 {
				break
			}
		}
		// explicit error if "toBe" falls at the end of a sentence?
		if err == nil {
			switch {
			case matched.toBe > 0 && matched.macro > 0:
				// ex. the closed box "is" "in" the lobby
				// ex. the box "is" a closed container "in" the lobby
				macroLen := len(macros[matched.whichMacro])
				if beStart, macStart := matched.toBe-1, matched.macro-1; beStart < macStart {
					lhs, rhs, mid := words[:beStart], words[macStart+macroLen:], words[beStart+1:macStart]
					// fix? parse the source while first looping over the words.
					if e := genSubjects(out, lhs); e != nil {
						err = errutil.New("parsing source", e)
					} else {
						// fix? this isnt satisfying
						if matched.whichMacro < 3 {
							if e := applyExtras(out, mid); e != nil {
								err = errutil.New("parsing extras", e)
							} else if e := applyKind(out, rhs); e != nil {
								err = errutil.New("parsing target", e)
							}
						} else {
							if rest, e := genNoun(&out.targets, rhs); e != nil {
								err = errutil.New("parsing target", e)
							} else if len(rest) > 0 {
								err = makeWordError(rest[0], "unexpected trailing text")
							} else if e := applyExtras(out, mid); e != nil {
								err = errutil.New("parsing extras", e)
							}
						}
					}
				}
			}
		}
	}
	if err == nil {
		ret = *out
	}
	return
}

type noun struct {
	det  []word
	name []word
	//kind   []word
	traits [][]word
}

// starting simple:
// <subjects> <are> <extras> <target>
// the <traits> <name> [and] ...
// ex. the closed box and the bottle ...
// or. the closed container called the box and the bottle....
func genSubjects(out *results, words []word) (err error) {
	if len(words) == 0 {
		err = errutil.New("expected at least one subject")
	} else {
		for len(words) > 0 && err == nil {
			if rem, e := genNoun(&out.subjects, words); e != nil {
				err = e
			} else {
				words, err = eatAnd(rem)
			}
		}
	}
	return
}

// extras apply equally to all subjects:
// <subjects> <are> <extras> <target>
// ex. the box is [a] closed [container] in the lobby.
func applyExtras(out *results, words []word) (err error) {
	if len(words) > 0 {
		if ts, e := parseTraits(words); e != nil {
			err = e
		} else {
			for i, n := range out.subjects {
				out.subjects[i].traits = append(n.traits, ts...)
			}
		}
	}
	return
}

// a kind of <traits> <kind>.
// future: and has; with value?
func applyKind(out *results, words []word) (err error) {
	if len(words) == 0 {
		err = errutil.New("expected some kind words")
	} else {
		err = applyExtras(out, words)
	}
	return
}

// <subjects> <are> <target>
// ex. ... <the traits> called <the noun> [and ....]
// ex. ... <the noun> [and ....]
// note: without the "called" leading words become the name of the noun.
// returns any unprocessed words
func genNoun(out *[]noun, words []word) (rest []word, err error) {
	var traits [][]word
	var andAt int // one indexed
	for i, t := range words {
		// track possible traits
		if t.equals(keywords.and) || t.equals(keywords.comma) { // fix something to error if nothing preceeds and
			if i == 0 {
				err = makeWordError(t, "missing words")
				break
			}
			andAt = i + 1
		} else if t.equals(keywords.called) {
			lhs, rhs := words[:i], words[i+1:]
			if ts, e := parseTraits(lhs); e != nil {
				err = e
				break
			} else if len(rhs) == 0 {
				err = makeWordError(t, "missing a noun")
				break
			} else {
				traits = ts
				words = rhs // shorten the parsing to skip the "called"
				andAt = 0   // forget any "and(s)", those would have been for traits
				break
			}
		}
	}
	if err == nil {
		// the and was the end of the noun, we just didn't know it till we hit the end of the loop
		if andAt > 0 {
			words, rest = words[:andAt-1], words[andAt:]
		}
		if n, e := parseNoun(words, traits); e != nil {
			err = e
		} else {
			*out = append(*out, n)
		}
	}
	return
}

func parseNoun(words []word, traits [][]word) (ret noun, err error) {
	if det, rest := known.determiners.cut(words); len(rest) == 0 {
		err = makeWordError(words[0], "missing a noun")
	} else {
		ret = noun{
			det:    det,
			name:   rest,
			traits: traits,
		}
	}
	return
}

// ex. "the open container"
// ex. "the open and openable container"
// tbd: should this more specifically test for kinds in the last slot?
// ( some situations, like 'extras' allow any trait, some like 'targets' require a kind )
func parseTraits(words []word) (traits [][]word, err error) {
	// for kinds, we dont really care about the determiner
	if _, rest := known.determiners.cut(words); len(rest) == 0 {
		err = errutil.New("expected some sort of name")
	} else {
		for len(rest) > 0 && err == nil {
			if trait, rem := known.traits.cut(rest); len(trait) == 0 {
				err = makeWordError(rest[0], "unknown trait")
			} else {
				traits = append(traits, trait)
				rest, err = eatAnd(rem)
			}
		}
	}
	return
}

// note: when words were parsed commas became their own word
// either of these patterns will be ate ", and" or "and"
func eatAnd(words []word) (ret []word, err error) {
	var comma bool
	var and bool
	ate := words
Loop:
	for len(ate) > 0 {
		switch w := ate[0]; w.hash {
		default:
			break Loop

		case keywords.comma:
			if comma || and {
				err = makeWordError(w, "unexpected comma")
				break Loop
			}
			ate, comma = ate[1:], true

		case keywords.and:
			if and {
				err = makeWordError(w, "unexpected and")
				break Loop
			}
			ate, and = ate[1:], true
		}
	}
	// nothingness is okay, but not nothingness after a comma or and.
	if err == nil && (and || comma) && len(ate) == 0 {
		err = makeWordError(words[0], "unexpected ending")
	} else {
		ret = ate
	}
	return
}
