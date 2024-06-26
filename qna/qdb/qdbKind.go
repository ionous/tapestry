package qdb

import (
	"database/sql"
	"fmt"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/qna/decoder"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables"
)

func (q *Query) GetKindByName(k string) (ret *rt.Kind, err error) {
	key := query.MakeKey("kinds", k, "")
	if c, e := q.constVals.Ensure(key, func() (ret any, err error) {
		ret, err = q.buildKind(k)
		return
	}); e != nil {
		err = e
	} else {
		ret = c.(*rt.Kind)
	}
	return
}

func (q *Query) getKindByNoun(fullname string) (ret *rt.Kind, err error) {
	if k, e := scanString(q.nounKind, fullname); e != nil {
		err = e
	} else {
		ret, err = q.GetKindByName(k)
	}
	return
}

func (q *Query) buildKind(k string) (ret *rt.Kind, err error) {
	if path, e := q.KindOfAncestors(k); e != nil {
		err = fmt.Errorf("error while getting kind %q, %w", k, e)
	} else if len(path) == 0 {
		err = fmt.Errorf("invalid kind %q", k)
	} else {
		k := path[0] // use the returned name in favor of the given name (ex. plurals)
		if fields, e := q.getCachedFields(k); e != nil {
			err = fmt.Errorf("error while building kind %q, %w", k, e)
		} else {
			var aspects []rt.Aspect // fix? currently only allow traits for objects. hrm.
			if objectLike := path[len(path)-1] == kindsOf.Kind.String(); objectLike {
				aspects = rt.MakeAspects(q, fields)
			}
			ret = &rt.Kind{Path: path, Fields: fields, Aspects: aspects}
		}
	}
	return
}

// cached fields exclusive to a kind
func (q *Query) getCachedFields(kind string) (ret []rt.Field, err error) {
	key := query.MakeKey("fields", kind, "")
	if c, e := q.constVals.Ensure(key, func() (ret any, err error) {
		return q.FieldsOf(kind)
	}); e != nil {
		err = e
	} else {
		ret = c.([]rt.Field)
	}
	return
}

// get uncached fields exclusive to a particular kind
func (q *Query) FieldsOf(kind string) (ret []rt.Field, err error) {
	// select  kind, field, value from (
	// 	select *, max(final) over (PARTITION by kind,field) as best
	// 	from mdl_value_kind
	// ) where final = best
	if rows, e := q.fieldsOf.Query(kind); e != nil {
		err = e
	} else {
		var last string
		var f struct {
			Name     string
			Affinity affine.Affinity
			Class    sql.NullString
			Init     []byte
		}
		err = tables.ScanAll(rows, func() (err error) {
			// the same field might be listed twice:
			// the final value ( listed first )
			// and a non final value ( listed second )
			if last != f.Name {
				if init, e := decodeInit(q.dec, f.Affinity, f.Init); e != nil {
					err = e
				} else {
					ret = append(ret, rt.Field{
						Name:     f.Name,
						Affinity: f.Affinity,
						Type:     f.Class.String,
						Init:     init,
					})
					last = f.Name
				}
			}
			return
		}, &f.Name, &f.Affinity, &f.Class.String, &f.Init)
	}
	return
}

// decode the passed assignment, if it exists.
func decodeInit(d decoder.Decoder, aff affine.Affinity, b []byte) (ret rt.Assignment, err error) {
	if len(b) > 0 {
		ret, err = d.DecodeAssignment(aff, b)
	}
	return
}
