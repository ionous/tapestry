package story

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"

	"git.sr.ht/~ionous/tapestry/dl/grammar"
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
