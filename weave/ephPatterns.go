package weave

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/weave/assert"

	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"github.com/ionous/errutil"
)

func (cat *Catalog) AssertResult(patternName, fieldName, class string, aff affine.Affinity, init assign.Assignment) error {
	return cat.Schedule(assert.AncestryPhase, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		if name, ok := UniformString(patternName); !ok {
			err = InvalidString(patternName)
		} else {
			name := d.singularize(name)
			if e := d.addKind(name, kindsOf.Pattern.String(), at); e != nil {
				err = e
			} else if init != nil {
				err = errutil.New("return values don't currently support initial values")
			} else if res, e := MakeUniformField(aff, fieldName, class, at); e != nil {
				err = e
			} else {
				err = d.schedule(at, assert.ResultPhase, func(ctx *Weaver) error {
					return cat.writer.Result(d.name, name, res.Name, res.Affinity, res.Type, at)
				})
			}
		}
		return
	})
}

// writes a definition of patternName?args=arg1,arg2,arg3
func (cat *Catalog) AssertParam(patternName, fieldName, class string, aff affine.Affinity, init assign.Assignment) error {
	return cat.Schedule(assert.AncestryPhase, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		if name, ok := UniformString(patternName); !ok {
			err = InvalidString(patternName)
		} else {
			name := d.singularize(name)
			if e := d.addKind(name, kindsOf.Pattern.String(), at); e != nil {
				err = e
			} else if init != nil {
				err = errutil.New("parameters values don't currently support initial values")
			} else if arg, e := MakeUniformField(aff, fieldName, class, at); e != nil {
				err = e
			} else {
				err = cat.Schedule(assert.ParamPhase, func(ctx *Weaver) error {
					return cat.writer.Parameter(d.name, name, arg.Name, arg.Affinity, arg.Type, at)
				})
			}
		}
		return
	})
}
