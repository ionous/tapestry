package grok

import (
	"github.com/ionous/errutil"
)

// fix: should be customizable
var known = struct {
	determiners, macros, traits spans
}{
	determiners: makeSpans([]string{
		"the", "a", "an", "some", "our",
		// ex. kettle of fish
		"a kettle of",
	}),
	macros: makeSpans([]string{
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
	traits: makeSpans([]string{
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
		"transparent",
		"fixed in place",
	},
	),
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

func Grok(p string) (ret results, err error) {
	out := &results{}
	if words, e := hashWords(p); e != nil {
		err = e
	} else {
		// scan for "is/are" or a macro verb, which ever comes first:
		for i, w := range words {
			if w.equals(keywords.is) || w.equals(keywords.are) {
				err = afterIs(out, words[:i], words[i+1:])
				break
			} else {
				if at, skip := known.macros.findPrefix(words[i:]); skip > 0 {
					out.macro = known.macros[at]
					err = afterMacro(out, words[:i], words[i+skip:])
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

// after finding is/are, scan for a macro;
// handle either an "is-verb" or a pure "is" statement.
func afterIs(out *results, lhs, rhs []word) (err error) {
	var traits [][]word
	at, cnt := 0, len(rhs)
	for at < cnt {
		// break on extra is.
		if w := rhs[at]; w.equals(keywords.is) || w.equals(keywords.are) {
			err = makeWordError(w, "only one is/are expected")
			break
		}
		// scan for traits so multi word traits don't get confused for verbs
		// ( ex. "fixed in place" vs. "in the lobby" )
		// but these are dropped if a macro is actually detected.
		if trait, skip := known.traits.findPrefix(rhs[at:]); skip > 0 {
			if skipAnd, e := countAnd(rhs[skip:]); e != nil {
				err = e
				break
			} else {
				traits = append(traits, known.traits[trait])
				at += skip + skipAnd
			}
		} else if macro, skip := known.macros.findPrefix(rhs[at:]); skip == 0 {
			at++
		} else {
			// fix? parse the subject while first looping over the words?
			if e := genSubjects(out, lhs); e != nil {
				err = errutil.New("parsing subject", e)
			} else {
				// before the macro, and after
				// ex. the closed box "is" "in" the lobby
				// ex. the box "is" a closed container "in" the lobby
				pre, post := rhs[:at], rhs[at+skip:]

				// fix? this isnt satisfying
				if macro < 3 {
					if e := applyExtras(out, pre); e != nil {
						err = errutil.New("parsing extras", e)
					} else if e := applyKind(out, post); e != nil {
						err = errutil.New("parsing target", e)
					}
				} else {
					if rest, e := genNounCalled(&out.targets, post); e != nil {
						err = errutil.New("parsing target", e)
					} else if len(rest) > 0 {
						err = makeWordError(rest[0], "unexpected trailing text")
					} else if e := applyExtras(out, pre); e != nil {
						err = errutil.New("parsing extras", e)
					}
				}
				if err == nil {
					out.macro = known.macros[macro]
				}
				break // done now regardless
			}
		}
	}
	// never hit a macro?
	// ex. (the box) is (closed).
	if nothingLeft := at == cnt; nothingLeft {
		if e := genSubjects(out, lhs); e != nil {
			err = errutil.New("parsing subject", e)
		} else {
			applyTraits(out.subjects, traits)
		}
	}
	return
}

// XXXXXXX haven't actually looked at this.
func afterMacro(out *results, lhs, rhs []word) (err error) {
	panic("not implemented")
	// at, cnt := 0, len(ws)
	// for ; at < cnt; at++ {
	// 	if w := ws[at]; w.equals(keywords.is) || w.equals(keywords.are) {
	// 		err = makeWordError(w, "only one is/are expected")
	// 		break
	// 	} else {
	// 		err = foundMacroIs(out, ws[at+at:])
	// 		break
	// 	}
	// }
	// // never is?
	// if nothingLeft := at == len(ws); nothingLeft {
	// 	// //err = foundIsOnly(out, ws)
	// 	// // lhs are the nouns; rhs are the traits
	// 	// if e := genSubjects(out, lhs); e != nil {
	// 	// 	err = errutil.New("parsing subject", e)
	// 	// } else if e := applyExtras(out, rhs); e != nil {
	// 	// 	err = errutil.New("parsing extras", e)
	// 	// }
	// }
	return
}

// ex. the bottle) is (closed.
func foundIsOnly(out *results, lhs, rhs []word) (err error) {
	// lhs are the nouns; rhs are the traits
	if e := genSubjects(out, lhs); e != nil {
		err = errutil.New("parsing subject", e)
	} else if e := applyExtras(out, rhs); e != nil {
		err = errutil.New("parsing extras", e)
	}
	return
}

func foundMacroIs(out *results, ws []word) (err error) {
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
			if rem, e := genNounCalled(&out.subjects, words); e != nil {
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
		if traits, e := parseTraits(words); e != nil {
			err = e
		} else {
			applyTraits(out.subjects, traits)
		}
	}
	return
}

func applyTraits(out []noun, traits [][]word) {
	for i := range out {
		out[i].traits = append(out[i].traits, traits...)
	}
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
func genNounCalled(out *[]noun, words []word) (rest []word, err error) {
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

// like countAnd, but returns a new slice
func eatAnd(words []word) (ret []word, err error) {
	if skip, e := countAnd(words); e != nil {
		err = e
	} else {
		ret = words[skip:]
	}
	return
}

// note: when words were parsed commas became their own word
// either of these patterns will be ate ", and" or "and"
// returns the number of eaten words
func countAnd(ws []word) (ret int, err error) {
	var comma bool
	var and bool
	at, cnt := 0, len(ws)
Loop:
	for at < cnt {
		switch w := ws[at]; w.hash {
		default:
			break Loop

		case keywords.comma:
			if comma || and {
				err = makeWordError(w, "unexpected comma")
				break Loop
			}
			at, comma = at+1, true

		case keywords.and:
			if and {
				err = makeWordError(w, "unexpected and")
				break Loop
			}
			at, and = at+1, true
		}
	}
	// nothingness is okay, but not nothingness after a comma or and.
	if nothingLeft := at == len(ws); nothingLeft && (and || comma) {
		// to get nothing, we must have had no error...
		err = makeWordError(ws[0], "unexpected ending")
	} else {
		ret = at
	}
	return
}
