package eph

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/imp/assert"

	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/ionous/errutil"
)

// write the pattern table in a reasonable order
func (c *Catalog) WritePatterns(w Writer) (err error) {
	if deps, e := c.ResolveKinds(); e != nil {
		err = e
	} else {
		for _, dep := range deps {
			if k := dep.Leaf().(*ScopedKind); k.HasAncestor(kindsOf.Pattern) {
				pat := k.name
				result := k.domain.GetDefinition(MakeKey("pat", pat, "res"))
				labels := k.patternHeader.labels()
				//
				if e := w.Write(mdl.Pat, k.domain.name, k.name, labels, result.value); e != nil {
					err = e
					break
				}
			}
		}
	}
	return
}

// pattern commands generate ancestry, fields, and pattern entries....
func (op *EphPatterns) Phase() assert.Phase { return assert.AncestryPhase }

func (op *EphPatterns) Weave(k assert.Assertions) (err error) {
	name := op.PatternName
	if e := k.AssertAncestor(name, kindsOf.Pattern.String()); e != nil {
		err = e
	} else {
		if ps := op.Params; err == nil && len(ps) > 0 {
			err = weaveFields(name, ps, k.AssertParam)
		}
		if ps := op.Result; err == nil && ps != nil {
			err = weaveField(name, *ps, k.AssertResult)
		}
		if ps := op.Locals; err == nil {
			err = weaveFields(name, ps, k.AssertLocal)
		}
	}
	return
}

func (ctx *Context) AssertLocal(patternName, fieldName, class string, aff affine.Affinity, init assign.Assignment) (err error) {
	d, at := ctx.d, ctx.at
	if name, ok := UniformString(patternName); !ok {
		err = InvalidString(patternName)
	} else {
		k := d.EnsureKind(name, at)
		k.AddRequirement(kindsOf.Pattern.String())

		if p, e := MakeUniformField(affineToAffinity(aff), fieldName, class, at); e != nil {
			err = e
		} else if e := p.setAssignment(init); e != nil {
			err = e
		} else {
			var locals []UniformField
			locals = append(locals, p)
			//
			err = d.QueueEphemera(at, PhaseFunction{assert.PropertyPhase,
				func(assert.World, assert.Assertions) (err error) {
					k.pendingFields = append(k.pendingFields, k.patternHeader.flush()...)
					k.pendingFields = append(k.pendingFields, locals...)
					return
				}})
		}
	}
	return
}

// accumulate the various bits of pattern data
// ensure they get written correctly, and in a good order.
type patternHeader struct {
	res, args []UniformField
	written   bool
}

func (pd *patternHeader) labels() (ret string) {
	var b strings.Builder
	for i, el := range pd.args {
		if i > 0 {
			b.WriteRune(',')
		}
		b.WriteString(el.Name)
	}
	return b.String()
}

func (pd *patternHeader) flush() (ret []UniformField) {
	if !pd.written {
		ret = append(pd.args, pd.res...)
		pd.written = true
	}
	return ret
}

func (ctx *Context) AssertResult(patternName, fieldName, class string, aff affine.Affinity, init assign.Assignment) (err error) {
	d, at := ctx.d, ctx.at
	if name, ok := UniformString(patternName); !ok {
		err = InvalidString(patternName)
	} else {
		k := d.EnsureKind(name, at)
		k.AddRequirement(kindsOf.Pattern.String())

		if k.domain != d {
			err = errutil.New("can only declare results in the original domain")
		} else if init != nil {
			err = errutil.New("return values dont currently support initial values")
		} else if res, e := MakeUniformField(affineToAffinity(aff), fieldName, class, at); e != nil {
			err = e
		} else if e := addPatternDef(d, k, "res", at, res.Name); e != nil {
			err = e
		} else {
			k.patternHeader.res = []UniformField{res}
		}
	}
	return
}

// writes a definition of patternName?args=arg1,arg2,arg3
func (ctx *Context) AssertParam(patternName, fieldName, class string, aff affine.Affinity, init assign.Assignment) (err error) {
	d, at := ctx.d, ctx.at
	if name, ok := UniformString(patternName); !ok {
		err = InvalidString(patternName)
	} else {
		k := d.EnsureKind(name, at)
		k.AddRequirement(kindsOf.Pattern.String())
		if k.domain != d {
			err = errutil.New("can only declare args in the original domain")
		} else if init != nil {
			err = errutil.New("return values dont currently support initial values")
		} else if arg, e := MakeUniformField(affineToAffinity(aff), fieldName, class, at); e != nil {
			err = e
		} else {
			// there used to be one set of args, now there are individual args
			// if e := addPatternDef(d, k, "args", at, patlabels); e != nil {
			// else...
			// fix: this should probably check that no locals have been written yet
			// and/or use the "result" to seal in the args.
			k.patternHeader.args = append(k.patternHeader.args, arg)
		}
	}
	return
}

func addPatternDef(d *Domain, k *ScopedKind, key, at, v string) (err error) {
	if k.domain != d {
		err = domainError{d.name, errutil.Fmt("expected the pattern %q and its %s to be defined in the same domain (%q)", k.name, key, k.domain.name)}
	} else if e := d.AddDefinition(MakeKey("pat", k.name, key), at, v); e != nil {
		err = e // use definition to block the pattern from defining different args
	}
	return
}
