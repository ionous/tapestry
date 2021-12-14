package story

import (
	"git.sr.ht/~ionous/iffy/dl/eph"
	"git.sr.ht/~ionous/iffy/dl/literal"
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
	if text, e := ConvertText(k, op.Lines.String()); e != nil {
		err = e
	} else {
		// give "things" an "description"
		if once := "summary"; k.Once(once) {
			k.WriteOnce(&eph.EphKinds{Kinds: "things", Contain: []eph.EphParams{{Name: "description", Affinity: eph.Affinity{eph.Affinity_Text}}}})
		}
		noun := LastNameOf(k.Env().Recent.Nouns.Subjects)
		k.Write(&eph.EphValues{Noun: noun, Field: "description", Value: &literal.TextValue{text}})
	}
	return
}

func LastNameOf(n []string) (ret string) {
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
