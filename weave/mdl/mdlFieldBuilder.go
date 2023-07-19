package mdl

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/lang"
	"github.com/ionous/errutil"
)

type FieldInfo struct {
	Name, Class string
	Affinity    affine.Affinity
	Init        assign.Assignment
}

type FieldBuilder struct {
	Fields
}

type Fields struct {
	kind string
	fieldSet
}

// supports enough slot for patterns
// see PatternBuilder.
type fieldSet struct {
	fields [NumFieldTypes][]FieldInfo
}

func NewFieldBuilder(kind string) *FieldBuilder {
	return &FieldBuilder{Fields: Fields{
		// tbd: feels like it'd be best to have spec flag names that need normalization,
		// and convert all the names at load time ( probably storing the original somewhere )
		// ( ex. store the normalized names in the meta data )
		kind: lang.Normalize(kind),
	}}
}

// defers execution; so no return value.
func (b *FieldBuilder) AddField(fn FieldInfo) {
	b.fields[PatternLocals] = append(b.fields[PatternLocals], fn)
}

func (b *FieldBuilder) AddAspect(ak string) {
	b.AddField(FieldInfo{
		Name:     ak,
		Class:    ak,
		Affinity: affine.Text,
	})
}

func (p *Fields) writeFields(pen *Pen) (err error) {
	if kind, e := pen.findRequiredKind(p.kind); e != nil {
		err = e
	} else if cache, e := p.fieldSet.cache(pen); e != nil {
		err = e
	} else if e := p.fieldSet.writeFieldSet(pen, kind, cache); e != nil {
		err = e
	}
	if err != nil {
		err = errutil.Fmt("%w in pattern %q domain %q", err, p.kind, pen.domain)
	}
	return
}

// future: wrap up the in progress field set with a "promise"
// to avoid looking up the same info repeatedly
type classCache map[string]kindInfo

// for patterns, waits to create the pattern after all fields are known
// that ensures that "extend pattern" (to add locals) happens after define pattern (for parameters and locals)
func (p *fieldSet) cache(pen *Pen) (ret classCache, err error) {
	for ft, fields := range p.fields {
		for _, field := range fields {
			if e := field.validate(FieldType(ft)); e != nil {
				err = e
			} else if clsName := field.getClass(); len(clsName) > 0 {
				if cls, e := pen.findOptionalKind(clsName); e != nil {
					err = e
				} else {
					if ret == nil {
						ret = make(classCache)
					}
					ret[clsName] = cls
				}
			}
			if err != nil {
				err = errutil.Fmt("%w trying to write field %q", err, field.Name)
				break
			}
		}
	}
	return
}

func (p *fieldSet) writeFieldSet(pen *Pen, kid kindInfo, cache classCache) (err error) {
	type fn func(kid, cls kindInfo, field string, aff affine.Affinity) error
	var out = [3]fn{
		pen.addParameter, // PatternParameters
		pen.addResult,
		pen.addField,
	}
	for ft, fields := range p.fields {
		call := out[ft]
		for _, field := range fields {
			cls := cache[field.getClass()]
			e := call(kid, cls, field.Name, field.Affinity)
			if e := eatDuplicates(pen.warn, e); e != nil {
				err = e
			} else if field.Init != nil {
				e := pen.addDefault(kid, field.Name, field.Init)
				if e := eatDuplicates(pen.warn, e); e != nil {
					err = e
				}
			}
			if err != nil {
				err = errutil.Fmt("%w trying to write field %q", err, field.Name)
				break
			}
		}
	}

	return
}

// shortcut: if we specify a field name for a record and no class, we'll expect the class to be the name.
func (p *FieldInfo) getClass() (ret string) {
	if cls := p.Class; len(cls) > 0 {
		ret = cls
	} else if isRecordAffinity(p.Affinity) {
		ret = p.Name
	}
	return
}

func (p *FieldInfo) validate(ft FieldType) (err error) {
	if len(p.Name) == 0 {
		err = errutil.New("missing name")
	} else if len(p.Affinity) == 0 {
		err = errutil.New("missing affinity")
	} else {
		switch ft {
		case PatternResults:
			if p.Init != nil {
				err = errutil.New("pattern returns currently don't support initial values")
			}
		}
	}
	return
}

// if there is a class specified, only certain affinities are allowed.
func isRecordAffinity(a affine.Affinity) (okay bool) {
	switch a {
	case affine.Record, affine.RecordList:
		okay = true
	}
	return
}
