package jess

import (
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

// generates story fragments for weaving expects normalized names.
// matches "pen" so it can write fragments straight to the db;
// alternatively, an implementation could generate story commands if that were useful for some reason.
type Registrar interface {
	AddPlural(many, one string) error
	AddGrammar(name string, prog *grammar.Directive) error
	AddKind(kind, ancestor string) error
	AddKindFields(kind string, fields []mdl.FieldInfo) error
	// apply a default trait to a kind
	AddKindTrait(kind, trait string) error
	// set the possible traits of an aspect
	AddAspectTraits(aspect string, traits []string) error
	AddNounKind(noun, kind string) error
	AddNounName(noun, name string, rank int) error
	AddNounTrait(noun, trait string) error
	AddNounValue(noun, prop string, val rt.Assignment) error
	AddNounPath(noun string, path []string, val literal.LiteralValue) error
	AddNounPair(rel, noun, otherNoun string) error
	AddFact(key string, partsAndValue ...string) error
	// generates a unique name with a passed category
	GenerateUniqueName(category string) string
	// return a list of related nouns
	GetRelativeNouns(noun, relation string, primary bool) ([]string, error)
	// given a plural name, return its singular version.
	GetPlural(string) string
	// given a singular name, return its plural.
	GetSingular(string) string
	// give a word, return its pre-defined opposite;
	// most often used with directions to know if we are coming or going.
	GetOpposite(string) (string, error)
}

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
