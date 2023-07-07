package story

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

// readNounsWithProperties -
// reads ancillary information about nouns from their names and declares properties for them.
// ex. proper or plural names, etc.
func readNounsWithProperties(w *weave.Weaver, nouns []string) (ret []string, err error) {
	ret = make([]string, 0, len(nouns))
	for _, noun := range nouns {
		if next, e := importNoun(w, noun, ret); e != nil {
			err = e
			break
		} else {
			ret = next
		}
	}
	return
}

// readNouns - reads noun names without declaring any properties....
// fix? unless they are counted nouns ( for backwards compatibility of "two cats whereabouts the kitchen" )
func readNouns(w *weave.Weaver, nouns []string) (ret []string, err error) {
	ret = make([]string, 0, len(nouns))
	for _, noun := range nouns {
		if a, e := makeArticleName(w, noun); e != nil {
			err = e
			break
		} else if a.count == 0 {
			ret = append(ret, a.name)
		} else if ns, e := importCountedNoun(w.Catalog, a.count, a.name); e == nil {
			ret = append(ret, ns...)
		} else {
			err = errutil.New("couldn't import counted nouns", a.count, a.name, e)
			break
		}
	}
	return
}

// helper to simplify setting the values of nouns
func assertNounValue(a assert.Assertions, val literal.LiteralValue, noun string, path ...string) error {
	last := len(path) - 1
	field, parts := path[last], path[:last]
	return a.AssertNounValue(noun, field, parts, val)
}

func importNoun(w *weave.Weaver, noun string, nouns []string) (ret []string, err error) {
	if a, e := makeArticleName(w, noun); e != nil {
		err = e
	} else if a.count > 0 {
		if ns, e := importCountedNoun(w.Catalog, a.count, a.name); e != nil {
			err = e
		} else {
			ret = append(nouns, ns...)
		}
	} else {
		if a.isProper() {
			err = assertNounValue(w.Catalog, B(true), a.name, "proper_named")
		} else if customDet, ok := a.customArticle(); ok && len(customDet) > 0 {
			err = assertNounValue(w.Catalog, T(customDet), a.name, "indefinite_article")
		}
		ret = append(nouns, a.name)
	}
	return
}

type articleName struct {
	article, name string
	count         int
}

func makeArticleName(w *weave.Weaver, name string) (ret articleName, err error) {
	parts := strings.Fields(strings.TrimSpace(name))
	if last := len(parts) - 1; last == 0 {
		ret = articleName{
			name: parts[0],
		}
	} else if last > 0 {
		first, rest := parts[0], parts[1:]
		if count, counted := lang.WordsToNum(first); counted && count > 0 {
			ret = articleName{
				count: count,
				name:  strings.Join(rest, " "),
			}
		} else if cnt, e := w.MatchArticle(parts); e != nil {
			err = e
		} else if cnt == len(parts) {
			err = errutil.New("missing name from", name)
		} else {
			ret = articleName{
				article: strings.Join(parts[:cnt], " "),
				name:    strings.Join(parts[cnt:], " "),
			}
		}
	}
	return
}

// this isn't a correct test... but it will work for now...
func (an *articleName) isProper() (okay bool) {
	if n := an.name; len(n) > 1 || an.article == "our" {
		first, _ := utf8.DecodeRuneInString(n)
		okay = unicode.ToUpper(first) == first
	}
	return
}

func (an *articleName) customArticle() (ret string, okay bool) {
	switch an.article {
	case "a", "an", "the":
	default:
		ret, okay = an.article, true
	}
	return
}

// ex. "two triangles" -> triangle is a kind of thing
func importCountedNoun(cat *weave.Catalog, cnt int, kindOrKinds string) (ret []string, err error) {
	if cnt > 0 {
		// generate unique names for each of the counted nouns.
		// fix: we probably want nouns to "stack", and be have individually duplicated objects.
		// ie. a single stackable "cats" with a value of 5, rather than cat_1, cat_2, etc.
		// and when you pick up one cat now you have two object stacks, both referring to the kind cats
		// an empty stack acts like no object, and gets collected in some fashion.
		names := make([]string, cnt)
		for i := 0; i < cnt; i++ {
			names[i] = cat.NewCounter(kindOrKinds, nil)
		}
		if e := cat.Schedule(assert.RequirePlurals, func(w *weave.Weaver) (err error) {
			var kind, kinds string
			// note: kind is phrased in the singular here when count is 1, plural otherwise.
			if cnt == 1 {
				kind = kindOrKinds
				kinds = w.PluralOf(kindOrKinds)
			} else {
				kinds = kindOrKinds
				kind = w.SingularOf(kindOrKinds)
			}
			if e := cat.AssertAncestor(kinds, "thing"); e != nil {
				err = e
			} else {
				for _, n := range names {
					if e := cat.AssertNounKind(n, kindOrKinds); e != nil {
						err = e
					} else if e := cat.AssertAlias(n, kind); e != nil {
						err = e // ^ so that typing "triangle" means "triangles_1"
						break
					} else if e := assertNounValue(cat, B(true), n, "counted"); e != nil {
						err = e
						break
					} else if e := assertNounValue(cat, T(kind), n, "printed name"); e != nil {
						err = e // so that printing "triangles_1" yields "triangle"
						break   // FIX: itd make a lot more sense to have a default value for the kind
					}
				}
			}
			return
		}); e != nil {
			err = e
		} else {
			ret = names
		}
	}
	return
}
