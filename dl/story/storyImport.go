package story

import (
	"errors"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"

	"git.sr.ht/~ionous/tapestry/dl/grammar"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// Execute - called by the macro runtime during weave.
func (op *DefineTraits) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// (the) colors are red, blue, or green.
func (op *DefineTraits) Weave(cat *weave.Catalog) error {
	return cat.Schedule(assert.RequireDeterminers, func(w *weave.Weaver) (err error) {
		if traits, e := safe.GetTextList(w, op.Traits); e != nil {
			err = e
		} else if aspect, e := safe.GetText(w, op.Aspect); e != nil {
			err = e
		} else {
			err = cat.AssertAspectTraits(aspect.String(), traits.Strings())
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *GrammarDecl) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *GrammarDecl) Weave(cat *weave.Catalog) (err error) {
	switch el := op.Grammar.(type) {
	// fix: why have a generic "grammar" decl, just to switch on two sub decls
	// they should be top level.
	case *grammar.Alias:
		err = cat.AssertAlias(el.AsNoun, el.Names...)
	case *grammar.Directive:
		name := strings.Join(el.Lede, "/")
		err = cat.AssertGrammar(name, el)
	default:
		err = errutil.Fmt("unknown grammar %T", el)
	}
	return
}

// Execute - called by the macro runtime during weave.
func (op *DefineNounTraits) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefineNounTraits) Weave(cat *weave.Catalog) error {
	return cat.Schedule(assert.RequireDeterminers, func(w *weave.Weaver) (err error) {
		if nouns, e := safe.GetTextList(w, op.Nouns); e != nil {
			err = e
		} else if kind, e := safe.GetOptionalText(w, op.Kind, ""); e != nil {
			err = e
		} else if traits, e := safe.GetTextList(w, op.Traits); e != nil {
			err = e
		} else {
			names := nouns.Strings()
			if kind, e := grok.StripArticle(kind.String()); e != nil {
				err = e
			} else if len(kind) > 0 {
				for i, name := range names {
					if name, e := grok.StripArticle(name); e != nil {
						err = errutil.Append(err, e)
					} else {
						names[i] = name // replace for the traits loop
						if e := cat.AssertNounKind(name, kind); e != nil {
							err = errutil.Append(err, e)
						}
					}
				}
			}
			if traits := traits.Strings(); len(traits) > 0 {
				for _, t := range traits {
					t := lang.Normalize(t)
					for _, n := range names {
						if e := assertNounValue(cat, B(true), n, t); e != nil {
							err = errutil.Append(err, e)
							break // out of the traits to the next noun
						}
					}
				}
			}
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *DefinePhrase) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefinePhrase) Weave(cat *weave.Catalog) error {
	return cat.Schedule(assert.RequireAncestry, func(w *weave.Weaver) (err error) {
		if phrase, e := safe.GetText(w, op.Phrase); e != nil {
			err = e
		} else if macro, e := safe.GetText(w, op.Macro); e != nil {
			err = e
		} else if rev, e := safe.GetOptionalBool(w, op.Reversed, false); e != nil {
			err = e
		} else if macro := lang.Normalize(macro.String()); len(macro) == 0 {
			err = errutil.New("missing macro name")
		} else {
			d, at := w.Domain.Name(), w.At
			err = cat.AddPhrase(d, macro, phrase.String(), rev.Bool(), at)
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *DeclareStatement) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DeclareStatement) Weave(cat *weave.Catalog) error {
	return cat.Schedule(assert.RequireRules, func(w *weave.Weaver) (err error) {
		if text, e := safe.GetText(w, op.Text); e != nil {
			err = e
		} else if e := op.grok(w, text.String()); e != nil {
			err = errutil.Fmt("%w grokking %q", e, text.String())
		}
		return
	})
}

func (op *DeclareStatement) grok(w *weave.Weaver, text string) (err error) {
	if res, e := w.Grok(text); e != nil {
		err = e
	} else if multiSrc, e := validSources(res.Sources, res.Macro.Type); e != nil {
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
			if name, e := importNamedNoun(w, n); e != nil && !errors.Is(e, mdl.Duplicate) {
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

// Execute - called by the macro runtime during weave.
func (op *DefineNouns) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefineNouns) Weave(cat *weave.Catalog) error {
	return cat.Schedule(assert.RequireDeterminers, func(w *weave.Weaver) (err error) {
		if nouns, e := safe.GetTextList(w, op.Nouns); e != nil {
			err = e
		} else if kind, e := safe.GetText(w, op.Kind); e != nil {
			err = e
		} else {
			names := nouns.Strings()
			if kind := kind.String(); len(kind) > 0 {
				if kind, e := grok.StripArticle(kind); e != nil {
					err = e
				} else {
					for _, noun := range names {
						if noun, e := grok.StripArticle(noun); e != nil {
							err = errutil.Append(err, e)
						} else if e := cat.AssertNounKind(noun, kind); e != nil {
							err = errutil.Append(err, e)
						}
					}
				}
			}
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *NounAssignment) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// ex. The description of the nets is xxx
func (op *NounAssignment) Weave(cat *weave.Catalog) error {
	return cat.Schedule(assert.RequireDeterminers, func(w *weave.Weaver) (err error) {
		if nouns, e := safe.GetTextList(w, op.Nouns); e != nil {
			err = e
		} else if field, e := safe.GetText(w, op.FieldName); e != nil {
			err = e
		} else if lines, e := ConvertText(op.Lines.String()); e != nil {
			err = e
		} else {
			subjects := nouns.Strings()
			field, lines := field.String(), T(lines)
			for _, noun := range subjects {
				if noun, e := grok.StripArticle(noun); e != nil {
					err = errutil.Append(err, e)
				} else if e := assertNounValue(cat, lines, noun, field); e != nil {
					err = errutil.Append(err, e)
				}
			}
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *DefineRelatives) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefineRelatives) Weave(cat *weave.Catalog) error {
	return cat.Schedule(assert.RequireDeterminers, func(w *weave.Weaver) (err error) {
		if nouns, e := safe.GetTextList(w, op.Nouns); e != nil {
			err = e
		} else if kind, e := safe.GetOptionalText(w, op.Kind, ""); e != nil {
			err = e
		} else if relation, e := safe.GetText(w, op.Relation); e != nil {
			err = e
		} else if otherNouns, e := safe.GetTextList(w, op.OtherNouns); e != nil {
			err = e
		} else {
			a, b := nouns.Strings(), otherNouns.Strings()
			for _, subject := range a {
				if subject, e := grok.StripArticle(subject); e != nil {
					err = errutil.Append(err, e)
				} else if kind := kind.String(); len(kind) > 0 {
					if e := cat.AssertNounKind(subject, kind); e != nil {
						err = errutil.New(err, e)
					}
				}
				if rel := relation.String(); len(rel) > 0 {
					for _, object := range b {
						if object, e := grok.StripArticle(object); e != nil {
							err = errutil.Append(err, e)
						} else if e := cat.AssertRelative(rel, object, subject); e != nil {
							err = errutil.New(err, e)
						}
					}
				}
			}
		}

		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *DefineOtherRelatives) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefineOtherRelatives) Weave(cat *weave.Catalog) error {
	return cat.Schedule(assert.RequireDeterminers, func(w *weave.Weaver) (err error) {
		if nouns, e := safe.GetTextList(w, op.Nouns); e != nil {
			err = e
		} else if relation, e := safe.GetText(w, op.Relation); e != nil {
			err = e
		} else if otherNouns, e := safe.GetTextList(w, op.OtherNouns); e != nil {
			err = e
		} else {
			a, b := nouns.Strings(), otherNouns.Strings()
			if rel := relation.String(); len(rel) > 0 {
				for _, subject := range a {
					if subject, e := grok.StripArticle(subject); e != nil {
						err = errutil.New(err, e)
					} else {
						for _, object := range b {
							if object, e := grok.StripArticle(object); e != nil {
								err = errutil.New(err, e)
							} else if e := cat.AssertRelative(rel, object, subject); e != nil {
								err = errutil.New(err, e)
							}
						}
					}
				}
			}
		}
		return
	})
}
