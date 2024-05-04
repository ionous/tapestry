package qna

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/qna/decoder"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/aspects"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"github.com/ionous/errutil"
)

func (run *Runner) GetKindByName(rawName string) (ret *rt.Kind, err error) {
	if name := inflect.Normalize(rawName); len(name) == 0 {
		err = errutil.New("no kind of empty name")
	} else if cached, e := run.getKind(name); e != nil {
		err = e
	} else {
		ret = cached
	}
	return
}

func (run *Runner) getKindOf(kn, kt string) (ret *rt.Kind, err error) {
	if ck, e := run.getKind(kn); e != nil {
		err = e
	} else if !ck.Implements(kt) {
		err = errutil.New(kn, "not a kind of", kt)
	} else {
		ret = ck
	}
	return
}

func (run *Runner) getAncestry(k string) (ret []string, err error) {
	run.ensureBaseKinds()
	if c, e := run.values.cache(func() (ret any, err error) {
		if path, e := run.query.KindOfAncestors(k); e != nil {
			err = errutil.Fmt("error while getting kind %q, %w", k, e)
		} else {
			ret = path
		}
		return
	}, meta.KindAncestry, k); e != nil {
		err = e
	} else {
		ret = c.([]string)
	}
	return
}

func (run *Runner) getKind(k string) (ret *rt.Kind, err error) {
	run.ensureBaseKinds()
	if c, e := run.values.cache(func() (ret any, err error) {
		ret, err = run.buildKind(k)
		return
	}, "kinds", k); e != nil {
		err = e
	} else {
		ret = c.(*rt.Kind)
	}
	return
}

// tbd: maybe macros and actions shouldnt have the parent  "pattern";
// it would simplify this, and re: Categorize the base shouldnt be needed anymore.
func (run *Runner) ensureBaseKinds() {
	key := makeKey("kinds", kindsOf.Kind.String())
	if _, ok := run.values.store[key]; !ok {
		for _, k := range kindsOf.DefaultKinds {
			var err error
			var kind *rt.Kind
			// note: responses have fields, even though the other base kinds dont
			if fs, e := run.getAllFields(k.String()); e != nil {
				err = errutil.Fmt("error while building kind %q, %w", k, e)
			} else {
				path := []string{k.String(), k.Parent().String()}
				kind = rt.NewKind(path, fs, nil)
			}
			key := makeKey("kinds", k.String())
			run.values.store[key] = cachedValue{kind, err}
		}
	}
	return
}

// fix? this is maybe a little odd... because when the domain changes, so will the kinds
// ( unless maybe we precache them all or change kind query to use a fixed (set of) domains
// - and record the domain into the cache; and/or build an in memory tree of kinds as a cache. )
func (run *Runner) buildKind(k string) (ret *rt.Kind, err error) {
	if path, e := run.getAncestry(k); e != nil {
		err = e
	} else if cnt := len(path); cnt < 2 {
		// should have a name and some base type
		err = errutil.Fmt("invalid kind %q", k)
	} else {
		k := path[0] // use the returned name in favor of the given name (ex. plurals)
		if fields, e := run.getAllFields(k); e != nil {
			err = errutil.Fmt("error while building kind %q, %w", k, e)
		} else {
			var traits []rt.Aspect // fix? currently only allow traits for objects. hrm.
			if objectLike := path[len(path)-1] == kindsOf.Kind.String(); objectLike {
				traits = aspects.MakeAspects(run, fields)
			}
			ret = rt.NewKind(path, fields, traits)
		}
	}
	return
}

// cached fields exclusive to a kind
func (run *Runner) getAllFields(kind string) (ret []rt.Field, err error) {
	if c, e := run.values.cache(func() (ret any, err error) {
		return run.query.FieldsOf(kind)
	}, "fields", kind); e != nil {
		err = e
	} else {
		fs := c.([]query.FieldData)
		for _, f := range fs {
			if init, e := decodeInit(run.decode, f.Affinity, f.Init); e != nil {
				err = e
			} else {
				ret = append(ret, rt.Field{Name: f.Name, Affinity: f.Affinity, Type: f.Class, Init: init})
			}
		}
	}
	return
}

// decode the passed assignment, if it exists.
func decodeInit(d decoder.Decoder, aff affine.Affinity, b []byte) (ret rt.Assignment, err error) {
	if len(b) > 0 {
		ret, err = d.DecodeAssignment(aff, b)
	}
	return
}
