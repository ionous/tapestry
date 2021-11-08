package testutil

import (
	"sort"

	g "git.sr.ht/~ionous/iffy/rt/generic"
)

type Objects map[string]*g.Record

func (or *Objects) AddObjects(kind *g.Kind, names ...string) {
	if *or == nil {
		*or = make(Objects)
	}
	for _, name := range names {
		(*or)[name] = kind.NewRecord()
	}
}

func (or *Objects) Names() (ret []string) {
	for n := range *or {
		ret = append(ret, n)
	}
	sort.Strings(ret)
	return
}
