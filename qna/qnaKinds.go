package qna

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/generic"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"github.com/ionous/errutil"
)

func (run *Runner) GetKindByName(rawName string) (ret *g.Kind, err error) {
	if name := lang.Normalize(rawName); len(name) == 0 {
		err = errutil.New("no kind of empty name")
	} else if cached, e := run.getKind(name); e != nil {
		err = e
	} else {
		ret = cached
	}
	return
}

func (run *Runner) getKindOf(kn, kt string) (ret *g.Kind, err error) {
	if ck, e := run.getKind(kn); e != nil {
		err = e
	} else if !ck.Implements(kt) {
		err = errutil.New(kn, "not a kind of", kt)
	} else {
		ret = ck
	}
	return
}

func (run *Runner) getKind(k string) (ret *g.Kind, err error) {
	if c, e := run.values.cache(func() (ret any, err error) {
		ret, err = run.buildKind(k)
		return
	}, "kinds", k); e != nil {
		err = e
	} else {
		ret = c.(*g.Kind)
	}
	return
}

func (run *Runner) buildKind(k string) (ret *g.Kind, err error) {
	if path, e := run.query.KindOfAncestors(k); e != nil {
		err = errutil.Fmt("error while getting kind %q, %w", k, e)
	} else if fs, e := run.getFields(k, path); e != nil {
		err = errutil.Fmt("error while building kind %q, %w", k, e)
	} else {
		// fix? kinds and fields can be zero for both empty kinds and non-existent kinds
		// this tries to figure out which is which to properly report errors
		// ideally, this could be done at query time ( they should be the lowest N rows )
		var okay bool
		if len(path) > 0 || len(fs) > 0 {
			okay = true
		} else {
			for _, x := range kindsOf.DefaultKinds {
				if k == x.String() {
					okay = true
					break
				}
			}
		}
		if !okay {
			err = errutil.Fmt("unknown kind %q", k)
		} else {
			// fix? this is maybe a little odd... because when the domain changes, so will the kinds
			// ( unless maybe we precache them all or change kind query to use a fixed (set of) domains
			//   and record the domain into the cache; and/or build an in memory tree of kinds as a cache. )
			kinds := generic.Kinds(run)
			kind := g.NewKind(kinds, k, path, fs)
			ret = kind
		}
	}
	return
}

// fields for a full hierarchy
// fields for kinds are "flattened" so all of the info for a hierarchy is duplicated in each kind
func (run *Runner) getFields(kind string, path []string) (ret []g.Field, err error) {
	var out []g.Field
	for _, kind := range path {
		if fs, e := run.getExclusiveFields(kind); e != nil {
			err = e
			break
		} else {
			out = append(out, fs...)
		}
	}
	if err == nil {
		// fix? would probably make more sense if kindOfAncestors included the kind
		if fs, e := run.getExclusiveFields(kind); e != nil {
			err = e
		} else {
			ret = append(out, fs...)
		}
	}
	return
}

// cached fields exclusive to a kind
func (run *Runner) getExclusiveFields(kind string) (ret []g.Field, err error) {
	if c, e := run.values.cache(func() (ret any, err error) {
		return run.query.FieldsOf(kind)
	}, "fields", kind); e != nil {
		err = e
	} else {
		fs := c.([]query.FieldData)
		for _, f := range fs {
			ret = append(ret, g.Field{Name: f.Name, Affinity: f.Affinity, Type: f.Class})
		}
	}
	return
}

// cached initialization exclusive to a kind
func (run *Runner) getKindValues(kind *g.Kind) (ret []kindValue, err error) {
	name := kind.Name()
	if c, e := run.values.cache(func() (ret any, err error) {
		if kv, e := run.query.KindValues(name); e != nil {
			err = e
		} else {
			prev := -1
			ks := make([]kindValue, len(kv))
			for i, el := range kv {
				if len(el.Path) > 0 {
					err = errutil.New("unexpected dot in field", name, el.Field)
				} else if f, at := findNextField(kind, el.Field, prev); at < 0 {
					err = errutil.New("unknown field", name, el.Field)
					break
				} else if v, e := run.decode.DecodeAssignment(f.Affinity, el.Value); e != nil {
					err = errutil.New("error decoding field", name, el.Field, e)
					break
				} else {
					ks[i] = kindValue{ /*f: f, */ i: at, val: v}
				}
			}
			ret = ks
		}
		return
	}, "init", name); e != nil {
		err = e
	} else {
		ret, _ = c.([]kindValue) // can also be nil
	}
	return
}

func findNextField(k *g.Kind, name string, prev int) (ret g.Field, next int) {
	next = -1 // provisionally
	for i, cnt := prev+1, k.NumField(); i < cnt; i++ {
		if f := k.Field(i); f.Name == name {
			ret, next = f, i
			break
		}
	}
	return
}

type kindValue struct {
	i   int // index in kind
	val assign.Assignment
}

// only inits if the field is unset
func (kv *kindValue) initValue(run rt.Runtime, rec *g.Record) (err error) {
	if !rec.HasValue(kv.i) {
		if val, e := kv.val.GetAssignedValue(run); e != nil {
			err = e
		} else if e := rec.SetIndexedField(kv.i, val); e != nil {
			err = e
		}
	}
	return
}

// assumes the record in in scope so it can read from its own values when needed
func initRecord(run *Runner, rec *g.Record) (err error) {
	if vs, e := run.getKindValues(rec.Kind()); e != nil {
		err = e
	} else {
		for _, kv := range vs {
			if e := kv.initValue(run, rec); e != nil {
				err = e
				break
			}
		}
	}
	return
}
