package rift

// used to preserve the specified order of keys
type MapValues []MapValue

type MapValue struct {
	Key   string
	Value any
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

// type MappingTarget struct {
// 	Values MapValues // tbd: possibly a pointer to the slice?

// 	}

// type KeyValueTarget struct {
// 	Values MapValues // tbd: possibly a pointer to the slice?
// 	Key    string
// }

// var _ CollectionTarget = (*Document)(nil)
// var _ CollectionTarget = (*Sequence)(nil)
// var _ CollectionTarget = (*KeyValue)(nil)

// // func (d *DocumentTarget) PopValue() any {
// 	return d.Value
// }

// func (d *SequenceTarget) PopValue() any {
// 	return d.Values
// }

// func (d *KeyValueTarget) PopValue() any {
// 	return d.Values
// // }

// func (d *KeyValueTarget) WriteValue(val any) (err error) {
// 	if at := d.Values.FindIndex(d.Key); at < 0 {
// 		err = errutil.Fmt("key %q already has a value", d.Key)
// 	} else {
// 		d.Values = append(d.Values, MapValue{Key: d.Key, Value: val})
// 	}
// 	return
// }
