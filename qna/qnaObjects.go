package qna

import (
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

// given an id, return the (a) name defined for it by the author
func (run *Runner) getObjectName(id string) (ret string, err error) {
	if c, e := run.values.cache(func() (ret interface{}, err error) {
		ret, err = run.qdb.NounName(id)
		return
	}, "objectName", id); e != nil {
		err = e
	} else {
		ret = c.(string)
	}
	return
}

// given an object name, return its id and kind.
func (run *Runner) getObjectInfo(name string) (ret qdb.NounInfo, err error) {
	if c, e := run.values.cache(func() (ret interface{}, err error) {
		if info, e := run.qdb.NounInfo(name); e != nil {
			err = e
		} else if !info.IsValid() {
			err = g.UnknownObject(name)
		} else {
			ret = info
		}
		return
	}, "objectInfo", name); e != nil {
		err = e
	} else {
		ret = c.(qdb.NounInfo)
	}
	return
}
