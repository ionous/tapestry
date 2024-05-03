package qna

import (
	"database/sql"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"github.com/ionous/errutil"
)

type counters map[string]int

func (c *counters) getCounter(name string) (ret rt.Value, err error) {
	// alt: a global $counters object with fields?
	i := (*c)[name]
	ret = rt.IntOf(i)
	return
}

func (c *counters) setCounter(name string, val rt.Value) (err error) {
	if aff := val.Affinity(); aff != affine.Number {
		err = errutil.Fmt("counter %q expected a number got %s", name, aff)
	} else {
		(*c)[name] = val.Int()
	}
	return
}

func (c counters) writeCounters(db *sql.DB) error {
	return writeValues(db, func(q *sql.Stmt) (err error) {
		for k, v := range c {
			if _, e := q.Exec(meta.Counter, k, v); e != nil {
				err = e
				break
			}
		}
		return
	})
}
