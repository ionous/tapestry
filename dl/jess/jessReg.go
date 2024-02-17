package jess

import "git.sr.ht/~ionous/tapestry/support/inflect"

// expects normalized names
type Registrar interface {
	AddKind(kind, ancestor string) error
	AddKindTrait(kind, trait string) error
}

func AddTraitsToKind(rar Registrar, kind string, traits Traitor) (err error) {
	for ts := traits; ts.HasNext(); {
		t := ts.GetNext()
		str := t.Matched.String()
		if e := rar.AddKindTrait(kind, inflect.Normalize(str)); e != nil {
			err = e
			break
		}
	}
	return
}
