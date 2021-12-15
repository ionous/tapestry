package eph

import (
	"strings"

	"git.sr.ht/~ionous/iffy/tables/mdl"
	"github.com/ionous/errutil"
)

// write the pattern table in a reasonable order
func (c *Catalog) WritePatterns(w Writer) (err error) {
	if deps, e := c.ResolveKinds(); e != nil {
		err = e
	} else {
		for _, dep := range deps {
			if k := dep.Leaf().(*ScopedKind); k.HasAncestor(KindsOfPattern) {
				pat := k.name
				result := k.domain.GetDefinition(AncestryPhase, pat+"?res")
				labels := k.domain.GetDefinition(AncestryPhase, pat+"?args")
				//
				if e := w.Write(mdl.Pat, k.domain.name, k.name, labels, result); e != nil {
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
		k.AddRequirement(KindsOfPattern)
		// size of fields is usually > 0, so dont worry too much about the edge case
		fields := make([]UniformField, 0, measurePattern(op))
		if e := op.assembleRet(d, k, at, &fields); e != nil {
			err = e
		} else if e := op.assembleArgs(d, k, at, &fields); e != nil {
			err = e
		} else if e := reduceLocals(op.Locals, &fields); e != nil {
			err = e
		} else {
			err = d.AddEphemera(EphAt{at, PhaseFunction{FieldPhase,
				func(c *Catalog, d *Domain, at string) (err error) {
					for _, p := range fields {
						if e := p.AssembleField(k, at); e != nil {
							err = e
							break
						}
					}
					return
				}}})
		}
	}
	return
}

// writes a definition of patternName?res=name
func (op *EphPatterns) assembleRet(d *Domain, k *ScopedKind, at string, outp *[]UniformField) (err error) {
	if patres, e := reduceRes(op.Result, outp); e != nil {
		err = e
	} else if len(patres) > 0 {
		err = addPatternDef(d, k, "res", at, patres)
	}
	return
}

// writes a definition of patternName?args=arg1,arg2,arg3
func (op *EphPatterns) assembleArgs(d *Domain, k *ScopedKind, at string, outp *[]UniformField) (err error) {
	if patlabels, e := reduceArgs(op.Params, outp); e != nil {
		err = e
	} else if len(patlabels) > 0 {
		err = addPatternDef(d, k, "args", at, patlabels)
	}
	return
}

func addPatternDef(d *Domain, k *ScopedKind, key, at, v string) (err error) {
	if k.domain != d {
		err = DomainError{d.name, errutil.Fmt("expected the pattern %q and its %s to be defined in the same domain (%q)", k.name, key, k.domain.name)}
	} else if e := d.AddDefinition(k.name+"?"+key, at, v); e != nil {
		err = e // use definition to block the pattern from defining different args
	}
	return
}

func reduceRes(param *EphParams, outp *[]UniformField) (ret string, err error) {
	if param != nil {
		if param.Initially != nil {
			err = errutil.New("return values dont currently support initial values")
		} else if p, e := MakeUniformField(param.Affinity, param.Name, param.Class); e != nil {
			err = e
		} else {
			*outp = append(*outp, p)
			ret = p.name
		}
	}
	return
}

func reduceArgs(params []EphParams, outp *[]UniformField) (ret string, err error) {
	var labels strings.Builder
	for i, param := range params {
		if param.Initially != nil {
			err = errutil.New("args dont currently support initial values")
		} else if p, e := MakeUniformField(param.Affinity, param.Name, param.Class); e != nil {
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

func reduceLocals(params []EphParams, outp *[]UniformField) (err error) {
	for _, param := range params {
		if p, e := MakeUniformField(param.Affinity, param.Name, param.Class); e != nil {
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

func measurePattern(op *EphPatterns) (ret int) {
	ret = len(op.Params) + len(op.Locals)
	if op.Result != nil {
		ret++
	}
	return
}
