package eph

import (
	"strings"

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
				defs := k.domain.phases[AncestryPhase].defs
				result := defs[pat+"?res"].value
				labels := defs[pat+"?args"].value
				//
				if e := w.Write(mdl_pat, k.domain.name, k.name, labels, result); e != nil {
					err = e
					break
				}
			}
		}
	}
	return
}

// pattern commands generate ancestry, fields, and pattern entries....
func (el *EphPatterns) Phase() Phase { return AncestryPhase }

func (el *EphPatterns) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if name, ok := UniformString(el.Name); !ok {
		err = InvalidString(el.Name)
	} else {
		k := d.EnsureKind(name, at)
		k.AddRequirement(KindsOfPattern)
		// size of fields is usually > 0, so dont worry too much about the edge case
		fields := make([]UniformField, 0, measurePattern(el))
		if e := el.assembleRet(d, k, at, &fields); e != nil {
			err = e
		} else if e := el.assembleArgs(d, k, at, &fields); e != nil {
			err = e
		} else if e := reduceLocals(el.Locals, &fields); e != nil {
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

func (el *EphPatterns) assembleRet(d *Domain, k *ScopedKind, at string, outp *[]UniformField) (err error) {
	if patres, e := reduceRes(el.Result, outp); e != nil {
		err = e
	} else if len(patres) > 0 {
		err = addPatternDef(d, k, "res", at, patres)
	}
	return
}

func (el *EphPatterns) assembleArgs(d *Domain, k *ScopedKind, at string, outp *[]UniformField) (err error) {
	if patlabels, e := reduceArgs(el.Args, outp); e != nil {
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
		if p, e := MakeUniformField(*param); e != nil {
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
		if p, e := MakeUniformField(param); e != nil {
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
		if p, e := MakeUniformField(param); e != nil {
			err = e
			break
		} else {
			*outp = append(*outp, p)
		}
	}
	return
}

func measurePattern(el *EphPatterns) (ret int) {
	ret = len(el.Args) + len(el.Locals)
	if el.Result != nil {
		ret++
	}
	return
}
