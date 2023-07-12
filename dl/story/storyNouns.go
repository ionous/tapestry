package story

import (
	"errors"
	"unicode"
	"unicode/utf8"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/assert"
)

// helper to simplify setting the values of nouns
func assertNounValue(a assert.Assertions, val literal.LiteralValue, noun string, path ...string) error {
	last := len(path) - 1
	field, parts := path[last], path[:last]
	return a.AssertNounValue(noun, field, parts, val)
}

// this isn't a correct test... but it will work for now...
func isProper(article grok.Article, name string) (okay bool) {
	a := lang.Normalize(article.String())
	if len(name) > 1 || a == "our" {
		first, _ := utf8.DecodeRuneInString(name)
		okay = unicode.ToUpper(first) == first
	}
	return
}

func getCustomArticle(article grok.Article) (ret string) {
	switch a := lang.Normalize(article.String()); a {
	case "a", "an", "the":
	default:
		ret = a
	}
	return
}

func importNamedNoun(w *weave.Weaver, n grok.Noun) (ret string, err error) {
	var noun *weave.ScopedNoun
	og := n.Name.String()
	if name := lang.Normalize(og); name == "you" {
		// tdb: the current thought is that "the player" should be a variable;
		// currently its an "agent".
		noun, err = w.Domain.GetExactNoun("self")
	} else {
		if n.Exact { // ex. ".... called the spatula."
			noun, err = w.Domain.GetExactNoun(name)
		} else {
			// if it doesnt exist; we create it.
			if fold, e := w.Domain.GetClosestNoun(name); e != nil {
				err = e
			} else {
				noun = fold
			}
		}
		if errors.Is(err, mdl.Missing) {
			// ugh
			base := "things"
			if len(n.Kinds) > 0 {
				base = lang.Normalize(n.Kinds[0].String())
			}
			noun, err = w.Domain.AddNoun(og, name, base, w.At)
		}
	}
	// assign kinds
	if err == nil {
		for _, k := range n.Kinds {
			k := lang.Normalize(k.String())
			if e := w.Catalog.AddNoun(noun.Domain(), noun.Name(), k, w.At); e != nil && !errors.Is(e, mdl.Duplicate) {
				err = e
				break
			}
		}
	}
	// add articles:
	if err == nil {
		if isProper(n.Article, og) {
			if e := noun.WriteValue(w.At, "proper named", nil, B(true)); e != nil && !errors.Is(e, mdl.Duplicate) {
				err = e
			}
		} else if a := getCustomArticle(n.Article); len(a) > 0 {
			if e := noun.WriteValue(w.At, "indefinite article", nil, T(a)); e != nil && !errors.Is(e, mdl.Duplicate) {
				err = e
			}
		}
	}
	// add traits:
	if err == nil {
		err = assignTraits(w, noun, n.Traits)
	}
	// return
	if err == nil {
		ret = noun.Name()
	}
	return
}

func assignTraits(w *weave.Weaver, noun *weave.ScopedNoun, traits []grok.Match) (err error) {
	for _, t := range traits {
		// FIX: this passes through "GetClosestNoun" which seems wrong here.
		// the issue is the noun might not exist;
		// so we'd have to break some of this open to handle it.
		if e := noun.WriteValue(w.At, t.String(), nil, B(true)); e != nil && !errors.Is(e, mdl.Duplicate) {
			err = e
			break
		}
	}
	return
}

// ex. "two triangles"
// - adds ( and returns ) nouns: triangle_1, triangle_2, etc. of kind "triangle/s"
// - uses "triangle" as an alias and printed name for each of the new nouns
// - flags them all as "counted.
// - ensures "triangle/s" are things
func importCountedNoun(cat *weave.Catalog, noun grok.Noun) (ret []string, err error) {
	// ..kindOrKinds string, article grok.Article, traits []grok.Match
	if cnt := noun.Article.Count; cnt > 0 {
		kindOrKinds := lang.Normalize(noun.Kinds[0].String())
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
					if n, e := w.Domain.AddNoun(n, n, kindOrKinds, w.At); e != nil {
						err = e
					} else if e := cat.AddName(n.Domain(), n.Name(), kind, -1, w.At); e != nil {
						err = e // ^ so that typing "triangle" means "triangles-1"
						break
					} else if e := n.WriteValue(w.At, "counted", nil, B(true)); e != nil {
						err = e
						break
					} else if e := n.WriteValue(w.At, "printed name", nil, T(kind)); e != nil {
						err = e // so that printing "triangles-1" yields "triangle"
						break   // FIX: itd make a lot more sense to have a default value for the kind
					} else if e := assignTraits(w, n, noun.Traits); e != nil {
						err = e
						break
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
