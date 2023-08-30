package qna

import (
	"git.sr.ht/~ionous/tapestry/qna/query"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

// given an id, return the (a) name defined for it by the author
func (run *Runner) getObjectName(id string) (ret string, err error) {
	if c, e := run.values.cache(func() (ret any, err error) {
		ret, err = run.query.NounName(id)
		return
	}, "objectName", id); e != nil {
		err = e
	} else {
		ret = c.(string)
	}
	return
}

// given an id, return the names defined for it by the author
func (run *Runner) getObjectNames(id string) (ret []string, err error) {
	if c, e := run.values.cache(func() (ret any, err error) {
		ret, err = run.query.NounNames(id)
		return
	}, "objectAliases", id); e != nil {
		err = e
	} else {
		ret = c.([]string)
	}
	return
}

// given an object name, return its id and kind.
func (run *Runner) getObjectInfo(name string) (ret query.NounInfo, err error) {
	if c, e := run.values.cache(func() (ret any, err error) {
		if info, e := run.query.NounInfo(name); e != nil {
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
		ret = c.(query.NounInfo)
	}
	return
}
