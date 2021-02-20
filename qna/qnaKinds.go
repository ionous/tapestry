package qna

import (
	"database/sql"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/object"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

type qnaKinds struct {
	kinds                       map[string]*g.Kind
	values                      valueMap  // kind.field => rt.Assignment
	typeOf, fieldsOf, traitsFor *sql.Stmt // selects field, type for a named kind
}

// aspects are a specific kind of record where every field is a boolean trait
func (q *qnaKinds) addKind(n string, k *g.Kind) *g.Kind {
	if q.kinds == nil {
		q.kinds = make(map[string]*g.Kind)
	}
	q.kinds[n] = k
	return k
}

func (q *qnaKinds) GetKindByName(name string) (ret *g.Kind, err error) {
	// fix? breakcase doesnt normalize casing ( mainly b/c of case-aware patterns and kinds )
	// different cased names might therefore point to the same kind.
	name = lang.Breakcase(name)
	if k, ok := q.kinds[name]; ok {
		ret = k
	} else {
		var n, path string
		if e := q.typeOf.QueryRow(name).Scan(&n, &path); e != nil {
			err = e
		} else {
			if len(path) == 0 {
				err = errutil.New("cant determine typeOf", name)
			} else if path == object.Aspect {
				if ts, e := q.queryTraits(n); e != nil {
					err = e
				} else {
					ret = q.addKind(n, g.NewKind(q, n, ts))
				}
			} else {
				if ts, e := q.queryFields(n); e != nil {
					err = e
				} else {
					ret = q.addKind(n, g.NewKind(q, path, ts))
				}
			}
		}
	}
	if err != nil {
		err = errutil.Fmt("error while getting kind %q, %w", name, err)
	}
	return
}

func (q *qnaKinds) queryFields(kind string) (ret []g.Field, err error) {
	// creates the kind if it needs to.
	if rows, e := q.fieldsOf.Query(kind); e != nil {
		err = e
	} else {
		var field, fieldType string
		var affinity affine.Affinity
		var i interface{}
		var hasInit bool
		err = tables.ScanAll(rows, func() (err error) {
			// by default the type and the affinity are the same
			// ( like reflect, where type and kind are the same for primitive types )
			if len(affinity) == 0 {
				affinity = affine.Affinity(fieldType)
			}
			ret = append(ret, g.Field{
				Name:     field,
				Affinity: affinity,
				Type:     fieldType,
			})
			if i != nil {
				if val, e := decodeValue(affinity, i); e != nil {
					err = errutil.New("error while decoding", field, e)
				} else {
					key := keyType{kind, field}
					q.values[key] = val
				}
				hasInit = true // for debugging
			}
			return
		}, &field, &fieldType, &affinity, &i)
		//
		hasInit = hasInit
	}
	return
}

func (q *qnaKinds) queryTraits(name string) (ret []g.Field, err error) {
	var trait string
	if rows, e := q.traitsFor.Query(name); e != nil {
		err = e
	} else {
		err = tables.ScanAll(rows, func() (err error) {
			ret = append(ret, g.Field{
				Name:     trait,
				Affinity: affine.Bool,
				Type:     "trait",
			})
			return
		}, &trait)
	}
	return
}
