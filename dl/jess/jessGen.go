package jess

import (
	"errors"
	"unicode"
	"unicode/utf8"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/ionous/errutil"
)

func grokNounPhrase(rar Registrar, res localResults) (err error) {
	if src, e := genNouns(rar, res.Primary); e != nil {
		err = e
	} else if tgt, e := genNouns(rar, res.Secondary); e != nil {
		err = e
	} else {
		// note: some phrases "the box is open" dont have macros.
		// in that case, genNouns does all the work.
		if macro := res.Macro.Name; len(macro) > 0 {
			err = rar.Apply(res.Macro, src, tgt)
		}
	}
	return
}

// tbd: add this as a factory function to Registrar?
func text(value, kind string) rt.Assignment {
	return &assign.FromText{
		Value: &literal.TextValue{Value: value, Kind: kind},
	}
}

// determine whether the noun seems to be a proper name
func isProper(article articleResult, name string) (okay bool) {
	a := inflect.Normalize(article.String())
	if len(name) > 1 || a == "our" {
		first, _ := utf8.DecodeRuneInString(name)
		okay = unicode.ToUpper(first) == first
	}
	return
}

// determine whether the noun will need a custom indefinite property
// this uses a subset of the known articles, due to way object printing works.
func getCustomArticle(article articleResult) (ret string) {
	switch a := inflect.Normalize(article.String()); a {
	case "a", "an", "the":
	default:
		ret = a
	}
	return
}

// add nouns and values
func genNouns(rar Registrar, ns []resultName) (ret []string, err error) {
	names := make([]string, 0, len(ns))
	for _, n := range ns {
		if n.Article.Count > 0 {
			if ns, e := importCountedNoun(rar, n); e != nil {
				err = e
				break
			} else {
				names = append(names, ns...)
			}
		} else {
			if name, e := importNamedNoun(rar, n); e != nil {
				err = e
				break
			} else {
				names = append(names, name)
			}
		}
	}
	if err == nil {
		ret = names
	}
	return
}

func importNamedNoun(rar Registrar, n resultName) (ret string, err error) {
	var noun string
	fullName := n.String()
	if name := inflect.Normalize(fullName); name == "you" {
		// tdb: the current thought is that "the player" should be a variable;
		// currently its an "agent".
		noun, err = rar.GetExactNoun("self")
	} else {
		if n.Exact { // ex. ".... called the spatula."
			noun, err = rar.GetExactNoun(name)
		} else if fold, e := rar.GetClosestNoun(name); e != nil {
			err = e
		} else {
			noun = fold
		}
		// if it doesnt exist; we create it.
		if errors.Is(err, mdl.Missing) {
			base := "things" // ugh
			if len(n.Kinds) > 0 {
				base = inflect.Normalize(n.Kinds[0].String())
			}
			if e := rar.AddNoun(name, fullName, base); e != nil {
				err = e
			} else {
				noun = name
			}
		}
	}
	// assign kinds
	// fix consider a "noun builder" instead
	if err == nil {
		for _, k := range n.Kinds {
			k := inflect.Normalize(k.String())
			// since noun already exists: this ensures that the noun inherits from all of the specified kinds
			if e := rar.AddNoun(noun, "", k); e != nil {
				err = e
				break
			}
		}
	}
	// add articles:
	if err == nil {
		if isProper(n.Article, fullName) {
			if e := rar.AddNounTrait(noun, "proper named"); e != nil {
				err = e
			}
		} else if a := getCustomArticle(n.Article); len(a) > 0 {
			if e := rar.AddNounValue(noun, "indefinite article", text(a, "")); e != nil {
				err = e
			}
		}
	}
	// add traits:
	if err == nil {
		for _, t := range n.Traits {
			t := inflect.Normalize(t.String())
			if e := rar.AddNounTrait(noun, t); e != nil {
				err = errutil.Append(err, e)
				break // out of the traits to the next noun
			}
		}
	}
	// return
	if err == nil {
		ret = noun
	}
	return
}

// ex. "two triangles"
// - adds ( and returns ) nouns: triangle_1, triangle_2, etc. of kind "triangle/s"
// - uses "triangle" as an alias and printed name for each of the new nouns
// - flags them all as "counted.
// - ensures "triangle/s" are things
func importCountedNoun(rar Registrar, noun resultName) (ret []string, err error) {
	// ..kindOrKinds string, article ArticleResult, traits []match.Match
	if cnt := noun.Article.Count; cnt > 0 {
		// generate unique names for each of the counted nouns.
		// fix: we probably want nouns to "stack", and be have individually duplicated objects.
		// ie. a single stackable "cats" with a value of 5, rather than cat_1, cat_2, etc.
		// and when you pick up one cat now you have two object stacks, both referring to the kind cats
		// an empty stack acts like no object, and gets collected in some fashion.
		name, parent := noun.String(), "thing"
		if len(name) == 0 {
			name = noun.Kinds[0].String()
		} else if len(noun.Kinds) > 0 {
			// ex. ""An empire apple, a pen, and two triangles are props in the lab."
			// fix: grok should return that as an object *called* two triangles, not something counted.
			parent = noun.Kinds[0].String()
		}
		name = inflect.Normalize(name)

		names := make([]string, cnt)
		for i := 0; i < cnt; i++ {
			names[i] = rar.GetUniqueName(name)
		}

		var kind, kinds string
		// note: kind is phrased in the singular here when count is 1, plural otherwise.
		if cnt == 1 {
			kind = name
			kinds = rar.GetPlural(name)
		} else {
			kinds = name
			kind = rar.GetSingular(name)
		}
		if e := rar.AddKind(kinds, parent); e != nil {
			err = e
		} else {
		Loop:
			for _, n := range names {
				if e := rar.AddNoun(n, n, kinds); e != nil {
					err = e
				} else if e := rar.AddNounAlias(n, kind, -1); e != nil {
					err = e // ^ so that typing "triangle" means "triangles-1"
					break
				} else if e := rar.AddNounTrait(n, "counted"); e != nil {
					err = e
					break
				} else if e := rar.AddNounValue(n, "printed name", text(kind, "")); e != nil {
					err = e // so that printing "triangles-1" yields "triangle"
					break   // FIX: itd make a lot more sense to have a default value for the kind
				} else {
					for _, t := range noun.Traits {
						if e := rar.AddNounTrait(n, t.String()); e != nil {
							err = e
							break Loop
						}
					}
				}
			}
			if err == nil {
				ret = names
			}
		}
	}
	return
}
