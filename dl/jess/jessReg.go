package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

// generates story fragments for weaving expects normalized names.
// matches "pen" so it can write fragments straight to the db;
// alternatively, an implementation could generate story commands if that were useful for some reason.
type Registrar interface {
	AddKind(kind, ancestor string) error
	AddKindTrait(kind, trait string) error
	//
	AddNoun(short, long, kind string) error
	AddNounAlias(noun, name string, rank int) error
	// tbd: or would g.Value be better at this point?
	// do we expect jess to be able to generate evaluations?
	AddNounTrait(noun, trait string) error
	AddNounValue(noun, prop string, val rt.Assignment) error
	//
	Apply(verb Macro, lhs, rhs []string) error

	// fix: move to "Query"?
	GetClosestNoun(name string) (string, error)
	GetExactNoun(name string) (string, error)

	GetPlural(string) string
	GetSingular(string) string
	GetUniqueName(category string) string
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
