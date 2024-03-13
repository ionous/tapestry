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
	AddFields(kind string, fields []mdl.FieldInfo) error
	AddGrammar(name string, prog *grammar.Directive) error
	AddKind(kind, ancestor string) error
	AddKindTrait(kind, trait string) error
	AddPlural(many, one string) error
	AddNounKind(noun, kind string) error
	AddNounName(noun, name string, rank int) error
	AddNounTrait(noun, trait string) error
	AddNounValue(noun, prop string, val rt.Assignment) error
	AddNounPath(noun string, path []string, val literal.LiteralValue) error
	AddNounPair(rel, noun, otherNoun string) error
	AddTraits(aspect string, traits []string) error
	AddFact(key string, partsAndValue ...string) error
	//
	GetOpposite(string) (string, error)
	GetPlural(string) string
	GetRelativeNouns(noun, relation string, primary bool) ([]string, error)
	GetSingular(string) string
	GetUniqueName(category string) string
	// apply the passed macro to the passed nouns
	Apply(verb Macro, lhs, rhs []string) error
}

type Context struct {
	Registrar
	proc ProcessingList
}

func (ctx *Context) PostProcess(i Priority, p Process) (_ error) {
	ctx.proc.AddToList(i, p)
	return
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
