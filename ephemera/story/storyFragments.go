package story

import (
	"git.sr.ht/~ionous/iffy/ephemera/eph"
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

func (op *Lede) ImportNouns(k *Importer) (err error) {
	if e := k.Env().Recent.Nouns.CollectSubjects(func() (err error) {
		for _, nn := range op.Nouns {
			if e := nn.ImportNouns(k); e != nil {
				err = errutil.Append(err, e)
			}
		}
		return
	}); e != nil {
		err = e
	} else {
		err = op.NounPhrase.ImportNouns(k)
	}
	return
}

func (op *Summary) ImportNouns(k *Importer) (err error) {
	if text, e := ConvertText(op.Lines.String()); e != nil {
		err = e
	} else {
		// give "things" an "description"
		if once := "summary"; k.Once(once) {
			domain := k.Env().Game.Domain
			things := k.NewDomainName(domain, "things", tables.NAMED_KINDS, once)
			appear := k.NewDomainName(domain, "description", tables.NAMED_FIELD, once)
			k.NewField(things, appear, tables.PRIM_TEXT, "")
		}
		prop := k.NewName("description", tables.NAMED_FIELD, op.At.String())
		noun := LastNameOf(k.Env().Recent.Nouns.Subjects)
		k.NewValue(noun, prop, text)
	}
	return
}

func LastNameOf(n []eph.Named) (ret eph.Named) {
	if cnt := len(n); cnt > 0 {
		ret = (n)[cnt-1]
	}
	return
}

func (op *Tail) ImportNouns(k *Importer) (err error) {
	if e := op.Pronoun.ImportNouns(k); e != nil {
		err = e
	} else if e := op.NounPhrase.ImportNouns(k); e != nil {
		err = e
	}
	return
}

func (op *TraitPhrase) ImportTraits(k *Importer, aspect eph.Named) (err error) {
	for rank, trait := range op.Trait {
		if t, e := NewTrait(k, trait); e != nil {
			err = errutil.Append(err, e)
		} else {
			k.NewTrait(t, aspect, rank)
		}
	}
	return
}
