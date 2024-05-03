package testutil

import (
	"sort"

	"git.sr.ht/~ionous/tapestry/rt"
)

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

func (or *Objects) Names() (ret []string) {
	for n := range *or {
		ret = append(ret, n)
	}
	sort.Strings(ret)
	return
}
