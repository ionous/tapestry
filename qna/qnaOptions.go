package qna

import (
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"github.com/ionous/errutil"
)

type qnaOptions map[string]g.Value

func (m qnaOptions) SetOption(name string, v g.Value) (err error) {
	if was, ok := m[name]; !ok {
		err = errutil.New("no such option", name)
	} else if a, va := was.Affinity(), v.Affinity(); a != va {
		err = errutil.New("option", name, "requires", a, "had", va, v)
	} else {
		m[name] = v
	}
	return
}

func (m qnaOptions) Option(name string) (ret g.Value, err error) {
	if was, ok := m[name]; !ok {
		err = errutil.New("no such option", name)
	} else {
		ret = was
	}
	return
}
