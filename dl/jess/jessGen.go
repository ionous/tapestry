package jess

import (
	"errors"
	"unicode"
	"unicode/utf8"

	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

func applyResults(rar Registrar, res localResults) error {
	return genNouns(rar, res.Primary, res.Secondary, func(src, tgt []string) (err error) {
		// note: some phrases "the box is open" dont have macros.
		// in that case, genNouns does all the work.
		if macro := res.Macro.Name; len(macro) > 0 {
			err = rar.Apply(res.Macro, src, tgt)
		}
		return
	})
}

func startsUpper(str string) bool {
	first, _ := utf8.DecodeRuneInString(str)
	return unicode.IsUpper(first) // this works okay even if the string was empty
}

// even one name can generate several nouns ( ex. "two things" )
// after gets called for each one.1
func genNoun(rar Registrar, n DesiredNoun, after postGenOne) error {
	return rar.PostProcess(NounSettings, func(Query, Registrar) (err error) {
		if names, e := importNoun(rar, n, nil); e != nil {
			err = e
		} else {
			for _, name := range names {
				if e := after(name); e != nil {
					err = e
					break
				}
			}
		}
		return
	})
}

type postGenOne func(a string) error
type postGenMany func(a, b []string) error

// add nouns and values
func genNouns(rar Registrar, a, b []DesiredNoun, after postGenMany) error {
	return rar.PostProcess(NounSettings, func(Query, Registrar) (err error) {
		if src, e := importNouns(rar, a); e != nil {
			err = e
		} else if tgt, e := importNouns(rar, b); e != nil {
			err = e
		} else {
			err = after(src, tgt)
		}
		return
	})
}

func importNouns(rar Registrar, ns []DesiredNoun) (ret []string, err error) {
	if cnt := len(ns); cnt > 0 {
		ret = make([]string, 0, cnt)
		for _, n := range ns {
			if ret, err = importNoun(rar, n, ret); err != nil {
				break
			}
		}
	}
	return
}

func importNoun(rar Registrar, n DesiredNoun, names []string) (ret []string, err error) {
	if n.Count > 0 {
		if res, e := importCountedNoun(rar, n); e != nil {
			err = e
		} else {
			ret = append(names, res...)
		}
	} else {
		if res, e := importNamedNoun(rar, n); e != nil {
			err = e
		} else {
			ret = append(names, res)
		}
	}
	return
}

func importNamedNoun(rar Registrar, n DesiredNoun) (ret string, err error) {
	if noun := n.Noun; len(noun) > 0 {
		if e := registerKinds(rar, noun, n.Kinds); e != nil {
			err = e
		} else if e := registerTraits(rar, noun, n.Traits); e != nil {
			err = e
		} else {
			ret = noun
		}
	} else {
		noun := inflect.Normalize(n.DesiredName)
		kinds, traits := n.Kinds, n.Traits
		if len(kinds) == 0 {
			kinds = []string{Things} // default
		}
		// the lack of a recognized article makes something proper-named.
		if len(n.Article) == 0 {
			traits = append(traits, ProperNameTrait)
		}
		if e := registerKinds(rar, noun, kinds); e != nil {
			err = e
		} else if e := registerNames(rar, noun, n.DesiredName); e != nil {
			err = e
		} else if e := registerTraits(rar, noun, traits); e != nil {
			err = e
		} else if e := registerArticle(rar, noun, n.Article, n.Flags); e != nil {
			err = e
		} else {
			ret = noun
		}
	}
	return
}

func registerArticle(rar Registrar, noun, a string, f ArticleFlags) (err error) {
	if f.Plural {
		err = rar.AddNounTrait(noun, PluralNamed)
	}
	if f.Indefinite && err == nil {
		err = rar.AddNounValue(noun, IndefiniteArticle, text(a, ""))
	}
	return
}

// ensure the noun inherits from all of the specified kinds
func registerKinds(rar Registrar, noun string, kinds []string) (err error) {
	// fix? it's a little hard in some cases to filter out common kinds (eg. rooms)
	// it doesnt really matter -- AddNounKind can handle it; but it makes the tests ugly.
	dedupe := make(map[string]bool)
	for _, k := range kinds {
		if !dedupe[k] {
			dedupe[k] = true
			if e := rar.AddNounKind(noun, k); e != nil && !errors.Is(e, mdl.Duplicate) {
				err = e
				break
			}
		}
	}
	return
}

func registerNames(rar Registrar, noun, name string) (err error) {
	names := mdl.MakeNames(name)
	for i, n := range names {
		if e := rar.AddNounName(noun, n, i); e != nil {
			err = e
			break
		}
	}
	return
}

func registerTraits(rar Registrar, noun string, traits []string) (err error) {
	for _, t := range traits {
		if e := rar.AddNounTrait(noun, t); e != nil {
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
func importCountedNoun(rar Registrar, noun DesiredNoun) (ret []string, err error) {
	if cnt := noun.Count; cnt > 0 {
		kinds := noun.Kinds[0]
		kind := rar.GetSingular(kinds)
		names := make([]string, cnt)
	Loop:
		for i := 0; i < cnt; i++ {
			n := rar.GetUniqueName(kind)
			names[i] = n
			if e := rar.AddNounKind(n, kinds); e != nil {
				err = e
			} else if e := rar.AddNounName(n, n, 0); e != nil {
				err = e // ^ so authors can refer to it by the dashed name
				break
			} else if e := rar.AddNounName(n, kind, -1); e != nil {
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
					if e := rar.AddNounTrait(n, t); e != nil {
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
