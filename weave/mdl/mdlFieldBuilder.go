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
	typeOfs,
	fieldOfs int
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
	} else if e := p.fieldSet.writeFieldSet(pen, kind); e != nil {
		err = e
	}
	return
}

func (p *fieldSet) writeFieldSet(pen *Pen, kid kindInfo) (err error) {
	type fn func(kid, cls kindInfo, field string, aff affine.Affinity) error
	var out = [3]fn{
		pen.addParameter, // PatternParameters
		pen.addResult,
		pen.addField,
	}
	for typeCnt := len(p.fields); p.typeOfs < typeCnt; {
		fields, call, ft := p.fields[p.typeOfs], out[p.typeOfs], FieldType(p.typeOfs)
		for elCount := len(fields); p.fieldOfs < elCount; {
			field := fields[p.fieldOfs]
			if e := field.validate(ft); e != nil {
				err = e
			} else if cls, e := pen.findOptionalKind(field.getClass()); e != nil {
				err = e
			} else {
				e := call(kid, cls, field.Name, field.Affinity)
				if err = eatDuplicates(pen.warn, e); err == nil && field.Init != nil {
					e := pen.addDefault(kid, field.Name, field.Init)
					err = eatDuplicates(pen.warn, e)
				}
			}

			if err != nil {
				err = errutil.Fmt("%w trying to write field %q in kind %q, domain %q", err, field.Name,
					kid.name, pen.domain)
				break
			}
			p.fieldOfs++
		}
		if err != nil {
			break
		}
		p.fieldOfs = 0 // reset for the next go around.
		p.typeOfs++
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
	switch ft {
	case PatternResults:
		if p.Init != nil {
			err = errutil.New("pattern returns currently don't support initial values")
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
