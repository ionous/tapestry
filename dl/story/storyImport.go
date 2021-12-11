package story

import (
	"git.sr.ht/~ionous/iffy/dl/grammar"
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
	if aspect, e := NewAspect(k, op.Aspect); e != nil {
		err = e
	} else {
		err = op.TraitPhrase.ImportTraits(k, aspect)
	}
	return
}

// horses are usually fast.
func (op *Certainties) ImportPhrase(k *Importer) (err error) {
	if certainty, e := op.Certainty.ImportString(k); e != nil {
		err = e
	} else if trait, e := NewTrait(k, op.Trait); e != nil {
		err = e
	} else if kind, e := NewPluralKinds(k, op.PluralKinds); e != nil {
		err = e
	} else {
		k.NewCertainty(certainty, trait, kind)
	}
	return
}

func (op *Comment) ImportPhrase(k *Importer) (err error) {
	// do nothing for now.
	return
}

func (op *GrammarDecl) ImportPhrase(k *Importer) error {
	_, e := k.NewProg("grammar", &grammar.Grammar{op.Grammar})
	return e
}

// ex. The description of the nets is xxx
func (op *NounAssignment) ImportPhrase(k *Importer) (err error) {
	if prop, e := NewProperty(k, op.Property); e != nil {
		err = e
	} else if text, e := ConvertText(k, op.Lines.String()); e != nil {
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
		for _, noun := range k.Env().Recent.Nouns.Subjects {
			k.NewValue(noun, prop, text)
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
	if relation, e := NewRelation(k, op.Relation); e != nil {
		err = e
	} else if e := k.Env().Recent.Nouns.CollectObjects(func() error {
		return ImportNamedNouns(k, op.Nouns)
	}); e != nil {
		err = e
	} else if e := k.Env().Recent.Nouns.CollectSubjects(func() error {
		return ImportNamedNouns(k, op.Nouns1)
	}); e != nil {
		err = e
	} else {
		domain := k.Env().Current.Domain
		for _, object := range k.Env().Recent.Nouns.Objects {
			for _, subject := range k.Env().Recent.Nouns.Subjects {
				k.NewRelative(subject, relation, object, domain)
			}
		}
	}
	return
}
