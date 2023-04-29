package story

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/assert"
)

// ImportNounProperties -
// reads ancillary information about nouns from their names and declares properties for them.
// ex. proper or plural names, etc.
func ImportNounProperties(cat *weave.Catalog, nouns []string) (ret []string, err error) {
	ret = make([]string, 0, len(nouns))
	for _, noun := range nouns {
		if next, e := importNoun(cat, noun, ret); e != nil {
			err = e
			break
		} else {
			ret = next
		}
	}
	return
}

// ReadNouns - reads noun names without declaring any properties....
// fix? unless they are counted nouns ( for backwards compatibility of "two cats whereabouts the kitchen" )
func ReadNouns(cat *weave.Catalog, nouns []string) (ret []string, err error) {
	ret = make([]string, 0, len(nouns))
	for _, noun := range nouns {
		if a := makeArticleName(noun); a.count == 0 {
			ret = append(ret, a.name)
		} else if ns, e := importCountedNoun(cat, a.count, a.name); e == nil {
			ret = append(ret, ns...)
		} else {
			err = e
			break
		}
	}
	return
}

func importNoun(cat *weave.Catalog, noun string, nouns []string) (ret []string, err error) {
	if a := makeArticleName(noun); a.count > 0 {
		if ns, e := importCountedNoun(cat, a.count, a.name); e != nil {
			err = e
		} else {
			ret = append(nouns, ns...)
		}
	} else {
		if a.isProper() {
			err = assert.AssertNounValue(cat, B(true), a.name, "proper_named")
		} else if customDet, ok := a.customArticle(); ok && len(customDet) > 0 {
			err = assert.AssertNounValue(cat, T(customDet), a.name, "indefinite_article")
		}
		ret = append(nouns, a.name)
	}
	return
}

type articleName struct {
	article, name string
	count         int
}

// fix: this will never be correct until the indefinite articles are driven by script, parsed first, and the noun munging is in weave.
func makeArticleName(name string) (ret articleName) {
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
		} else {
			// tried using "of" to grab mass nouns; but it doesnt work well for "a can of soup"
			//split := 1
			//for i, s := range parts {
			//	if s == "of" && i < len(parts)-1 {
			//		first = strings.Join(parts[:i+1], " ")
			//		split = i + 1
			//		break
			//	}
			//}
			ret = articleName{
				article: first,
				name:    strings.Join(parts[1:], " "),
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
// fix? consider a specific counted noun phrase; the noun phrase needs more work.
// also, we probably want noun stacks not individually duplicated names
func importCountedNoun(cat *weave.Catalog, cnt int, kindOrKinds string) (ret []string, err error) {
	if cnt > 0 {
		// note: kind is phrased in the singular here when count is 1, plural otherwise.
		// but, because of "Recent.Nouns" processing we have to generate some sort of noun name *immediately*
		// ( itd be nice to have a more start and stop importer, where we could delay processing of branches of the tree. )
		names := make([]string, cnt)
		for i := 0; i < cnt; i++ {
			noun := cat.NewCounter(kindOrKinds, nil)
			if e := assert.AssertNounValue(cat, B(true), noun, "counted"); e != nil {
				err = e
				break
			}
			names[i] = noun
		}
		if e := cat.Schedule(assert.AncestryPhase, func(w *weave.Weaver) (err error) {
			// by now, plurals will be determined, so we can determine which is which.
			var kind, kinds string
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
					} else if e := assert.AssertNounValue(cat, T(kind), n, "printed name"); e != nil {
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
