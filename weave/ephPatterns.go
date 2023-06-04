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
			pat := d.EnsureKind(name, at)
			pat.AddRequirement(kindsOf.Pattern.String())
			// schedule results after than parameters for the sake of sorting
			err = d.schedule(at, assert.ResultPhase, func(ctx *Weaver) (err error) {
				if pat.domain != d {
					err = errutil.New("can only declare results in the original domain")
				} else if init != nil {
					err = errutil.New("return values don't currently support initial values")
				} else if res, e := MakeUniformField(aff, fieldName, class, at); e != nil {
					err = e
				} else {
					err = cat.writer.Result(d.name, pat.name, res.Name, res.Affinity, res.Type, at)
				}
				return
			})
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
			pat := d.EnsureKind(name, at)
			pat.AddRequirement(kindsOf.Pattern.String())
			if pat.domain != d {
				err = errutil.New("can only declare args in the original domain")
			} else if init != nil {
				err = errutil.New("return values dont currently support initial values")
			} else if arg, e := MakeUniformField(aff, fieldName, class, at); e != nil {
				err = e
			} else {
				err = cat.Schedule(assert.ParamPhase, func(ctx *Weaver) error {
					return cat.writer.Parameter(d.name, pat.name, arg.Name, arg.Affinity, arg.Type, at)
				})
			}
		}
		return
	})
}
