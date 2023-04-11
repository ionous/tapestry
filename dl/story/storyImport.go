package story

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/imp/assert"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"

	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// Execute - called by the macro runtime during weave.
func (op *DefineTraits) Execute(macro rt.Runtime) error {
	return imp.StoryStatement(macro, op)
}

// (the) colors are red, blue, or green.
func (op *DefineTraits) PostImport(k *imp.Importer) (err error) {
	if traits, e := safe.GetTextList(k, op.Traits); e != nil {
		err = e
	} else if aspect, e := safe.GetText(k, op.Aspect); e != nil {
		err = e
	} else {
		err = k.AssertAspectTraits(aspect.String(), traits.Strings())
	}
	return
}

// Execute - called by the macro runtime during weave.
func (op *GrammarDecl) Execute(macro rt.Runtime) error {
	return imp.StoryStatement(macro, op)
}

func (op *GrammarDecl) PostImport(k *imp.Importer) (err error) {
	switch el := op.Grammar.(type) {
	case *grammar.Alias:
		err = k.AssertAlias(el.AsNoun, el.Names...)
	case *grammar.Directive:
		name := strings.Join(el.Lede, "/")
		err = k.AssertGrammar(name, el)
	}
	return
}

// Execute - called by the macro runtime during weave.
func (op *DefineNounTraits) Execute(macro rt.Runtime) error {
	return imp.StoryStatement(macro, op)
}

func (op *DefineNounTraits) PostImport(k *imp.Importer) (err error) {
	if nouns, e := safe.GetTextList(k, op.Nouns); e != nil {
		err = e
	} else if kind, e := safe.GetOptionalText(k, op.Kind, ""); e != nil {
		err = e
	} else if traits, e := safe.GetTextList(k, op.Traits); e != nil {
		err = e
	} else if bareNames, e := ImportNounProperties(k, nouns.Strings()); e != nil {
		err = e
	} else {
		if kind := kind.String(); len(kind) > 0 {
			for _, n := range bareNames {
				if e := k.AssertNounKind(n, kind); e != nil {
					err = errutil.Append(err, e)
				}
			}
		}
		if traits := traits.Strings(); len(traits) > 0 {
			for _, t := range traits {
				for _, n := range bareNames {
					if e := assert.AssertNounValue(k, B(true), n, t); e != nil {
						err = errutil.Append(err, e)
						break // out of the traits to the next noun
					}
				}
			}
		}

	}
	return
}

// Execute - called by the macro runtime during weave.
func (op *DefineNouns) Execute(macro rt.Runtime) error {
	return imp.StoryStatement(macro, op)
}

func (op *DefineNouns) PostImport(k *imp.Importer) (err error) {
	if nouns, e := safe.GetTextList(k, op.Nouns); e != nil {
		err = e
	} else if kind, e := safe.GetText(k, op.Kind); e != nil {
		err = e
	} else if bareNames, e := ImportNounProperties(k, nouns.Strings()); e != nil {
		err = e
	} else {
		if kind := kind.String(); len(kind) > 0 {
			for _, n := range bareNames {
				if e := k.AssertNounKind(n, kind); e != nil {
					err = errutil.Append(err, e)
				}
			}
		}
	}
	return
}

// Execute - called by the macro runtime during weave.
func (op *NounAssignment) Execute(macro rt.Runtime) error {
	return imp.StoryStatement(macro, op)
}

// ex. The description of the nets is xxx
func (op *NounAssignment) PostImport(k *imp.Importer) (err error) {
	if nouns, e := safe.GetTextList(k, op.Nouns); e != nil {
		err = e
	} else if field, e := safe.GetText(k, op.FieldName); e != nil {
		err = e
	} else if lines, e := ConvertText(k, op.Lines.String()); e != nil {
		err = e
	} else if subjects, e := ReadNouns(k, nouns.Strings()); e != nil {
		err = e
	} else {
		field, lines := field.String(), T(lines)
		for _, noun := range subjects {
			if e := assert.AssertNounValue(k, lines, noun, field); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	return
}

// Execute - called by the macro runtime during weave.
func (op *DefineRelatives) Execute(macro rt.Runtime) error {
	return imp.StoryStatement(macro, op)
}

func (op *DefineRelatives) PostImport(k *imp.Importer) (err error) {
	if nouns, e := safe.GetTextList(k, op.Nouns); e != nil {
		err = e
	} else if kind, e := safe.GetOptionalText(k, op.Kind, ""); e != nil {
		err = e
	} else if relation, e := safe.GetText(k, op.Relation); e != nil {
		err = e
	} else if otherNouns, e := safe.GetTextList(k, op.OtherNouns); e != nil {
		err = e
	} else if a, e := ReadNouns(k, nouns.Strings()); e != nil {
		err = e
	} else if b, e := ReadNouns(k, otherNouns.Strings()); e != nil {
		err = e
	} else {
		for _, subject := range a {
			if kind := kind.String(); len(kind) > 0 {
				if e := k.AssertNounKind(subject, kind); e != nil {
					err = errutil.New(err, e)
				}
			}
			if rel := relation.String(); len(rel) > 0 {
				for _, object := range b {
					if e := k.AssertRelative(rel, object, subject); e != nil {
						err = errutil.New(err, e)
					}
				}
			}
		}
	}
	return
}

// Execute - called by the macro runtime during weave.
func (op *DefineOtherRelatives) Execute(macro rt.Runtime) error {
	return imp.StoryStatement(macro, op)
}

func (op *DefineOtherRelatives) PostImport(k *imp.Importer) (err error) {
	if nouns, e := safe.GetTextList(k, op.Nouns); e != nil {
		err = e
	} else if relation, e := safe.GetText(k, op.Relation); e != nil {
		err = e
	} else if otherNouns, e := safe.GetTextList(k, op.OtherNouns); e != nil {
		err = e
	} else if a, e := ReadNouns(k, nouns.Strings()); e != nil {
		err = e
	} else if b, e := ReadNouns(k, otherNouns.Strings()); e != nil {
		err = e
	} else {
		if rel := relation.String(); len(rel) > 0 {
			for _, subject := range a {
				for _, object := range b {
					if e := k.AssertRelative(rel, object, subject); e != nil {
						err = errutil.New(err, e)
					}
				}
			}
		}
	}
	return
}
