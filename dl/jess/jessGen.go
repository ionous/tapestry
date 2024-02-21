package jess

import (
	"errors"
	"unicode"
	"unicode/utf8"

	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/ionous/errutil"
)

// todo: remove. it should no longer be necessary to compile local results
// before applying them; the various generate functions should be able to use Registrar directly
func applyResults(rar Registrar, res localResults) (err error) {
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

// fix: it doesnt seem like this should be *reading* nouns:
// that should be part of matching; recording whatever it needs
// then output the results here.
func importNamedNoun(rar Registrar, n resultName) (ret string, err error) {
	var noun string
	fullName := n.String()
	if name := inflect.Normalize(fullName); name == PlayerYou {
		// tdb: the current thought is that "the player" should be a variable;
		// currently its an "agent".
		noun, err = rar.GetExactNoun(PlayerSelf)
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
			base := DefaultKind // ugh
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
			if e := rar.AddNounTrait(noun, ProperNameTrait); e != nil {
				err = e
			}
		} else if a := getCustomArticle(n.Article); len(a) > 0 {
			if e := rar.AddNounValue(noun, IndefiniteArticle, text(a, "")); e != nil {
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
func importCountedNoun(rar Registrar, noun resultName) (ret []string, err error) {
	if cnt := noun.Article.Count; cnt > 0 {
		kinds := noun.Kinds[0].String()
		kind := rar.GetSingular(kinds)
		names := make([]string, cnt)
	Loop:
		for i := 0; i < cnt; i++ {
			n := rar.GetUniqueName(kinds)
			names[i] = n
			if e := rar.AddNoun(n, n, kinds); e != nil {
				err = e
			} else if e := rar.AddNounAlias(n, kind, -1); e != nil {
				err = e // ^ so that typing "triangle" means "triangles-1"
				break
			} else if e := rar.AddNounTrait(n, CountedTrait); e != nil {
				err = e
				break
			} else if e := rar.AddNounValue(n, PrintedName, text(kind, "")); e != nil {
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
			if err == nil {
				ret = names
			}
		}
	}
	return
}
