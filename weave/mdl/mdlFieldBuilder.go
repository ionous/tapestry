package mdl

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/lang"
	"github.com/ionous/errutil"
	"github.com/ionous/sliceOf"
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
	aspects []string
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

// creates not/is implicit aspects
// tbd: would it be nicer to support single trait kinds?
// not_aspect would instead be: Not{IsTrait{PositiveName}}
func (b *FieldBuilder) AddAspect(ak string) {
	b.aspects = append(b.aspects, ak)
}

func (fs *Fields) writeFields(pen *Pen) (err error) {
	if kind, e := pen.findRequiredKind(fs.kind); e != nil {
		err = e
	} else if cache, e := fs.fieldSet.cache(pen); e != nil {
		err = e
	} else if e := fs.fieldSet.writeFieldSet(pen, kind, cache); e != nil {
		err = e
	} else {
		// generate implicit aspects
		for _, ak := range fs.aspects {
			if cls, e := pen.addAspect(ak, sliceOf.String("not "+ak, "is "+ak)); e != nil {
				err = e
				break
			} else if e := pen.addField(kind, cls, ak, affine.Text); e != nil {
				err = e
				break
			}
		}
	}
	if err != nil {
		err = errutil.Fmt("%w in pattern %q domain %q", err, fs.kind, pen.domain)
	}
	return
}

// future: wrap up the in progress field set with a "promise"
// to avoid looking up the same info repeatedly
type classCache map[string]kindInfo

// for patterns, waits to create the pattern after all fields are known
// that ensures that "extend pattern" (to add locals) happens after define pattern (for parameters and locals)
func (fs *fieldSet) cache(pen *Pen) (ret classCache, err error) {
	for ft, fields := range fs.fields {
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

func (fs *fieldSet) writeFieldSet(pen *Pen, kid kindInfo, cache classCache) (err error) {
	type fn func(kid, cls kindInfo, field string, aff affine.Affinity) error
	var out = [3]fn{
		pen.addParameter, // PatternParameters
		pen.addResult,
		pen.addField,
	}
	for ft, fields := range fs.fields {
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
func (fs *FieldInfo) getClass() (ret string) {
	if cls := fs.Class; len(cls) > 0 {
		ret = cls
	} else if isRecordAffinity(fs.Affinity) {
		ret = fs.Name
	}
	return
}

func (fs *FieldInfo) validate(ft FieldType) (err error) {
	if len(fs.Name) == 0 {
		err = errutil.New("missing name")
	} else if len(fs.Affinity) == 0 {
		err = errutil.New("missing affinity")
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
