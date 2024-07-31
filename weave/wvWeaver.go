package weave

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

type penCreator struct {
	d *Domain
}

func (p penCreator) NewPen(when weaver.Phase, pos compact.Source) *mdl.Pen {
	return p.d.cat.Modeler.PinPos(p.d.name, pos)
}

type localWeaver struct {
	d *Domain
	*mdl.Pen
}

func (ja localWeaver) GenerateUniqueName(category string) string {
	// FIX: this should probably be part of pen.
	return ja.d.cat.NewCounter(category)
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
		ja.d.AddInitialValue(u.Noun, u.Field, u.Value)
	}
	return
}

func truly() rt.Assignment {
	return &call.FromBool{
		Value: &literal.BoolValue{Value: true},
	}
}
