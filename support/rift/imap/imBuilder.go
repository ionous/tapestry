package imap

import "git.sr.ht/~ionous/tapestry/support/rift/maps"

// return a builder can generate ItemMap
func Build(reserve int) maps.Builder {
	return mapBuilder{values: make(ItemMap, reserve)}
}

type mapBuilder struct {
	values ItemMap
}

func (b mapBuilder) Add(key string, val any) maps.Builder {
	b.values = append(b.values, MapItem{key, val})
	return b
}

// returns ItemMap
func (b mapBuilder) Map() any {
	return b.values
}
