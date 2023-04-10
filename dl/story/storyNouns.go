package story

import (
	"git.sr.ht/~ionous/tapestry/imp/assert"
	"strings"
	"unicode"
	"unicode/utf8"

	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
)

// ImportNounProperties -
// reads ancillary information about nouns from their names and declares properties for them.
// ex. proper or plural names, etc.
func ImportNounProperties(k *imp.Importer, nouns []string) (ret []string, err error) {
	// declare a noun class that has several default fields
	// should be moved to script.
	if once := "noun"; k.Once(once) {
		k.WriteOnce(&eph.EphKinds{Kind: "objects", Ancestor: kindsOf.Kind.String()})
		// common or proper nouns ( rabbit, vs. Roger )
		k.AddImplicitAspect("noun_types", "objects", "common_named", "proper_named", "counted")
		// whether a player can refer to an object by its name.
		k.AddImplicitAspect("private_names", "objects", "publicly_named", "privately_named")
		k.AddImplicitAspect("private_names", "objects", "publicly_named", "privately_named")
	}
	ret = make([]string, 0, len(nouns))
	for _, noun := range nouns {
		if next, e := importNoun(k, noun, ret); e != nil {
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
func ReadNouns(k *imp.Importer, nouns []string) (ret []string, err error) {
	ret = make([]string, 0, len(nouns))
	for _, noun := range nouns {
		if a := makeArticleName(noun); a.count == 0 {
			ret = append(ret, a.name)
		} else if ns, e := importCountedNoun(k, a.count, a.name); e == nil {
			ret = append(ret, ns...)
		} else {
			err = e
			break
		}
	}
	return
}

func importNoun(k *imp.Importer, noun string, nouns []string) (ret []string, err error) {
	if a := makeArticleName(noun); a.count > 0 {
		if ns, e := importCountedNoun(k, a.count, a.name); e != nil {
			err = e
		} else {
			ret = append(nouns, ns...)
		}
	} else {
		if a.isProper() {
			k.WriteEphemera(&eph.EphValues{Noun: a.name, Field: "proper_named", Value: B(true)})

		} else if customDet, ok := a.customArticle(); ok {
			// setup the indefinite article
			// create a "indefinite article" field for all objects
			if k.Once("named_noun") {
				k.WriteOnce(&eph.EphKinds{Kind: "objects", Contain: []eph.EphParams{{Name: "indefinite_article", Affinity: eph.Affinity{eph.Affinity_Text}}}})
			}
			// set the indefinite article field for this object.
			k.WriteEphemera(&eph.EphValues{Noun: a.name, Field: "indefinite_article", Value: T(customDet)})
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
func importCountedNoun(k *imp.Importer, cnt int, kindOrKinds string) (names []string, err error) {
	if once := "printed_name"; k.Once(once) {
		k.WriteOnce(&eph.EphKinds{Kind: "objects", Contain: []eph.EphParams{{Name: "printed_name", Affinity: eph.Affinity{Str: eph.Affinity_Text}}}})
	}
	if cnt > 0 {
		// note: kind is phrased in the singular here when count is 1, plural otherwise.
		// but, because of "Recent.Nouns" processing we have to generate some sort of noun name *immediately*
		// ( itd be nice to have a more start and stop importer, where we could delay processing of branches of the tree. )
		names = make([]string, cnt)
		for i := 0; i < cnt; i++ {
			noun := k.NewCounter(kindOrKinds, nil)
			k.WriteEphemera(&eph.EphValues{Noun: noun, Field: "counted", Value: B(true)})
			names[i] = noun
		}
		k.WriteEphemera(eph.PhaseFunction{OnPhase: assert.AncestryPhase,
			Do: func(c *eph.Catalog, d *eph.Domain, at string) (err error) {
				// by now, plurals will be determined, so we can determine which is which.
				var kind, kinds string
				if cnt == 1 {
					kind = kindOrKinds
					kinds, err = d.Pluralize(kindOrKinds)
				} else {
					kinds = kindOrKinds
					kind, err = d.Singularize(kindOrKinds)
				}
				if err == nil {
					if e := d.AddEphemera(at, &eph.EphKinds{Kind: kinds, Ancestor: "thing"}); e != nil {
						err = e
					} else {
						for _, n := range names {
							if e := d.AddEphemera(at, &eph.EphNouns{Noun: n, Kind: kindOrKinds}); e != nil {
								err = e
							} else if e := d.AddEphemera(at, &eph.EphAliases{
								// so that typing "triangle" means "triangles_1"
								ShortName: n, Aliases: []string{kind},
							}); e != nil {
								err = e
								break
							} else if e := d.AddEphemera(at, &eph.EphValues{
								// so that printing "triangles_1" yields "triangle"
								// FIX: itd make a lot more sense to have a default value for the kind
								Noun: n, Field: "printed_name", Value: T(kind),
							}); e != nil {
								err = e
								break
							}
						}
					}
				}
				return
			}},
		)
	}
	return
}
