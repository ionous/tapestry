package weave

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/imp/assert"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/weave/eph"
	"github.com/ionous/errutil"
)

func (ctx *Context) AssertAlias(opShortName string, opAliases ...string) (err error) {
	d, at := ctx.d, ctx.at
	if noun, e := getClosestNoun(d, opShortName); e != nil {
		err = e
	} else {
		for _, a := range opAliases {
			if a, ok := UniformString(a); !ok {
				err = errutil.Append(err, InvalidString(a))
			} else {
				if !noun.AddAlias(a, at) {
					LogWarning(errutil.Fmt("duplicate alias %q for %q at %s",
						a, noun.name, at))
				}
			}
		}
	}
	return
}

// generates traits and adds them to a custom aspect kind.
func (ctx *Context) AssertAspectTraits(opAspects string, opTraits []string) (err error) {
	d, at := ctx.d, ctx.at
	// we dont singularize aspects even thought its a kind;
	// most are really singularizable anyway, and some common things like "darkness" dont singularize correctly.
	if aspect, ok := UniformString(opAspects); !ok {
		err = InvalidString(opAspects)
	} else if traits, e := UniformStrings(opTraits); e != nil {
		err = e
	} else {
		kid := d.EnsureKind(aspect, at)
		kid.AddRequirement(kindsOf.Aspect.String())
		if len(traits) > 0 {
			err = d.QueueEphemera(at, eph.PhaseFunction{assert.AspectPhase,
				func(assert.World, assert.Assertions) (err error) {
					var conflict *Conflict // checks for conflicts, allows duplicates.
					if e := kid.AddField(&traitDef{at, aspect, traits}); errors.As(e, &conflict) && conflict.Reason == Duplicated {
						LogWarning(e) // warn if it was a duplicated definition
					} else {
						err = e // some other error ( or nil )
					}
					return
				}})
		}
	}
	return
}

func (ctx *Context) AssertCheck(opName string, opExe []rt.Execute, opExpect literal.LiteralValue) (err error) {
	c, at := ctx.c, ctx.at
	// fix. todo: this isnt very well thought out right now --
	// what if a check is part of a story scene? shouldnt it have access to those objects?
	// if checks always establish their own domain, why do they have a duplicate name?
	// there are some checks that have their own scenes, some that have expressions to run, some that have things to check.
	// and others that have just one of those things. ( are there expectations that can simply verify the existence of objects in a model? )
	// etc.
	if name, ok := UniformString(opName); !ok {
		err = InvalidString(opName)
	} else if d, e := c.EnsureDomain(name, at); e != nil {
		err = e
	} else {
		// uses directive phase just to be near the end somewhere...
		err = d.QueueEphemera(at, eph.PhaseFunction{assert.DirectivePhase,
			func(assert.World, assert.Assertions) (err error) {
				check := d.EnsureCheck(name, at)
				if e := check.setExpectation(opExpect); e != nil {
					err = e
				} else if e := check.setProg(opExe); e != nil {
					err = e
				}
				return
			}})
	}
	return
}

func (ctx *Context) AssertDefinition(path ...string) (err error) {
	d, at := ctx.d, ctx.at
	if end := len(path) - 1; end <= 0 {
		err = errutil.New("path too short", path)
	} else {
		path, value := path[:end], path[end]
		err = d.AddDefinition(MakeKey(path...), at, value)
	}
	return

}

func (ctx *Context) Schedule(when assert.Phase, do func(assert.World, assert.Assertions) error) (err error) {
	c, d, at := ctx.c, ctx.d, ctx.at
	var writeEph eph.WriterFun = func(op eph.Ephemera) {
		d.QueueEphemera(at, op)
	}
	return do(c, eph.NewCommandBuilder(writeEph))
}

func getClosestNoun(d *Domain, rawName string) (ret *ScopedNoun, err error) {
	if short, ok := UniformString(rawName); !ok {
		err = InvalidString(rawName)
	} else if noun, ok := d.GetClosestNoun(short); !ok {
		err = errutil.New("unknown noun", rawName)
	} else {
		ret = noun
	}
	return
}
