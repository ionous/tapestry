package qna

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/aspects"
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
	run.ensureBaseKinds()
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

// tbd: maybe macros and actions shouldnt have the parent  "pattern";
// it would simplify this, and re: Categorize the base shouldnt be needed anymore.
func (run *Runner) ensureBaseKinds() {
	key := makeKey("kinds", kindsOf.Kind.String())
	if _, ok := run.values[key]; !ok {
		for _, k := range kindsOf.DefaultKinds {
			var err error
			var kind *g.Kind
			// note: responses have fields, even though the other base kinds dont
			if fs, e := run.getExclusiveFields(k.String()); e != nil {
				err = errutil.Fmt("error while building kind %q, %w", k, e)
			} else {
				var parent *g.Kind
				if p := k.Parent().String(); len(p) > 0 {
					parent, err = run.getKind(p)
				}
				if err == nil {
					kind = g.NewKind(k.String(), parent, fs)
				}
			}
			key := makeKey("kinds", k.String())
			run.values[key] = cachedValue{kind, err}
		}
	}
	return
}

// gofmt will add extra lines if the "fix" comment below is here. :'(
func (run *Runner) buildKind(k string) (ret *g.Kind, err error) {
	// fix? this is maybe a little odd... because when the domain changes, so will the kinds
	// ( unless maybe we precache them all or change kind query to use a fixed (set of) domains
	//	and record the domain into the cache; and/or build an in memory tree of kinds as a cache. )
	if path, e := run.query.KindOfAncestors(k); e != nil {
		err = errutil.Fmt("error while getting kind %q, %w", k, e)
	} else if cnt := len(path); cnt == 0 {
		err = errutil.Fmt("invalid kind %q", k)
	} else if parent, e := run.getKind(path[cnt-1]); e != nil {
		err = e
	} else if fields, e := run.getExclusiveFields(k); e != nil {
		err = errutil.Fmt("error while building kind %q, %w", k, e)
	} else {
		// we never actually use the field values of the kind:
		// instead we pull individual defaults from the db.
		// because of that, weave generates reasonable default values for kinds with traits

		if objectLike := path[0] == kindsOf.Kind.String(); objectLike {
			traits := aspects.MakeAspects(run, fields)
			ret = g.NewKindWithTraits(k, parent, fields, traits)
		} else {
			ret = g.NewKind(k, parent, fields)
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
