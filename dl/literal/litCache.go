package literal

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

// actually, record literals should already be encodedable right?
// its just that the compact encoder doesnt have a *special* handler.
// Name  string       `if:"label=name,type=text"`
// Value LiteralValue `if:"label=value"`

type RecordCache struct {
	rec *rt.Record // it might be more accurate to store this at an arbitrary string key in the runtime
	// that way, if multiple runtimes are used within the same process ( ex. tests )
	// each can be isolated from the other.
}

type RecordsCache struct {
	recs []*rt.Record
}

func (rc *RecordCache) GetRecord(kinds rt.Kinds, kind string, fields []FieldValue) (ret rt.Value, err error) {
	if rc.rec == nil {
		if k, e := kinds.GetKindByName(kind); e != nil {
			err = e
		} else if v, e := buildRecord(kinds, k, fields); e != nil {
			err = e
		} else {
			rc.rec = v
		}
	}
	if err == nil {
		ret = rt.RecordOf(rc.rec)
	}
	return
}

func (rc *RecordsCache) GetRecords(kinds rt.Kinds, kind string, els []FieldList) (ret rt.Value, err error) {
	if rc.recs == nil {
		if k, e := kinds.GetKindByName(kind); e != nil {
			err = e
		} else if vs, e := buildRecords(kinds, k, els); e != nil {
			err = e
		} else {
			rc.recs = vs
		}
	}
	if err == nil {
		ret = rt.RecordsFrom(rc.recs, kind)
	}
	return
}

func buildRecords(kinds rt.Kinds, k *rt.Kind, els []FieldList) (ret []*rt.Record, err error) {
	var out []*rt.Record
	for _, el := range els {
		if v, e := buildRecord(kinds, k, el.Fields); e != nil {
			err = e
			break
		} else {
			out = append(out, v)
		}
	}
	if err == nil {
		ret = out
	}
	return
}

// create a new record
// note: this doesnt translate traits to aspects under the theory there should be only one of each field in list of fields
func buildRecord(kinds rt.Kinds, k *rt.Kind, fields []FieldValue) (ret *rt.Record, err error) {
	out := rt.NewRecord(k)
	// fix? it might? make more sense to be able to create record with FieldValue(s) directly
	// to avoid the extra allocation -- to handle the slice conversion -- since you cant cast slices of types in go:
	// 1. have generic depend on an autogenerated type
	// 2. allow external types to be referenced by the autogeneration
	// 2. pass some sort of an iterator nextField() (string, rt.Value) (
	set := make([]bool, k.FieldCount())
	// fields of name, literal value
	for _, fv := range fields {
		if idx := k.FieldIndex(fv.Field); idx < 0 {
			err = errutil.Fmt("unknown field %q in kind %q", fv.Field, k.Name())
			break
		} else if set[idx] {
			err = errutil.New("duplicate fields set by literal", fv.Field)
			break
		} else if v, e := makeValue(kinds, k.Field(idx), fv.Value); e != nil {
			err = e
			break
		} else if e := out.SetIndexedField(idx, v); e != nil {
			err = e
		} else {
			set[idx] = true
		}
	}
	if err == nil {
		ret = out
	}
	return
}

func makeValue(kinds rt.Kinds, ft rt.Field, val LiteralValue) (ret rt.Value, err error) {
	if fvs, ok := val.(*FieldList); !ok {
		ret, err = val.GetLiteralValue(kinds)
	} else if fvk, e := kinds.GetKindByName(ft.Type); e != nil {
		err = e
	} else if c, e := buildRecord(kinds, fvk, fvs.Fields); e != nil {
		err = e
	} else {
		ret = rt.RecordOf(c)
	}
	return
}
