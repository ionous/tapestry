package qna

import (
	"git.sr.ht/~ionous/iffy/qna/pdb"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

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

func (run *Runner) getObjectByName(name string) (ret *qnaObject, err error) {
	// note: if were able to get the object kind, then the object is in scope.
	if ok, e := run.getObjectInfo(name); e != nil {
		err = e
	} else if c, e := run.values.cache(func() (ret interface{}, err error) {
		if k, e := run.GetKindByName(ok.Kind); e != nil {
			err = e
		} else {
			ret = &qnaObject{run: run, domain: ok.Domain, name: ok.Id, kind: k}
		}
		return
	}, "objectByName", name); e != nil {
		err = e
	} else {
		ret = c.(*qnaObject)
	}
	return
}

// given an object name, return its id and kind.
func (run *Runner) getObjectInfo(name string) (ret pdb.NounInfo, err error) {
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
		ret = c.(pdb.NounInfo)
	}
	return
}
