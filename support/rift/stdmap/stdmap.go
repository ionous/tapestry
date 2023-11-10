package stdmap

import "git.sr.ht/~ionous/tapestry/support/rift/maps"

type StdMap map[string]any

func Build(reserve int) maps.Builder {
	return make(StdMap)
}

func (m StdMap) Add(key string, val any) maps.Builder {
	m[key] = val
	return m
}

func (m StdMap) Map() any {
	return (map[string]any)(m)
}
