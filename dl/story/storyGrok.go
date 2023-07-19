package story

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/ionous/errutil"

	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// Execute - called by the macro runtime during weave.
func (op *DeclareStatement) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DeclareStatement) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequireRules, func(w *weave.Weaver) (err error) {
		if text, e := safe.GetText(w, op.Text); e != nil {
			err = e
		} else {
			if res, e := w.Grok(text.String()); e != nil {
				err = e
			} else if strings.HasPrefix(res.Macro.Name, "inherit") {
				// hack for handling "kinds of"
				// alt, could maybe look at the number of fields of the macro
				// it'd be amazing to handle all of validate and genNouns in macros
				// (so no decisions are needed here ) but dont think that's easily possible.
				err = grokKindPhrase(w, res)
			} else {
				err = grokNounPhrase(w, res)
			}
			//
			if err != nil {
				err = errutil.Fmt("%w grokking %q", err, text.String())
			}
		}
		return
	})
}

func grokKindPhrase(w *weave.Weaver, res grok.Results) (err error) {
	if len(res.Targets) != 0 {
		err = errutil.Fmt("%s only expected sources", res.Macro.Name)
	} else {
		pen := w.Pin()
		for _, src := range res.Sources {
			// two forms: one where the kind is already known to be a kind;
			// and one where its seen as a generic name.
			if grok.MatchLen(src.Name) == 0 {
				if len(src.Kinds) != 2 {
					err = errutil.Fmt("%s expected a single kind per source", res.Macro.Name)
				} else {
					kind := lang.Normalize(src.Kinds[0].String())
					ancestor := lang.Normalize(src.Kinds[1].String())
					if e := pen.AddKind(kind, ancestor); e != nil {
						err = e
						break
					}
				}
			} else {
				if len(src.Kinds) != 1 {
					err = errutil.Fmt("%s expected a single kind per source", res.Macro.Name)
				} else {
					kind := lang.Normalize(src.Name.String())
					ancestor := lang.Normalize(src.Kinds[0].String())
					if e := pen.AddKind(kind, ancestor); e != nil {
						err = e
						break
					}
				}
			}
		}
	}
	return
}

func grokNounPhrase(w *weave.Weaver, res grok.Results) (err error) {
	if multiSrc, e := validSources(res.Sources, res.Macro.Type); e != nil {
		err = e
	} else if multiTgt, e := validTargets(res.Targets, res.Macro.Type); e != nil {
		err = e
	} else if src, e := genNouns(w, res.Sources, multiSrc); e != nil {
		err = e
	} else if tgt, e := genNouns(w, res.Targets, multiTgt); e != nil {
		err = e
	} else {
		// note: some phrases "the box is open" dont have macros.
		// in that case, genNouns does all the work.
		if macro := res.Macro.Name; len(macro) > 0 {
			if kind, e := w.GetKindByName(macro); e != nil {
				err = e
			} else if !kind.Implements(kindsOf.Macro.String()) {
				err = errutil.Fmt("expected %q to be a macro", macro)
			} else if fieldCnt := kind.NumField(); fieldCnt < 2 {
				err = errutil.Fmt("expected macro %q to have at least two argument (not %d)", macro, fieldCnt)
			} else {
				rec := kind.NewRecord()
				args := []g.Value{src, tgt}
				for i := 0; i < 2; i++ {
					if val, e := safe.AutoConvert(w, kind.Field(i), args[i]); e != nil {
						err = e
						break
					} else if e := rec.SetIndexedField(i, val); e != nil {
						// note: set indexed field assigns without copying
						// but get value copies out, so this should be okay.
						err = errutil.Fmt("%w while setting %q arg %v", e, macro, i)
						break
					}
				}
				if err == nil {
					if v, e := w.Call(rec, affine.Text); e != nil && !errors.Is(e, rt.NoResult) {
						err = e
					} else if v != nil {
						if msg := v.String(); len(msg) > 0 {
							err = errutil.Fmt("Declare statement: %s", msg)
						}
					}
				}
			}
		}
	}
	return
}

// for logging errors
func reduceNouns(ns []grok.Noun) (ret string) {
	var s strings.Builder
	for i, n := range ns {
		if i > 1 {
			s.WriteString(", ")
		}
		s.WriteString(n.Name.String())

	}
	return s.String()
}

func validSources(ns []grok.Noun, mtype grok.MacroType) (multi bool, err error) {
	switch mtype {
	case grok.Macro_SourcesOnly, grok.Macro_ManySources, grok.Macro_ManyMany:
		if cnt := len(ns); cnt == 0 {
			err = errutil.New("expected at least one source noun")
		}
		multi = true
	case grok.Macro_ManyTargets:
		if cnt := len(ns); cnt > 1 {
			err = errutil.New("expected exactly one noun, have:", reduceNouns(ns))
		}
	default:
		err = errutil.New("invalid macro type")
	}
	return
}

func validTargets(ns []grok.Noun, mtype grok.MacroType) (multi bool, err error) {
	switch mtype {
	case grok.Macro_SourcesOnly:
		if cnt := len(ns); cnt != 0 {
			err = errutil.New("didn't expect any target nouns")
		}
	case grok.Macro_ManySources:
		if cnt := len(ns); cnt > 1 {
			err = errutil.New("expected at most one target noun")
		}
	case grok.Macro_ManyTargets, grok.Macro_ManyMany:
		// any number okay
		multi = true
	default:
		err = errutil.New("invalid macro type")
	}
	return
}

// determine whether the noun seems to be a proper name
func isProper(article grok.Article, name string) (okay bool) {
	a := lang.Normalize(article.String())
	if len(name) > 1 || a == "our" {
		first, _ := utf8.DecodeRuneInString(name)
		okay = unicode.ToUpper(first) == first
	}
	return
}

// determine whether the noun will need a custom indefinite property
// this uses a subset of the known articles, due to way object printing works.
func getCustomArticle(article grok.Article) (ret string) {
	switch a := lang.Normalize(article.String()); a {
	case "a", "an", "the":
	default:
		ret = a
	}
	return
}

// add nouns and values
func genNouns(w *weave.Weaver, ns []grok.Noun, multi bool) (ret g.Value, err error) {
	names := make([]string, 0, len(ns))
	for _, n := range ns {
		if n.Article.Count > 0 {
			if ns, e := importCountedNoun(w.Catalog, n); e != nil {
				err = e
				break
			} else {
				names = append(names, ns...)
			}
		} else {
			if name, e := importNamedNoun(w, n); e != nil {
				err = e
				break
			} else {
				names = append(names, name)
			}
		}
	} // end for loop
	// all done?
	if err == nil {
		if multi {
			ret = g.StringsOf(names)
		} else if len(names) > 0 {
			ret = g.StringOf(names[0])
		} else {
			ret = g.Empty
		}
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
			// note: this can return nil if it existed
			// but since we just checked for it; things in single threading are fine.
			noun, err = w.Domain.AddNoun(og, name, base)
		}
	}
	// assign kinds
	if err == nil {
		pen := w.Pin()
		for _, k := range n.Kinds {
			// fix dont listen to any errors since we cant filter dupes here
			k := lang.Normalize(k.String())
			pen.AddNoun(noun.Name(), k)
		}
	}
	// add articles:
	if err == nil {
		if isProper(n.Article, og) {
			if e := noun.WriteValue(w.At, "proper named", nil, B(true)); e != nil {
				err = e
			}
		} else if a := getCustomArticle(n.Article); len(a) > 0 {
			if e := noun.WriteValue(w.At, "indefinite article", nil, T(a)); e != nil {
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
		if e := cat.Schedule(weave.RequirePlurals, func(w *weave.Weaver) (err error) {
			var kind, kinds string
			// note: kind is phrased in the singular here when count is 1, plural otherwise.
			if cnt == 1 {
				kind = kindOrKinds
				kinds = w.PluralOf(kindOrKinds)
			} else {
				kinds = kindOrKinds
				kind = w.SingularOf(kindOrKinds)
			}
			if e := w.Pin().AddKind(kinds, "thing"); e != nil {
				err = e
			} else {
				pen := w.Pin()
				for _, n := range names {
					if n, e := w.Domain.AddNoun(n, n, kindOrKinds); e != nil {
						err = e
					} else if e := pen.AddName(n.Name(), kind, -1); e != nil {
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

func assignTraits(w *weave.Weaver, noun *weave.ScopedNoun, traits []grok.Match) (err error) {
	for _, t := range traits {
		// FIX: this passes through "GetClosestNoun" which seems wrong here.
		// the issue is the noun might not exist;
		// so we'd have to break some of this open to handle it.
		if e := noun.WriteValue(w.At, t.String(), nil, B(true)); e != nil {
			err = e
			break
		}
	}
	return
}
