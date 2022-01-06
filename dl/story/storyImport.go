package story

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"github.com/ionous/errutil"
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
	k.WriteEphemera(&eph.EphAspects{Aspects: op.Aspect.String(), Traits: ts})
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
		k.WriteEphemera(&eph.EphAliases{ShortName: el.AsNoun, Aliases: el.Names})
	case *grammar.Directive:
		name := strings.Join(el.Lede, "/")
		k.WriteEphemera(&eph.EphDirectives{Name: name, Directive: *el})
	}
	return
}

// ex. The description of the nets is xxx
func (op *NounAssignment) ImportPhrase(k *Importer) (err error) {
	if text, e := ConvertText(k, op.Lines.String()); e != nil {
		err = e
	} else if e := k.Env().Recent.Nouns.CollectSubjects(func() (err error) {
		for _, n := range op.Nouns {
			if e := n.ImportNouns(k); e != nil {
				err = errutil.Append(err, e)
			}
		}
		return
	}); e != nil {
		err = e
	} else {
		prop := op.Property.String()
		for _, noun := range k.Env().Recent.Nouns.Subjects {
			k.WriteEphemera(&eph.EphValues{Noun: noun, Field: prop, Value: T(text)})
		}
	}
	return
}

func (op *NounStatement) ImportPhrase(k *Importer) (err error) {
	if e := op.Lede.ImportNouns(k); e != nil {
		err = e
	} else {
		if els := op.Tail; els != nil {
			for _, el := range els {
				if e := el.ImportNouns(k); e != nil {
					err = errutil.Append(err, e)
				}
			}
		}
		if err == nil && op.Summary != nil {
			err = op.Summary.ImportNouns(k)
		}
	}
	return
}

// ex. On the beach are shells.
func (op *RelativeToNoun) ImportPhrase(k *Importer) (err error) {
	if e := k.Env().Recent.Nouns.CollectObjects(func() error {
		return ImportNamedNouns(k, op.Nouns)
	}); e != nil {
		err = e
	} else if e := k.Env().Recent.Nouns.CollectSubjects(func() error {
		return ImportNamedNouns(k, op.Nouns1)
	}); e != nil {
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
