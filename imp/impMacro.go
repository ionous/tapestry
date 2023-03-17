package imp

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/pattern"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// duplicates cachedKind.
type macroKind struct {
	*g.Kind
	init []assign.Assignment
	do   []rt.Execute
}

func (k macroKind) initializeRecord(run rt.Runtime, rec *g.Record) (err error) {
	for fieldIndex, init := range k.init {
		if init != nil {
			ft := k.Field(fieldIndex)
			if src, e := safe.GetAssignment(run, init); e != nil {
				err = errutil.New("error determining local", k.Name(), ft.Name, e)
				break
			} else if val, e := safe.AutoConvert(run, ft, src); e != nil {
				err = e
			} else if e := rec.SetIndexedField(fieldIndex, val); e != nil {
				err = errutil.New("error setting local", k.Name(), ft.Name, e)
				break
			}
		}
	}
	return
}

type macroReg map[string]macroKind

func (m macroReg) GetKindByName(n string) (ret *g.Kind, err error) {
	if a, ok := m[n]; !ok {
		err = errutil.New("no such kind")
	} else {
		ret = a.Kind
	}
	return
}

// ugh. see register macro notes.
func (k *Importer) Call(rec *g.Record, expectedReturn affine.Affinity) (ret g.Value, err error) {
	kind := rec.Kind()
	if macro, ok := k.macros[kind.Name()]; !ok {
		err = errutil.New("unknown macro", kind.Name())
	} else if res, e := pattern.NewResults(k, rec, expectedReturn); e != nil {
		err = e
	} else if e := macro.initializeRecord(k, rec); e != nil {
		err = e
	} else {
		oldScope := k.Stack.ReplaceScope(res)
		if e := macro.initializeRecord(k, rec); e != nil {
			err = e
		} else if e := safe.RunAll(k, macro.do); e != nil {
			err = e
		} else if v, e := res.GetResult(); e != nil {
			err = e
		} else {
			// warning: in order to generate appropriate defaults ( ex. a record of the right type )
			// while still informing the caller of lack of pattern decision in a concise manner
			// can return both a valid value and an error
			ret = v
			if !res.ComputedResult() {
				err = rt.NoResult{}
			}
		}
		k.Stack.ReplaceScope(oldScope)
	}
	return
}

// EphMacro - hijacks pattern registration for use with macros
type EphMacro struct {
	eph.EphPatterns
	MacroStatements []rt.Execute
}

// oh, the tangled webs we weave.
// normally we distill fields into ephParams,
// and then build up kinds, finally storing the results in the db.
// fix: ideally would flush to the db after each domain ( or phase ) so that the weave can see it.
// right now, only the play time's runtime reads from the db.
func (k *Importer) registerMacro(op *EphMacro) (err error) {
	// check not already registered.
	if name := op.PatternName; len(name) == 0 {
		err = errutil.Fmt("no macro name specified")
	} else if _, ok := k.macros[name]; ok {
		err = errutil.Fmt("macro %q already registered", name)
	} else {
		cnt := 1 + len(op.Params) + len(op.Locals)
		init := make([]assign.Assignment, cnt)
		fields := make([]g.Field, cnt)
		addParam := func(p eph.EphParams) (err error) {
			if u, e := p.Unify(op.PatternName); e != nil {
				err = e
			} else {
				init = append(init, u.Initially)
				fields = append(fields, g.Field{
					Name:     u.Name,
					Affinity: u.Affinity,
					Type:     u.Type,
				})
			}
			return
		}
		addParams := func(ps []eph.EphParams) (err error) {
			for _, el := range ps {
				if e := addParam(el); e != nil {
					err = e
					break
				}
			}
			return
		}
		if e := addParam(*op.Result); e != nil {
			err = e
		} else if e := addParams(op.Params); e != nil {
			err = e
		} else if e := addParams(op.Locals); e != nil {
			err = e
		} else {
			k.macros[name] = macroKind{
				Kind: g.NewKind(k.macros, name, nil, fields),
				init: init,
				do:   op.MacroStatements,
			}
		}
	}
	return
}
