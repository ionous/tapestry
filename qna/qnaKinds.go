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
	kinds                        map[string]*g.Kind
	typeOf, fieldsFor, traitsFor *sql.Stmt // selects field, type for a named kind
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
	name = lang.Breakcase(name)
	if k, ok := q.kinds[name]; ok {
		ret = k
	} else {
		var role string
		if e := q.typeOf.QueryRow(name).Scan(&role); e != nil {
			err = e
		} else {
			if len(role) == 0 {
				err = errutil.New("missing role")
			} else if role == object.Aspect {
				if ts, e := q.queryTraits(name); e != nil {
					err = e
				} else {
					ret = q.addKind(name, g.NewKind(q, name, ts))
				}
			} else {
				if ts, e := q.queryFields(name); e != nil {
					err = e
				} else {
					ret = q.addKind(name, g.NewKind(q, role, ts))
				}
			}
		}
	}
	if err != nil {
		err = errutil.Fmt("error while getting kind %q, %w", name, err)
	}
	return
}

func (q *qnaKinds) queryFields(name string) (ret []g.Field, err error) {
	// creates the kind if it needs to.
	var field, fieldType string
	var affinity affine.Affinity
	if rows, e := q.fieldsFor.Query(name); e != nil {
		err = e
	} else {
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
			return
		}, &field, &fieldType, &affinity)
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
