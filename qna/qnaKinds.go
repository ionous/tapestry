package qna

import (
	"database/sql"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

type qnaKinds struct {
	kinds                       qnaKindMap
	typeOf, fieldsOf, traitsFor *sql.Stmt // selects field, type for a named kind
	signatures                  []map[uint64]interface{}
}

type qnaKindMap map[string]qnaKind

type qnaKind struct {
	kind *g.Kind
	init []rt.Assignment
}

// aspects are a specific kind of record where every field is a boolean trait
func (q *qnaKinds) addKind(n string, k *g.Kind, init []rt.Assignment) *g.Kind {
	if q.kinds == nil {
		q.kinds = make(qnaKindMap)
	}
	q.kinds[n] = qnaKind{kind: k, init: init}
	return k
}

func (q *qnaKinds) GetKindByName(name string) (ret *g.Kind, err error) {
	// fix? breakcase doesnt normalize casing ( mainly b/c of case-aware patterns and kinds )
	// different cased names might therefore point to the same kind.
	name = lang.Breakcase(name)
	if k, ok := q.kinds[name]; ok {
		ret = k.kind
	} else {
		var n, path string
		if e := q.typeOf.QueryRow(name).Scan(&n, &path); e != nil {
			err = e
		} else {
			if len(path) == 0 {
				err = errutil.New("cant determine typeOf", name)
			} else if path == object.Aspect {
				// aspects also have kinds, but they are stored in their own table
				if ts, e := q.queryTraits(n); e != nil {
					err = e
				} else {
					// fix: inits would allow the "usually x" to vary per kind
					ret = q.addKind(n, g.NewKind(q, n, ts), nil)
				}
			} else {
				if ts, inits, e := q.queryFields(n); e != nil {
					err = e
				} else {
					ret = q.addKind(n, g.NewKind(q, path, ts), inits)
				}
			}
		}
	}
	if err != nil {
		err = errutil.Fmt("error while getting kind %q, %w", name, err)
	}
	return
}

func (q *qnaKinds) queryFields(kind string) (ret []g.Field, retInit []rt.Assignment, err error) {
	// creates the kind if it needs to.
	if rows, e := q.fieldsOf.Query(kind); e != nil {
		err = e
	} else {
		var field, fieldType string
		var affinity affine.Affinity
		var i interface{}
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
				if val, e := decodeValue(affinity, i, q.signatures); e != nil {
					err = errutil.New("error while decoding", field, e)
				} else {
					// add room for all earlier fields
					if retInit == nil {
						cnt := len(ret) // ; -1 so we can immediately append
						retInit = make([]rt.Assignment, cnt-1, cnt)
					}
					retInit = append(retInit, val)
				}
			}
			return
		}, &field, &fieldType, &affinity, &i)
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
