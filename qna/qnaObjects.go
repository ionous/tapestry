package qna

import (
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
)

// given an id, return the (a) name defined for it by the author
func (run *Runner) getObjectName(id string) (ret string, err error) {
	return run.query.NounName(id)
}

// given an id, return the names defined for it by the author
func (run *Runner) getObjectNames(id string) (ret []string, err error) {
	return run.query.NounNames(id)
}

// given an object name, return its id and kind.
func (run *Runner) getObjectInfo(name string) (ret query.NounInfo, err error) {
	if info, e := run.query.NounInfo(name); e != nil {
		err = e
	} else if !info.IsValid() {
		err = rt.UnknownObject(name)
	} else {
		ret = info
	}
	return
}
