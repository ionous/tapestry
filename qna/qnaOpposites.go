package qna

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

type kinds interface {
	GetKindByName(name string) (*rt.Kind, error)
}

// for a given aspect and one of its traits, find the opposite trait
// returns error if there is no such aspect, trait, or opposite value.
func oppositeDay(ks kinds, aspect, trait string, b bool) (ret string, err error) {
	if b {
		ret = trait
	} else if k, e := ks.GetKindByName(aspect); e != nil {
		err = e
	} else if cnt := k.FieldCount(); cnt != 2 {
		err = errutil.Fmt("couldn't determine the opposite of %s.%s", aspect, trait)
	} else if i := k.FieldIndex(trait); i < 0 {
		err = errutil.Fmt("couldn't find the trait %s.%s", aspect, trait)
	} else {
		field := k.Field((i + 1) & 1)
		ret = field.Name
	}
	return
}
