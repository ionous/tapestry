package qna

import (
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"github.com/ionous/errutil"
)

// uses its own cache to preserve values across domain changes
type qnaOptions map[string]g.Value

// change an existing option.
func (m qnaOptions) setOption(name string, v g.Value) (err error) {
	// see meta.Options: new options cannot be added dynamically at runtime
	if was, ok := m[name]; !ok {
		err = errutil.New("no such option", name)
	} else if a, va := was.Affinity(), v.Affinity(); a != va {
		err = errutil.New("option", name, "requires", a, "had", va, v)
	} else {
		m[name] = v
	}
	return
}

// return an existing option.
func (m qnaOptions) option(name string) (ret g.Value, err error) {
	if was, ok := m[name]; !ok {
		err = errutil.New("no such option", name)
	} else {
		ret = was
	}
	return
}
