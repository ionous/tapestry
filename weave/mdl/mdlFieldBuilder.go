package mdl

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"github.com/ionous/errutil"
)

type FieldType int

//go:generate stringer -type=FieldType -linecomment
const (
	PatternParameters FieldType = iota // pattern parameters
	PatternResults
	PatternLocals
	NumFieldTypes
)

type FieldInfo struct {
	Name, Class string
	Affinity    affine.Affinity
	Init        assign.Assignment
}

type fieldSet struct {
	fields [NumFieldTypes][]FieldInfo
	typeOfs,
	fieldOfs int
}

func (p *fieldSet) write(m *Pen, classes classCache, kid KindInfo) (err error) {
	// type fn func(kid, cls KindInfo, field string, aff affine.Affinity) error
	// var out = [3]fn{
	// 	m.AddParameterById, // PatternParameters
	// 	m.AddResultById,
	// 	m.AddField,
	// }
	// for cnt := len(p.fields); err == nil && p.typeOfs < cnt; p.typeOfs++ {
	// 	fields, call := p.fields[p.typeOfs], out[p.typeOfs]
	// 	for cnt := len(fields); err == nil && p.fieldOfs < cnt; p.fieldOfs++ {
	// 		field, ft := fields[p.fieldOfs], FieldType(p.fieldOfs)
	// 		if e := field.validate(ft); e != nil {
	// 			err = e
	// 		} else if cls, e := classes.getClass(field.Class); e != nil {
	// 			err = e
	// 		} else if e := call(kid, cls, field.Name, field.Affinity); e != nil {
	// 			err = e
	// 		}
	// 	}
	// }
	panic("not implemented")
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
