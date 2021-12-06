package eph

import (
	"strings"
)

// write the pattern table in a reasonable order
func (c *Catalog) WritePatterns(w Writer) (err error) {
	if deps, e := c.ResolveKinds(); e != nil {
		err = e
	} else {
		for _, dep := range deps {
			k, parents := dep.Leaf().(*ScopedKind), dep.Parents()
			if len(parents) == 1 && parents[0].Name() == KindsOfPattern {
				pat := k.name
				defs := k.domain.phases[FieldPhase].defs
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
		kind := d.EnsureKind(name, at)
		kind.AddRequirement(KindsOfPattern)
		// size of fields is usually > 0, so dont worry too much about the edge case
		fields := make([]UniformField, 0, measurePattern(el))
		if e := el.assembleRet(d, name, at, &fields); e != nil {
			err = e
		} else if e := el.assembleArgs(d, name, at, &fields); e != nil {
			err = e
		} else if e := reduceLocals(el.Locals, &fields); e != nil {
			err = e
		} else {
			err = d.AddEphemera(EphAt{at, PhaseFunction{FieldPhase,
				func(c *Catalog, d *Domain, at string) (err error) {
					for _, p := range fields {
						if e := p.AssembleField(kind, at); e != nil {
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

func (el *EphPatterns) assembleRet(d *Domain, name, at string, outp *[]UniformField) (err error) {
	if patret, e := reduceRet(el.Return, outp); e != nil {
		err = e
	} else if len(patret) > 0 {
		// use definition to block the pattern from defining a different return
		if e := d.AddDefinition(name+"?res", at, patret); e != nil {
			err = e
		}
	}
	return
}

func (el *EphPatterns) assembleArgs(d *Domain, name, at string, outp *[]UniformField) (err error) {
	if patlabels, e := reduceArgs(el.Args, outp); e != nil {
		err = e
	} else if len(patlabels) > 0 {
		// use definition to block the pattern from defining different args
		if e := d.AddDefinition(name+"?args", at, patlabels); e != nil {
			err = e
		}
	}
	return
}

func reduceRet(param *EphParams, outp *[]UniformField) (ret string, err error) {
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
	if el.Return != nil {
		ret++
	}
	return
}
