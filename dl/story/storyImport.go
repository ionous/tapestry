package story

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
)

// top level imports
type StoryStatement interface {
	ImportPhrase(*Importer) error
}

type nounImporter interface {
	importNouns(*Importer) error
}

// (the) colors are red, blue, or green.
func (op *AspectTraits) ImportPhrase(k *Importer) (err error) {
	var ts []string
	for _, t := range op.TraitPhrase.Trait {
		ts = append(ts, t.String())
	}
	k.Write(&eph.EphAspects{Aspects: op.Aspect.String(), Traits: ts})
	return
}

// horses are usually fast.
func (op *Certainties) ImportPhrase(k *Importer) (err error) {
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

func (op *Comment) ImportPhrase(k *Importer) (err error) {
	// do nothing for now.
	return
}

func (op *GrammarDecl) ImportPhrase(k *Importer) (err error) {
	switch el := op.Grammar.(type) {
	case *grammar.Alias:
		k.Write(&eph.EphAliases{ShortName: el.AsNoun, Aliases: el.Names})
	case *grammar.Directive:
		name := strings.Join(el.Lede, "/")
		k.Write(&eph.EphDirectives{Name: name, Directive: *el})
	}
	return
}

type NounContinuation interface {
	importNounPhrase(*Importer) error
}

func (op *NounKindStatement) ImportPhrase(k *Importer) error {
	return importNounPhrase(k, op.Nouns, &op.KindOfNoun, op.More)
}

func (op *NounTraitStatement) ImportPhrase(k *Importer) error {
	return importNounPhrase(k, op.Nouns, &op.NounTraits, op.More)
}

func (op *NounRelationStatement) ImportPhrase(k *Importer) error {
	return importNounPhrase(k, op.Nouns, &op.NounRelation, op.More)
}

func importNounPhrase(k *Importer, nouns []NamedNoun, first NounContinuation, rest []NounContinuation) (err error) {
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
func (op *NounAssignment) ImportPhrase(k *Importer) (err error) {
	if text, e := ConvertText(k, op.Lines.String()); e != nil {
		err = e
	} else if e := CollectSubjectNouns(k, op.Nouns); e != nil {
		err = e
	} else {
		prop := op.Property.String()
		for _, noun := range k.Env().Recent.Nouns.Subjects {
			k.Write(&eph.EphValues{Noun: noun, Field: prop, Value: T(text)})
		}
	}
	return
}

// ex. On the beach are shells.
func (op *RelativeToNoun) ImportPhrase(k *Importer) (err error) {
	if e := CollectObjectNouns(k, op.Nouns); e != nil {
		err = e
	} else if e := CollectSubjectNouns(k, op.Nouns1); e != nil {
		err = e
	} else {
		relation := op.Relation.String()
		for _, object := range k.Env().Recent.Nouns.Objects {
			for _, subject := range k.Env().Recent.Nouns.Subjects {
				k.Write(&eph.EphRelatives{Rel: relation, Noun: subject, OtherNoun: object})
			}
		}
	}
	return
}
