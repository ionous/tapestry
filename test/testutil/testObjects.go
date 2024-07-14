package testutil

import (
	"sort"

	"git.sr.ht/~ionous/tapestry/rt"
)

// nouns in the test runtime are represented by records
type Objects map[string]*rt.Record

func (or *Objects) AddObjects(kind *rt.Kind, names ...string) int {
	if *or == nil {
		*or = make(Objects)
	}
	for _, name := range names {
		(*or)[name] = rt.NewRecord(kind)
	}
	return len(names)
}

func (or *Objects) AddObject(kind *rt.Kind, name string) (ret *rt.Record) {
	if *or == nil {
		*or = make(Objects)
	} else if _, exists := (*or)[name]; exists {
		panic("can't add the same name twice")
	}
	obj := rt.NewRecord(kind)
	(*or)[name] = obj
	return obj
}

func (or *Objects) Names() (ret []string) {
	for n := range *or {
		ret = append(ret, n)
	}
	sort.Strings(ret)
	return
}
