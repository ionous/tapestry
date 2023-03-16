package story

import (
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/rt"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/eph"
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
		k.WriteEphemera(&eph.EphAspects{Aspects: aspect.String(), Traits: traits.Strings()})
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
		k.WriteEphemera(&eph.EphAliases{ShortName: el.AsNoun, Aliases: el.Names})
	case *grammar.Directive:
		name := strings.Join(el.Lede, "/")
		k.WriteEphemera(&eph.EphDirectives{Name: name, Directive: *el})
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
	} else if bareNames, e := ImportNouns(k, nouns.Strings()); e != nil {
		err = e
	} else {
		if kind := kind.String(); len(kind) > 0 {
			for _, n := range bareNames {
				k.WriteEphemera(&eph.EphNouns{Noun: n, Kind: kind})
			}
		}
		if traits := traits.Strings(); len(traits) > 0 {
			for _, t := range traits {
				for _, n := range bareNames {
					k.WriteEphemera(&eph.EphValues{Noun: n, Field: t, Value: B(true)})
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
	} else if bareNames, e := ImportNouns(k, nouns.Strings()); e != nil {
		err = e
	} else {
		if kind := kind.String(); len(kind) > 0 {
			for _, n := range bareNames {
				k.WriteEphemera(&eph.EphNouns{Noun: n, Kind: kind})
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
	} else if subjects, e := ImportNouns(k, nouns.Strings()); e != nil {
		err = e
	} else {
		field, lines := field.String(), T(lines)
		for _, noun := range subjects {
			k.WriteEphemera(&eph.EphValues{Noun: noun, Field: field, Value: lines})
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
	} else if a, e := ImportNouns(k, nouns.Strings()); e != nil {
		err = e
	} else if b, e := ImportNouns(k, otherNouns.Strings()); e != nil {
		err = e
	} else {
		for _, subject := range a {
			if kind := kind.String(); len(kind) > 0 {
				k.WriteEphemera(&eph.EphNouns{Noun: subject, Kind: kind})
			}
			if rel := relation.String(); len(rel) > 0 {
				for _, object := range b {
					k.WriteEphemera(&eph.EphRelatives{Rel: rel, Noun: object, OtherNoun: subject})
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
	} else if a, e := ImportNouns(k, nouns.Strings()); e != nil {
		err = e
	} else if b, e := ImportNouns(k, otherNouns.Strings()); e != nil {
		err = e
	} else {
		if rel := relation.String(); len(rel) > 0 {
			for _, subject := range a {
				for _, object := range b {
					k.WriteEphemera(&eph.EphRelatives{Rel: rel, Noun: object, OtherNoun: subject})
				}
			}
		}
	}
	return
}
