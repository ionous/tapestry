package story

import (
	"errors"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"

	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/literal"
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
		err = errutil.New("unknown grammar %T", el)
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
		} else if bareNames, e := readNounsWithProperties(w, nouns.Strings()); e != nil {
			err = e
		} else {
			if kind := kind.String(); len(kind) > 0 {
				for _, n := range bareNames {
					if e := cat.AssertNounKind(n, kind); e != nil {
						err = errutil.Append(err, e)
					}
				}
			}
			if traits := traits.Strings(); len(traits) > 0 {
				for _, t := range traits {
					for _, n := range bareNames {
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
		} else if macro, ok := weave.UniformString(macro.String()); !ok {
			err = errutil.New("missing macro name")
		} else if rev, e := safe.GetOptionalBool(w, op.Reversed, false); e != nil {
			err = e
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
		} else if res, e := w.Grok(text.String()); e != nil {
			err = e
		} else if src, e := genNouns(w, res.Sources, res.Macro.Type == grok.Macro_ManyTargets); e != nil {
			err = e
		} else if tgt, e := genNouns(w, res.Targets, res.Macro.Type != grok.Macro_ManyTargets &&
			res.Macro.Type != grok.Macro_ManyMany); e != nil {
			err = e
		} else {
			macro := res.Macro.Name
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
							err = errutil.New("Declare statement: %s", msg)
						}
					}
				}
			}
		}
		return
	})
}

// add nouns and values
func genNouns(w *weave.Weaver, ns []grok.Noun, wantOne bool) (ret g.Value, err error) {
	cat, domain := w.Catalog, w.Domain
	if cnt := len(ns); wantOne && cnt != 1 {
		err = errutil.New("expected exactly one noun")
	} else {
		names := make([]string, cnt)
	Out:
		for i := 0; i < cnt; i++ {
			n := ns[i]
			name := n.Name.String()
			if name == "you" {
				// tdb: the current thought is that "the player" should be a variable;
				// currently its an "agent".
				name = "self"
			} else if !n.Exact {
				if fold, e := domain.GetClosestNoun(name); e == nil {
					name = fold
				} else if !errors.Is(e, mdl.Missing) {
					err = e
					break Out
				}

			}
			names[i] = name
			// tbd: might be some nicer, faster ways of handling all this
			// instead of passing it through the highest level catalog interfaces
			for _, k := range n.Kinds {
				if e := cat.AssertNounKind(name, k.String()); e != nil {
					err = e
					break Out
				}
			}
			for _, t := range n.Traits {
				if e := cat.AssertNounValue(name, t.String(), nil, literal.B(true)); e != nil {
					err = e
					break Out
				}
			}
		} // end for loop
		// all done?
		if err == nil {
			if wantOne {
				ret = g.StringOf(names[0])
			} else {
				ret = g.StringsOf(names)
			}
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
		} else if bareNames, e := readNounsWithProperties(w, nouns.Strings()); e != nil {
			err = e
		} else {
			if kind := kind.String(); len(kind) > 0 {
				for _, n := range bareNames {
					if e := cat.AssertNounKind(n, kind); e != nil {
						err = errutil.Append(err, e)
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
		} else if subjects, e := readNouns(w, nouns.Strings()); e != nil {
			err = e
		} else {
			field, lines := field.String(), T(lines)
			for _, noun := range subjects {
				if e := assertNounValue(cat, lines, noun, field); e != nil {
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
		} else if a, e := readNouns(w, nouns.Strings()); e != nil {
			err = e
		} else if b, e := readNouns(w, otherNouns.Strings()); e != nil {
			err = e
		} else {
			for _, subject := range a {
				if kind := kind.String(); len(kind) > 0 {
					if e := cat.AssertNounKind(subject, kind); e != nil {
						err = errutil.New(err, e)
					}
				}
				if rel := relation.String(); len(rel) > 0 {
					for _, object := range b {
						if e := cat.AssertRelative(rel, object, subject); e != nil {
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
		} else if a, e := readNouns(w, nouns.Strings()); e != nil {
			err = e
		} else if b, e := readNouns(w, otherNouns.Strings()); e != nil {
			err = e
		} else {
			if rel := relation.String(); len(rel) > 0 {
				for _, subject := range a {
					for _, object := range b {
						if e := cat.AssertRelative(rel, object, subject); e != nil {
							err = errutil.New(err, e)
						}
					}
				}
			}
		}
		return
	})
}
