package weave_test

import (
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/test/eph"
)

func dd(names ...string) []string {
	return names
}

// relation, kind, cardinality, otherKinds
func newRelation(r, k, c, o string) *eph.Relations {
	var card eph.Cardinality
	switch c {
	case tables.ONE_TO_ONE:
		card = &eph.OneOne{Kind: k, OtherKind: o}
	case tables.ONE_TO_MANY:
		card = &eph.OneMany{Kind: k, OtherKinds: o}
	case tables.MANY_TO_ONE:
		card = &eph.ManyOne{Kinds: k, OtherKind: o}
	case tables.MANY_TO_MANY:
		card = &eph.ManyMany{Kinds: k, OtherKinds: o}
	default:
		panic("unknown cardinality")
	}
	return &eph.Relations{
		Rel:         r,
		Cardinality: card,
	}
}
