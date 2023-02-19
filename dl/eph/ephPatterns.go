package eph

import (
	"strings"

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
				labels := k.domain.GetDefinition(MakeKey("pat", pat, "args"))
				//
				if e := w.Write(mdl.Pat, k.domain.name, k.name, labels.value, result.value); e != nil {
					err = e
					break
				}
			}
		}
	}
	return
}

// pattern commands generate ancestry, fields, and pattern entries....
func (op *EphPatterns) Phase() Phase { return AncestryPhase }

func (op *EphPatterns) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if name, ok := UniformString(op.Name); !ok {
		err = InvalidString(op.Name)
	} else {
		k := d.EnsureKind(name, at)
		k.AddRequirement(kindsOf.Pattern.String())
		var locals []UniformField
		if e := op.assembleRes(d, k, at, &k.patternHeader); e != nil {
			err = e
		} else if e := op.assembleArgs(d, k, at, &k.patternHeader); e != nil {
			err = e
		} else if e := reduceLocals(op.Locals, at, &locals); e != nil {
			err = e
		} else {
			err = d.AddEphemera(at, PhaseFunction{PropertyPhase,
				func(c *Catalog, d *Domain, at string) (err error) {
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

func (pd *patternHeader) flush() (ret []UniformField) {
	if !pd.written {
		ret = append(pd.args, pd.res...)
		pd.written = true
	}
	return ret
}

// writes a definition of patternName?res=name
func (op *EphPatterns) assembleRes(d *Domain, k *ScopedKind, at string, outp *patternHeader) (err error) {
	var res []UniformField
	if op.Result != nil && k.domain != d {
		err = errutil.New("can only declare results in the original domain")
	} else if patres, e := reduceRes(op.Result, at, &res); e != nil {
		err = e
	} else if len(patres) > 0 {
		if e := addPatternDef(d, k, "res", at, patres); e != nil {
			err = e
		} else {
			outp.res = res
		}
	}
	return
}

// writes a definition of patternName?args=arg1,arg2,arg3
func (op *EphPatterns) assembleArgs(d *Domain, k *ScopedKind, at string, outp *patternHeader) (err error) {
	var args []UniformField
	if len(op.Params) > 0 && k.domain != d {
		err = errutil.New("can only declare args in the original domain")
	} else if patlabels, e := reduceArgs(op.Params, at, &args); e != nil {
		err = e
	} else if len(patlabels) > 0 {
		if e := addPatternDef(d, k, "args", at, patlabels); e != nil {
			err = e
		} else {
			outp.args = args
		}
	}
	return
}

func addPatternDef(d *Domain, k *ScopedKind, key, at, v string) (err error) {
	if k.domain != d {
		err = DomainError{d.name, errutil.Fmt("expected the pattern %q and its %s to be defined in the same domain (%q)", k.name, key, k.domain.name)}
	} else if e := d.AddDefinition(MakeKey("pat", k.name, key), at, v); e != nil {
		err = e // use definition to block the pattern from defining different args
	}
	return
}

func reduceRes(param *EphParams, at string, outp *[]UniformField) (ret string, err error) {
	if param != nil {
		if param.Initially.IsValid() {
			err = errutil.New("return values dont currently support initial values")
		} else if p, e := param.unify(at); e != nil {
			err = e
		} else {
			*outp = append(*outp, p)
			ret = p.name
		}
	}
	return
}

func reduceArgs(params []EphParams, at string, outp *[]UniformField) (ret string, err error) {
	var labels strings.Builder
	for i, param := range params {
		if param.Initially.IsValid() {
			err = errutil.New("args dont currently support initial values")
		} else if p, e := param.unify(at); e != nil {
			err = e
			break
		} else {
			*outp = append(*outp, p)
			if i > 0 {
				labels.WriteRune(',') // join
			}
			// fix: eventually, labels might be different than the field names
			labels.WriteString(p.name)
		}
	}
	if err == nil {
		ret = labels.String()
	}
	return
}

func reduceLocals(params []EphParams, at string, outp *[]UniformField) (err error) {
	for _, param := range params {
		if p, e := param.unify(at); e != nil {
			err = e
			break
		} else if e := p.setAssignment(param.Initially); e != nil {
			err = e
		} else {
			*outp = append(*outp, p)
		}
	}
	return
}
