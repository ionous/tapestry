package story

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// top level imports
type StoryStatement interface {
	PostImport(*imp.Importer) error
}

type nounImporter interface {
	importNouns(*imp.Importer) error
}

// (the) colors are red, blue, or green.
func (op *DefineTraits) PostImport(k *imp.Importer) (err error) {
	if traits, e := safe.GetTextList(nil, op.Traits); e != nil {
		err = e
	} else if aspect, e := safe.GetText(nil, op.Aspect); e != nil {
		err = e
	} else {
		k.WriteEphemera(&eph.EphAspects{Aspects: aspect.String(), Traits: traits.Strings()})
	}
	return
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

func (op *DefineNounTraits) PostImport(k *imp.Importer) (err error) {
	if nouns, e := safe.GetTextList(nil, op.Nouns); e != nil {
		err = e
	} else if kind, e := safe.GetOptionalText(nil, op.Kind, ""); e != nil {
		err = e
	} else if traits, e := safe.GetTextList(nil, op.Traits); e != nil {
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

func (op *DefineNouns) PostImport(k *imp.Importer) (err error) {
	if nouns, e := safe.GetTextList(nil, op.Nouns); e != nil {
		err = e
	} else if kind, e := safe.GetText(nil, op.Kind); e != nil {
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

// ex. The description of the nets is xxx
func (op *NounAssignment) PostImport(k *imp.Importer) (err error) {
	if text, e := ConvertText(k, op.Lines.String()); e != nil {
		err = e
	} else if e := CollectSubjectNouns(k, op.Nouns); e != nil {
		err = e
	} else {
		prop := op.Property.String()
		for _, noun := range k.Env().Recent.Nouns.Subjects {
			k.WriteEphemera(&eph.EphValues{Noun: noun, Field: prop, Value: T(text)})
		}
	}
	return
}

func (op *DefineRelatives) PostImport(k *imp.Importer) (err error) {
	if nouns, e := safe.GetTextList(nil, op.Nouns); e != nil {
		err = e
	} else if kind, e := safe.GetOptionalText(nil, op.Kind, ""); e != nil {
		err = e
	} else if relation, e := safe.GetText(nil, op.Relation); e != nil {
		err = e
	} else if otherNouns, e := safe.GetTextList(nil, op.OtherNouns); e != nil {
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

func (op *DefineOtherRelatives) PostImport(k *imp.Importer) (err error) {
	if nouns, e := safe.GetTextList(nil, op.Nouns); e != nil {
		err = e
	} else if relation, e := safe.GetText(nil, op.Relation); e != nil {
		err = e
	} else if otherNouns, e := safe.GetTextList(nil, op.OtherNouns); e != nil {
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
