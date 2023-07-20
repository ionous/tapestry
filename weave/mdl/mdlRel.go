package mdl

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/tables"
)

func makeCard(amany, bmany bool) (ret string) {
	switch {
	case !amany && !bmany:
		ret = tables.ONE_TO_ONE
	case !amany && bmany:
		ret = tables.ONE_TO_MANY
	case amany && !bmany:
		ret = tables.MANY_TO_ONE
	case amany && bmany:
		ret = tables.MANY_TO_MANY
	}
	return
}

type relInfo struct {
	one, other  string
	cardinality string
}

func (rel *relInfo) String() string {
	return fmt.Sprint(rel.other, rel.other, rel.cardinality)
}

func (rel *relInfo) makeRel() (first, second relKind) {
	switch rel.cardinality {
	case tables.ONE_TO_ONE:
		first = relKind{rel.one, false}
		second = relKind{rel.other, false}
	case tables.ONE_TO_MANY:
		first = relKind{rel.one, false}
		second = relKind{rel.other, true}
	case tables.MANY_TO_ONE:
		first = relKind{rel.one, true}
		second = relKind{rel.other, false}
	case tables.MANY_TO_MANY:
		first = relKind{rel.one, true}
		second = relKind{rel.other, true}
	default:
		panic("unknown cardinality")
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
		ret = "other kinds"
	} else {
		ret = "other kind"
	}
	return
}
