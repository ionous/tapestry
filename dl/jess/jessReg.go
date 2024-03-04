package jess

import (
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

// generates story fragments for weaving expects normalized names.
// matches "pen" so it can write fragments straight to the db;
// alternatively, an implementation could generate story commands if that were useful for some reason.
type Registrar interface {
	AddFields(kind string, fields []mdl.FieldInfo) error
	AddGrammar(name string, prog *grammar.Directive) error
	AddKind(kind, ancestor string) error
	AddKindTrait(kind, trait string) error
	AddPlural(many, one string) error
	AddNounKind(noun, kind string) error
	AddNounName(noun, name string, rank int) error
	AddNounTrait(noun, trait string) error
	AddNounValue(noun, prop string, val rt.Assignment) error // tbd: or would g.Value be better at this point?
	AddTraits(aspect string, traits []string) error
	//
	GetClosestNoun(name string) (string, error)
	GetPlural(string) string
	GetSingular(string) string
	GetUniqueName(category string) string
	// apply the passed macro to the passed nouns
	Apply(verb Macro, lhs, rhs []string) error
	// register a function for later processing
	PostProcess(Priority, Process) error
}

type Priority int
type Process func(Query, Registrar) error

const (
	NounSettings Priority = iota
	// NounPairings
	Understandings
	PriorityCount
)

// setup the default traits for the passed kind
func AddKindTraits(rar Registrar, kind string, traits Traitor) (err error) {
	for ts := traits; ts.HasNext(); {
		t := ts.GetNext()
		str := t.String()
		if e := rar.AddKindTrait(kind, inflect.Normalize(str)); e != nil {
			err = e
			break
		}
	}
	return
}
