package story

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/imp"
)

// top level imports
type StoryStatement interface {
	PostImport(*imp.Importer) error
}

type nounImporter interface {
	importNouns(*imp.Importer) error
}

// (the) colors are red, blue, or green.
func (op *AspectTraits) PostImport(k *imp.Importer) (err error) {
	var ts []string
	for _, t := range op.TraitPhrase.Trait {
		ts = append(ts, t.String())
	}
	k.WriteEphemera(&eph.EphAspects{Aspects: op.Aspect.String(), Traits: ts})
	return
}

// horses are usually fast.
func (op *Certainties) PostImport(k *imp.Importer) (err error) {
	certaintiesNotImplemented.PrintOnce()
	// if certainty, e := op.Certainty.ImportString(k); e != nil {
	// 	err = e
	// } else if trait, e := NewTrait(k, op.Trait); e != nil {
	// 	err = e
	// } else if kind, e := NewPluralKinds(k, op.PluralKinds); e != nil {
	// 	err = e
	// } else {
	// 	k.NewCertainty(certainty, trait, kind)
	// }
	return
}

var certaintiesNotImplemented eph.PrintOnce = "certainties not implemented"

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

type NounContinuation interface {
	importNounPhrase(*imp.Importer) error
}

func (op *NounKindStatement) PostImport(k *imp.Importer) error {
	var yuck []NamedNoun // cast the elements to their less specific type
	for _, n := range op.Nouns {
		yuck = append(yuck, n)
	}
	return importNounPhrase(k, yuck, &op.KindOfNoun, op.More)
}

func (op *NounTraitStatement) PostImport(k *imp.Importer) error {
	return importNounPhrase(k, op.Nouns, &op.NounTraits, op.More)
}

func (op *NounRelationStatement) PostImport(k *imp.Importer) error {
	return importNounPhrase(k, op.Nouns, &op.NounRelation, op.More)
}

func importNounPhrase(k *imp.Importer, nouns []NamedNoun, first NounContinuation, rest []NounContinuation) (err error) {
	if e := CollectSubjectNouns(k, nouns); e != nil {
		err = e
	} else if e := first.importNounPhrase(k); e != nil {
		err = e
	} else {
		for _, el := range rest {
			if e := el.importNounPhrase(k); e != nil {
				err = e
				break
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

// ex. On the beach are shells.
func (op *RelativeToNoun) PostImport(k *imp.Importer) (err error) {
	if e := CollectObjectNouns(k, op.Nouns); e != nil {
		err = e
	} else if e := CollectSubjectNouns(k, op.OtherNouns); e != nil {
		err = e
	} else {
		relation := op.Relation.String()
		for _, object := range k.Env().Recent.Nouns.Objects {
			for _, subject := range k.Env().Recent.Nouns.Subjects {
				k.WriteEphemera(&eph.EphRelatives{Rel: relation, Noun: subject, OtherNoun: object})
			}
		}
	}
	return
}
