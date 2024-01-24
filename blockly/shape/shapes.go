package shape

import (
	"sort"

	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// turn the passed types into a (json list of) shapes
func FromTypes(types []*typeinfo.TypeSet) (ret string, err error) {
	return FromTypeMap(MakeTypeMap(types))
}

func FromTypeMap(ts TypeMap) (ret string, err error) {
	w := NewShapeWriter(ts)
	ret = js.Embrace(js.Array, func(out *js.Builder) {
		var comma bool
		for _, key := range ts.Keys() {
			blockType := ts[key]
			if comma {
				out.R(js.Comma)
				comma = false
			}
			if w.WriteShape(out, blockType) {
				comma = true
				if flow, ok := blockType.(*typeinfo.Flow); ok {
					out.R(js.Comma)
					w.writeMutator(out, flow)
				}
			}
		}
	})
	return
}

type TypeMap map[string]typeinfo.T

func (m TypeMap) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func MakeTypeMap(types []*typeinfo.TypeSet) TypeMap {
	m := make(TypeMap)
	for _, g := range types {
		addTypes(m, g.Flow)
		addTypes(m, g.Slot)
		addTypes(m, g.Str)
		addTypes(m, g.Num)
	}
	return m
}

func addTypes[V typeinfo.T](m TypeMap, slice []V) {
	for _, x := range slice {
		m[x.TypeName()] = x
	}
}
