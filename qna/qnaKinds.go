package qna

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"github.com/ionous/errutil"
)

func (run *Runner) GetKindByName(rawName string) (ret *rt.Kind, err error) {
	if name := inflect.Normalize(rawName); len(name) == 0 {
		err = errutil.New("no kind of empty name")
	} else {
		ret, err = run.getKind(name)
	}
	return
}

func (run *Runner) getKind(k string) (*rt.Kind, error) {
	return run.query.GetKindByName(k)
}

func (run *Runner) getKindOf(kn, kt string) (ret *rt.Kind, err error) {
	if ck, e := run.getKind(kn); e != nil {
		err = e
	} else if !ck.Implements(kt) {
		err = errutil.New(kn, "not a kind of", kt)
	} else {
		ret = ck
	}
	return
}

func (run *Runner) getAncestry(k string) (ret []string, err error) {
	if path, e := run.query.KindOfAncestors(k); e != nil {
		err = errutil.Fmt("error while getting kind %q, %w", k, e)
	} else {
		ret = path
	}
	return
}
