package idlcmd

import (
	"git.sr.ht/~ionous/tapestry/tables/idl"
	"github.com/ionous/errutil"
)

// a terrible way to optimize database writes
type cachedWrite struct {
	tgt  string
	args []interface{}
}
type writeCache struct {
	// have to write all of the types first, so that we can ask for their keys
	// since we're caching things in memory first anyway...
	// split them into two groups.
	ops  []cachedWrite
	rest []cachedWrite
}

func (q *writeCache) Write(tgt string, args ...interface{}) (err error) {
	data := cachedWrite{tgt, args}
	if tgt == idl.Op {
		q.ops = append(q.ops, data)
	} else {
		q.rest = append(q.rest, data)
	}
	return
}

func (q *writeCache) Flush(w writer) (err error) {
	if e := writeSlice(w, q.ops); e != nil {
		err = e
	} else if e := writeSlice(w, q.rest); e != nil {
		err = e
	}
	return
}

func writeSlice(w writer, els []cachedWrite) (err error) {
	for _, el := range els {
		if e := w.Write(el.tgt, el.args...); e != nil {
			err = errutil.New("couldnt write to", el.tgt, e)
			break
		}
	}
	return
}
