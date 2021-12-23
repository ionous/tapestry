package qna

import (
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/generic"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"github.com/ionous/errutil"
)

func (run *Runner) GetKindByName(rawName string) (ret *g.Kind, err error) {
	name := lang.Underscore(rawName)
	if cached, e := run.getCachedKind(name); e != nil {
		err = e
	} else {
		ret = cached.kind
	}
	return
}
func (run *Runner) getCachedKind(k string) (ret cachedKind, err error) {
	if c, e := run.values.cache(func() (ret interface{}, err error) {
		ret, err = run.buildKind(k)
		return
	}, "kinds", k); e != nil {
		err = e
	} else {
		ret = c.(cachedKind)
	}
	return
}
func (run *Runner) buildKind(k string) (ret cachedKind, err error) {
	if path, e := run.qdb.KindOfAncestors(k); e != nil {
		err = errutil.Fmt("error while getting kind %q, %w", k, e)
	} else if fields, init, e := run.getFields(k); e != nil {
		err = errutil.Fmt("error while building kind %q, %w", k, e)
	} else {
		// fix? this is maybe a little odd... because when the domain changes, so will the kinds
		// ( unless maybe we precache them all or change kind query to use a fixed (set of) domains
		//   and record the domain into the cache; and/or build an in memory tree of kinds as a cache. )
		kinds := generic.Kinds(run)
		kind := g.NewKind(kinds, k, path, fields)
		ret = cachedKind{kind, init}
	}
	return
}

type cachedKind struct {
	kind *g.Kind
	init []rt.Assignment // FIX FIX FIX -- when is init actually getting used? should be with patterns... somewhere...
}

// for kind initialization
func (run *Runner) getFields(kind string) (ret []g.Field, retInit []rt.Assignment, err error) {
	if fs, e := run.qdb.FieldsOf(kind); e != nil {
		err = e
	} else {
		var init []rt.Assignment
		out := make([]g.Field, len(fs))
		for i, f := range fs {
			out[i] = g.Field{f.Name, f.Affinity, f.Class}
			//
			if prog := f.Init; len(prog) > 0 {
				if val, e := decodeAssignment(f.Affinity, prog, run.signatures); e != nil {
					err = errutil.New("error while decoding", f.Name, e)
					break
				} else {
					// retInit and ret need to be the same length
					// its possible we have had several fields without init though:
					if retInit == nil {
						cnt := len(ret) // make with a current size one less than cnt so we can append.
						init = make([]rt.Assignment, cnt-1, cnt)
					}
					init = append(init, val)
				}
			}
		}
		if err == nil {
			ret, retInit = out, init
		}
	}
	return
}
