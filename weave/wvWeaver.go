package weave

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

type ScheduledCallback func(weaver.Weaves, rt.Runtime) error

type localWeaver struct {
	d *Domain
	*mdl.Pen
}

func (ja localWeaver) GenerateUniqueName(category string) string {
	return ja.d.cat.NewCounter(category)
}

// fix: make the signature or caller transparent
func (ja localWeaver) AddFact(key string, parts ...string) (err error) {
	if ok, e := ja.Pen.AddFact(key, parts...); e != nil {
		err = e
	} else if !ok {
		err = mdl.Duplicate
	}
	return
}
func (ja localWeaver) AddNounTrait(noun, trait string) (err error) {
	return ja.AddNounValue(noun, trait, truly())
}
func (ja localWeaver) AddNounValue(noun, field string, value rt.Assignment) (err error) {
	// if we are adding an initial value for a different domain
	// then that gets changed into "set value" triggered on "begin domain"
	var u mdl.DomainValueError
	if e := ja.Pen.AddNounValue(noun, field, value); !errors.As(e, &u) {
		err = e // nil or unexpected error.
	} else {
		ja.d.initialValues = ja.d.initialValues.add(u.Noun, u.Field, u.Value)
	}
	return
}

func truly() rt.Assignment {
	return &assign.FromBool{
		Value: &literal.BoolValue{Value: true},
	}
}
