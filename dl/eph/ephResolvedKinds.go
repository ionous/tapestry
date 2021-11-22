package eph

import "sort"

// a list of the assembled kinds ( built from Kinds )
type ResolvedKinds []ResolvedKind

// the first element is the kind, the last is the root, ancestors in an appropriate order in-between.
type ResolvedKind struct {
	kind  string
	kinds []string
}

func SortKinds(ks ResolvedKinds) {
	sort.Slice(ks, func(i, j int) (less bool) {
		k1, k2 := ks[i], ks[j]
		if diff := len(k1.kinds) - len(k2.kinds); diff < 0 {
			less = true
		} else if diff == 0 {
			less = k1.kind < k2.kind
		}
		return
	})
}
