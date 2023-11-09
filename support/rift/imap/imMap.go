package imap

// used to preserve the specified order of keys
type ItemMap []MapItem

type MapItem struct {
	Key   string
	Value any
}

// returns a valid pointer into the slice, or nil if not found
func (m ItemMap) Find(k string) (ret *MapItem, okay bool) {
	if at := m.FindIndex(k); at >= 0 {
		ret, okay = &(m[at]), true
	}
	return
}

// returns the index of the item or -1 if not found
func (m ItemMap) FindIndex(k string) (ret int) {
	ret = -1 // provisionally
	for i, kv := range m {
		if kv.Key == k {
			ret = i
			break
		}
	}
	return
}
