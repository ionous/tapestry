package mdl

import (
	"database/sql"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

func (m *Writer) findCardinality(kind kindInfo) (ret string, err error) {
	if e := m.db.QueryRow(`
	select cardinality
	from mdl_rel
	where relKind = ?1 
	limit 1
	`, kind.id).Scan(&ret); e == sql.ErrNoRows {
		err = errutil.Fmt("unknown or invalid cardinality for %q in %q", kind.name, kind.domain)
	} else {
		err = e
	}
	return
}

type relKind struct {
	class  string
	plural bool
}

func (k *relKind) affinity() (ret affine.Affinity) {
	if k.plural {
		ret = affine.TextList
	} else {
		ret = affine.Text
	}
	return
}

func (k *relKind) lhs() (ret string) {
	if k.plural {
		ret = "kinds"
	} else {
		ret = "kind"
	}
	return
}

func (k *relKind) rhs() (ret string) {
	if k.plural {
		ret = "other_kinds"
	} else {
		ret = "other_kind"
	}
	return
}

func makeRel(a, b string, card string) (first, second relKind) {
	switch card {
	case tables.ONE_TO_ONE:
		first = relKind{a, false}
		second = relKind{b, false}
	case tables.ONE_TO_MANY:
		first = relKind{a, false}
		second = relKind{b, true}
	case tables.MANY_TO_ONE:
		first = relKind{a, true}
		second = relKind{b, false}
	case tables.MANY_TO_MANY:
		first = relKind{a, true}
		second = relKind{b, true}
	default:
		panic("unknown cardinality")
	}
	return
}
