package weave

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
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

func (k *Catalog) GetKindByName(n string) (ret *g.Kind, err error) {
	if a, ok := k.macros[n]; !ok {
		err = errutil.New("no such kind", n)
	} else {
		ret = a.Kind
	}
	return
}

// ugh. see register macro notes.
func (k *Catalog) Call(rec *g.Record, expectedReturn affine.Affinity) (ret g.Value, err error) {
	// kind := rec.Kind()
	// if macro, ok := k.macros[kind.Name()]; !ok {
	// 	err = errutil.New("unknown macro", kind.Name())
	// } else if res, e := pattern.NewMacroResults(k, rec, expectedReturn); e != nil {
	// 	err = e
	// } else if e := macro.initializeRecord(k, rec); e != nil {
	// 	err = e
	// } else {
	// 	oldScope := k.Stack.ReplaceScope(res)
	// 	if e := macro.initializeRecord(k, rec); e != nil {
	// 		err = e
	// 	} else if e := safe.RunAll(k, macro.do); e != nil {
	// 		err = e
	// 	} else if !res.ComputedResult() && expectedReturn != affine.None {
	// 		err = errutil.Fmt("%w calling %s pattern %q", rt.NoResult, aff, rec.Kind().Name())
	// 	} else if v, e := res.GetResult(); e != nil {
	// 		err = e
	// 	} else {
	// 		ret = v
	// 	}
	// 	k.Stack.ReplaceScope(oldScope)
	// }
	return
}

// oh, the tangled webs we weave.
// normally we distill fields into ephParams,
// and then build up kinds, finally storing the results in the db.
// fix: ideally would flush to the db after each domain ( or phase ) so that the weave can see it.
// right now, only the play time's runtime reads from the db.
// func (k *Catalog) registerMacro(op *EphMacro) (err error) {
// 	// check not already registered.
// 	if name := op.PatternName; len(name) == 0 {
// 		err = errutil.Fmt("no macro name specified")
// 	} else if _, ok := k.macros[name]; ok {
// 		err = errutil.Fmt("macro %q already registered", name)
// 	} else {
// 		cnt := 1 + len(op.Params) + len(op.Locals)
// 		init := make([]assign.Assignment, 0, cnt)
// 		fields := make([]g.Field, 0, cnt)
// 		addParam := func(p eph.Params) (err error) {
// 			if len(p.Name) > 0 { // for now, silent skip nothing fields
// 				if u, e := p.Unify(op.PatternName); e != nil {
// 					err = e
// 				} else {
// 					init = append(init, u.Initially)
// 					fields = append(fields, g.Field{
// 						Name:     u.Name,
// 						Affinity: u.Affinity,
// 						Type:     u.Type,
// 					})
// 				}
// 			}
// 			return
// 		}
// 		addParams := func(ps []eph.Params) (err error) {
// 			for _, el := range ps {
// 				if e := addParam(el); e != nil {
// 					err = e
// 					break
// 				}
// 			}
// 			return
// 		}
// 		if e := addParams(op.Params); e != nil {
// 			err = e
// 		} else if e := addParams(op.Locals); e != nil {
// 			err = e
// 		} else if e := addParam(*op.Result); e != nil {
// 			err = e // ^ to match patterns, result (if any) is last.
// 		} else {
// 			k.macros[name] = macroKind{
// 				Kind: g.NewKind(k, name, nil, fields),
// 				init: init,
// 				do:   op.MacroStatements,
// 			}
// 		}
// 	}
// 	return
// }
