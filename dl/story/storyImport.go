package story

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/rt"
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
		} else if bareNames, e := ImportNounProperties(cat, nouns.Strings()); e != nil {
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
						if e := assert.AssertNounValue(cat, B(true), n, t); e != nil {
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
func (op *DefineNouns) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefineNouns) Weave(cat *weave.Catalog) error {
	return cat.Schedule(assert.RequireDeterminers, func(w *weave.Weaver) (err error) {
		if nouns, e := safe.GetTextList(w, op.Nouns); e != nil {
			err = e
		} else if kind, e := safe.GetText(w, op.Kind); e != nil {
			err = e
		} else if bareNames, e := ImportNounProperties(cat, nouns.Strings()); e != nil {
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
		} else if subjects, e := ReadNouns(cat, nouns.Strings()); e != nil {
			err = e
		} else {
			field, lines := field.String(), T(lines)
			for _, noun := range subjects {
				if e := assert.AssertNounValue(cat, lines, noun, field); e != nil {
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
		} else if a, e := ReadNouns(cat, nouns.Strings()); e != nil {
			err = e
		} else if b, e := ReadNouns(cat, otherNouns.Strings()); e != nil {
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
		} else if a, e := ReadNouns(cat, nouns.Strings()); e != nil {
			err = e
		} else if b, e := ReadNouns(cat, otherNouns.Strings()); e != nil {
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
