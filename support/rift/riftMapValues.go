package rift

// used to preserve the specified order of keys
type MapValues []MapValue

type MapValue struct {
	Key   string
	Value any
}

func (m *MapValues) Add(key string, val any) {
	(*m) = append(*m, MapValue{key, val})
}

// returns a valid pointer into the slice, or nil if not found
func (m MapValues) Find(k string) (ret *MapValue, okay bool) {
	if at := m.FindIndex(k); at >= 0 {
		ret, okay = &(m[at]), true
	}
	return
}

// returns the index of the item or -1 if not found
func (m MapValues) FindIndex(k string) (ret int) {
	ret = -1 // provisionally
	for i, kv := range m {
		if kv.Key == k {
			ret = i
			break
		}
	}
	return
}
