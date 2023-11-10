// package orderedmap
// implements rift maps interface for ian coleman's ordered map implementation
// https://github.com/iancoleman/orderedmap
package orderedmap

import "git.sr.ht/~ionous/tapestry/support/rift/maps"

// return a builder which generates a ItemMap
func Build(reserve bool) maps.Builder {
	var keys []string
	if reserve {
		keys = make([]string, 1)
	}
	// orderedmap exposes New() which returns a pointer; we dont need the extra dereference
	// alt: the compiler might be smart enough to handle *New() as a non allocating copy
	// ( and values could init'd after creation )
	return sliceBuilder{values: OrderedMap{
		values:     make(map[string]any),
		escapeHTML: true,
		keys:       keys,
	}}
}

type sliceBuilder struct {
	values OrderedMap
}

func (b sliceBuilder) Add(key string, val any) maps.Builder {
	if len(key) == 0 { // there should be only one blank key; at the start
		if _, exists := b.values.Get(key); !exists {
			// could adjust the slice. but the program should know better.
			panic("map doesn't have space for a blank key")
		}
	}
	b.values.Set(key, val)
	return b
}

// returns an OrderedMap
func (b sliceBuilder) Map() any {
	return b.values
}

// shortcut to access the underlying ordered keys
func (b sliceBuilder) Keys() []string {
	return b.values.Keys()
}

// shortcut to access the underlying unordered map
func (b sliceBuilder) Values() map[string]any {
	return b.values.Values()
}
