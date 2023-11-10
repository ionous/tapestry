package imap

import "git.sr.ht/~ionous/tapestry/support/rift/maps"

// return a builder can generate ItemMap
func Build(reserve bool) maps.Builder {
	var cnt int
	if reserve {
		cnt = 1
	}
	return mapBuilder{values: make(ItemMap, cnt)}
}

type mapBuilder struct {
	values ItemMap
}

// panic if adding the blank key but no space for a blank key was reserved.
func (b mapBuilder) Add(key string, val any) maps.Builder {
	if len(key) == 0 { // there should be only one blank key; at the start
		if len(b.values) == 0 || len(b.values[0].Key) != 0 {
			// could adjust the slice. but the program should know better.
			panic("map doesn't have space for a blank key")
		}
		b.values[0] = MapItem{Value: val}
	}
	b.values = append(b.values, MapItem{key, val})
	return b
}

// returns ItemMap
func (b mapBuilder) Map() any {
	return b.values
}
