// package orderedmap
// implements rift maps interface for ian coleman's ordered map implementation
// https://github.com/iancoleman/orderedmap
package orderedmap

import "git.sr.ht/~ionous/tapestry/support/rift/maps"

// return a builder which generates a ItemMap
func Build(reserve int) maps.Builder {
	// orderedmap exposes New() but we shouldnt need the double dereference
	// alt: the compiler might be smart enough to handle *New() as a non allocating copy
	return sliceBuilder{values: OrderedMap{
		values:     make(map[string]any),
		escapeHTML: true,
	}}
}

type sliceBuilder struct {
	values OrderedMap
}

func (b sliceBuilder) Add(key string, val any) maps.Builder {
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
