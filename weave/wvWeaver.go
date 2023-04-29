package weave

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

func (cat *Catalog) AssertAlias(opShortName string, opAliases ...string) error {
	return cat.Schedule(assert.AliasPhase, func(ctx *Weaver) (err error) {
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
	})
}

// generates traits and adds them to a custom aspect kind.
func (cat *Catalog) AssertAspectTraits(opAspects string, opTraits []string) error {
	// uses the ancestry phase because it generates kinds ( one per aspect. )
	return cat.Schedule(assert.AncestryPhase, func(ctx *Weaver) (err error) {
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
				err = d.schedule(at, assert.FieldPhase, func(ctx *Weaver) (err error) {
					var conflict *Conflict // checks for conflicts, allows duplicates.
					if e := kid.AddField(&traitDef{at, aspect, traits}); errors.As(e, &conflict) && conflict.Reason == Duplicated {
						LogWarning(e) // warn if it was a duplicated definition
					} else {
						err = e // some other error ( or nil )
					}
					return
				})
			}
		}
		return
	})
}

func (cat *Catalog) AssertCheck(opName string, opExe []rt.Execute, opExpect literal.LiteralValue) error {
	// uses domain phase, because it needs to ensure a domain exists
	return cat.Schedule(assert.PostDomain, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		if name, ok := UniformString(opName); !ok {
			err = InvalidString(opName)
		} else {
			check := d.EnsureCheck(name, at)
			if e := check.setExpectation(opExpect); e != nil {
				err = e
			} else if e := check.setProg(opExe); e != nil {
				err = e
			}
		}
		return
	})
}

func (cat *Catalog) AssertDefinition(path ...string) error {
	return cat.Schedule(assert.PostDomain, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		if end := len(path) - 1; end <= 0 {
			err = errutil.New("path too short", path)
		} else {
			path, value := path[:end], path[end]
			err = d.AddDefinition(MakeKey(path...), at, value)
		}
		return
	})
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
