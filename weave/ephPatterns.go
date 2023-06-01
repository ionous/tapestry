package weave

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave/assert"

	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"github.com/ionous/errutil"
)

// write the pattern table in a reasonable order
func (c *Catalog) WritePatterns(m mdl.Modeler) (err error) {
	if deps, e := c.ResolveKinds(); e != nil {
		err = e
	} else {
		for _, dep := range deps {
			if k := dep.Leaf().(*ScopedKind); k.HasAncestor(kindsOf.Pattern) {
				pat := k.name
				result := k.domain.GetDefinition(MakeKey("pat", pat, "res"))
				labels := k.header.labels()
				//
				if e := m.Pat(k.domain.name, k.name, labels, result.value); e != nil {
					err = e
					break
				}
			}
		}
	}
	return
}

func (cat *Catalog) AssertResult(patternName, fieldName, class string, aff affine.Affinity, init assign.Assignment) error {
	return cat.Schedule(assert.AncestryPhase, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		if name, ok := UniformString(patternName); !ok {
			err = InvalidString(patternName)
		} else {
			pat := d.EnsureKind(name, at)
			pat.AddRequirement(kindsOf.Pattern.String())
			// schedule locals later than parameters for the sake of sorting
			err = d.schedule(at, assert.ResultPhase, func(ctx *Weaver) (err error) {
				if pat.domain != d {
					err = errutil.New("can only declare results in the original domain")
				} else if init != nil {
					err = errutil.New("return values dont currently support initial values")
				} else if res, e := MakeUniformField(aff, fieldName, class, at); e != nil {
					err = e
				} else if e := addPatternDef(d, pat, "res", at, res.Name); e != nil {
					err = e
				} else {
					pat.header.resList = []UniformField{res}
					err = cat.Schedule(assert.ResultPhase, func(ctx *Weaver) error {
						return cat.writeField(d.name, pat.name, res)
					})
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
				pat.header.paramList = append(pat.header.paramList, arg)
				err = cat.Schedule(assert.ParamPhase, func(ctx *Weaver) error {
					return cat.writeField(d.name, pat.name, arg)
				})
			}
		}
		return
	})
}

func addPatternDef(d *Domain, k *ScopedKind, key, at, v string) (err error) {
	if k.domain != d {
		err = domainError{d.name, errutil.Fmt("expected the pattern %q and its %s to be defined in the same domain (%q)", k.name, key, k.domain.name)}
	} else if e := d.AddDefinition(MakeKey("pat", k.name, key), at, v); e != nil {
		err = e // use definition to block the pattern from defining different args
	}
	return
}
