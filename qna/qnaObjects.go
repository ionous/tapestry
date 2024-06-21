package qna

import (
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
)

// given an id, return the (a) name defined for it by the author
func (run *Runner) getObjectName(id string) (ret string, err error) {
	key := query.MakeKey("objectName", id, "")
	if c, e := run.constVals.Ensure(key, func() (ret any, err error) {
		ret, err = run.query.NounName(id)
		return
	}); e != nil {
		err = e
	} else {
		ret = c.(string)
	}
	return
}

// given an id, return the names defined for it by the author
func (run *Runner) getObjectNames(id string) (ret []string, err error) {
	key := query.MakeKey("objectAliases", id, "")
	if c, e := run.constVals.Ensure(key, func() (ret any, err error) {
		ret, err = run.query.NounNames(id)
		return
	}); e != nil {
		err = e
	} else {
		ret = c.([]string)
	}
	return
}

// given an object name, return its id and kind.
func (run *Runner) getObjectInfo(name string) (ret query.NounInfo, err error) {
	key := query.MakeKey("objectInfo", name, "")
	if c, e := run.constVals.Ensure(key, func() (ret any, err error) {
		if info, e := run.query.NounInfo(name); e != nil {
			err = e
		} else if !info.IsValid() {
			err = rt.UnknownObject(name)
		} else {
			ret = info
		}
		return
	}); e != nil {
		err = e
	} else {
		ret = c.(query.NounInfo)
	}
	return
}
