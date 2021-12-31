package qna

import (
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/generic"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"github.com/ionous/errutil"
)

func (run *Runner) GetKindByName(rawName string) (ret *g.Kind, err error) {
	if name := lang.Underscore(rawName); len(name) == 0 {
		err = errutil.New("no kind of empty name")
	} else if cached, e := run.getKind(name); e != nil {
		err = e
	} else {
		ret = cached.Kind
	}
	return
}

func (run *Runner) getKindOf(kn, kt string) (ret cachedKind, err error) {
	if ck, e := run.getKind(kn); e != nil {
		err = e
	} else if !ck.Implements(kt) {
		err = errutil.New(kn, "not a kind of", kt)
	} else {
		ret = ck
	}
	return
}

func (run *Runner) getKind(k string) (ret cachedKind, err error) {
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
	} else if fs, e := run.getFieldSet(k, path); e != nil {
		err = errutil.Fmt("error while building kind %q, %w", k, e)
	} else {
		// fix? this is maybe a little odd... because when the domain changes, so will the kinds
		// ( unless maybe we precache them all or change kind query to use a fixed (set of) domains
		//   and record the domain into the cache; and/or build an in memory tree of kinds as a cache. )
		kinds := generic.Kinds(run)
		kind := g.NewKind(kinds, k, path, fs.fields)
		ret = cachedKind{kind, fs.init}
	}
	return
}

type cachedKind struct {
	*g.Kind
	init []rt.Assignment
}

type fieldSet struct {
	fields []g.Field       // the kind's own fields are first, the root fields are last
	init   []rt.Assignment // fix? move this into g.Kinds ( probably as a callback Init()g.Value to avoid dependency on package rt )
}

// get the fields and initialization settings of a hierarchy
// fields for kinds are "flattened" so all of the info for a hierarchy is duplicated in each kind
// tbd if that makes sense or not, but its how things are for now.
func (run *Runner) getFieldSet(kind string, path []string) (ret fieldSet, err error) {
	var out fieldSet
	for next, i := kind, len(path); ; {
		if fs, e := run.getFields(next); e != nil {
			err = e
			break
		} else {
			if cnt := len(fs.init); cnt > 0 {
				if out.init == nil {
					out.init = make([]rt.Assignment, len(out.fields))
				}
				out.init = append(out.init, fs.init...)
			}
			out.fields = append(out.fields, fs.fields...)
		}
		// do while
		if i--; i >= 0 {
			next = path[i] // root is at the start, and we visit it last.
		} else {
			break
		}
	}
	if err == nil {
		ret = out
	}
	return
}

// get the fields and initialization settings of a single kind
func (run *Runner) getFields(kind string) (ret fieldSet, err error) {
	if fs, e := run.qdb.FieldsOf(kind); e != nil {
		err = e
	} else {
		fields := make([]g.Field, len(fs))
		var init []rt.Assignment // init is so often empty, only allocate on demand.
		for i, f := range fs {
			fields[i] = g.Field{f.Name, f.Affinity, f.Class}
			if prog := f.Init; len(prog) > 0 {
				if val, e := decodeAssignment(f.Affinity, prog, run.signatures); e != nil {
					err = errutil.New("error while decoding", f.Name, e)
					break
				} else {
					if init == nil {
						init = make([]rt.Assignment, len(fs))
					}
					init[i] = val
				}
			}
		}
		if err == nil {
			ret.fields, ret.init = fields, init
		}
	}
	return
}
