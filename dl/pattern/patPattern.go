package pattern

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/term"
	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"git.sr.ht/~ionous/iffy/rt/scope"
	"github.com/ionous/errutil"
)

type Pattern struct {
	Name    string
	Params  []term.Preparer
	Locals  []term.Preparer
	Returns term.Preparer
	Rules   []*Rule
}

// NewRecord - create a new record capable of holding a pattern's parameters, locals, and return values.
// return the record and the index of the return value ( or negative 1 if none )
func (ps *Pattern) NewRecord(run rt.Runtime) (ret *g.Record, err error) {
	var ts term.Terms
	if _, e := ps.computeParams(run, &ts); e != nil {
		err = e
	} else if _, e := ps.computeLocals(run, &ts); e != nil {
		err = e
	} else if _, e := ps.computeReturn(run, &ts); e != nil {
		err = e
	} else if vs, e := ts.NewRecord(run); e != nil {
		err = e
	} else {
		ret = vs
	}
	return
}

// Run - given a record ( returned by new record ) and an expected return affinity.
// push the record onto the scope, execute the matching pattern rules, and return its result.
func (ps *Pattern) Run(run rt.Runtime, rec *g.Record, aff affine.Affinity) (ret g.Value, err error) {
	res := -1
	if ps.Returns != nil {
		res = rec.Kind().FieldIndex(ps.Returns.String())
	}
	//
	scope := scope.TargetRecord{object.Variables, rec}
	run.PushScope(&scope)
	if e := ps.run(run); e != nil {
		err = e
	} else if res, e := getResult(rec, res, aff); e != nil {
		err = e
	} else {
		ret = res
	}
	run.PopScope()
	return
}

// RunWithScope - note: assumes whatever scope is needed to run the pattern has already been setup.
func (ps *Pattern) run(run rt.Runtime) (err error) {
	if inds, e := splitRules(run, ps.Rules); e != nil {
		err = e
	} else {
		for _, i := range inds {
			if e := safe.Run(run, ps.Rules[i].Execute); e != nil {
				err = e
				break
			}
			// NOTE: if we need to differentiate between "ran" and "not found",
			// "didnt run" should probably become an error code.
		}
	}
	return
}

func (ps *Pattern) computeParams(run rt.Runtime, terms *term.Terms) (ret int, err error) {
	return prepareList(run, ps.Params, terms)
}

func (ps *Pattern) computeLocals(run rt.Runtime, terms *term.Terms) (ret int, err error) {
	return prepareList(run, ps.Locals, terms)
}

func (ps *Pattern) computeReturn(run rt.Runtime, terms *term.Terms) (ret int, err error) {
	if res := ps.Returns; res == nil {
		ret = -1
	} else if v, e := res.Prepare(run); e != nil {
		err = e
	} else {
		n := res.String()
		ret = terms.AddValue(n, v)
	}
	return
}

func prepareList(run rt.Runtime, list []term.Preparer, terms *term.Terms) (ret int, err error) {
	for _, n := range list {
		if v, e := n.Prepare(run); e != nil {
			err = errutil.Append(err, e)
		} else {
			terms.AddValue(n.String(), v)
			ret++
		}
	}
	return
}
