package pattern

import (
	"strconv"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

type Pattern struct {
	Name   string
	Return string            // name of return field; empty if none ( could be an index but slightly safer this way )
	Labels []string          // one label for every parameter
	Locals []core.Assignment // usually equal to the number of locals; or nil for testing.
	Fields []g.Field         // flat list of params and locals and an optional return
	Rules  []*Rule
}

func (pat *Pattern) Run(run rt.Runtime, args []*core.Argument, aff affine.Affinity) (ret g.Value, err error) {
	// create a container to hold results of args, locals, and the pending return value
	rec := g.NewAnonymousRecord(run, pat.Fields)
	// args run in the scope of their parent context
	// they write to the record that will become the new context
	if e := pat.determineArgs(run, rec, args); e != nil {
		err = e
	} else {
		// initializers ( and the pattern itself ) run in the scope of the pattern
		// ( with access to all locals and args)
		oldScope := run.ReplaceScope(g.RecordOf(rec))
		// locals ( by definition ) write to the record context
		if e := pat.initializeLocals(run, rec); e != nil {
			err = e
		} else if e := pat.executePattern(run); e != nil {
			err = e
		} else if res, e := pat.getResult(rec, aff); e != nil {
			err = e
		} else {
			ret = res
		}
		//
		run.ReplaceScope(oldScope)
	}
	if err != nil {
		err = errutil.New(pat.Name, err.Error())
	}
	return
}

func (pat *Pattern) determineArgs(run rt.Runtime, rec *g.Record, args []*core.Argument) (err error) {
	// future: args matching is predetermined in reading / parsing
	// note: templates (ex. print_article) dont always specify all the parameters...
	if paramCnt, argCnt := len(pat.Labels), len(args); paramCnt < argCnt {
		err = errutil.New("expected", paramCnt, "parameters(s), have", argCnt, "arguments")
	} else {
		// note: set indexed field assigns without copying
		for i, a := range args {
			if n, l := a.Name, pat.Labels[i]; n != l && n != argIndex(i) {
				err = errutil.New("has mismatched arg.", i, "expected", l, "have", n)
				break
			} else if val, e := core.GetAssignedValue(run, a.From); e != nil {
				err = errutil.New("error determining arg", i, n, e)
				break
			} else if v, e := filterText(run, pat.Fields[i], val); e != nil {
				err = errutil.New("error narrowing arg", i, n, e)
				break
			} else if e := rec.SetIndexedField(i, v); e != nil {
				err = errutil.New("error setting arg", i, n, e)
				break
			}
		}
	}
	return
}

// fix? allows callers to use positional arguments
// for lists could have a special RunWithVarArgs that uses a custom determineArgs
// or, allow blank names to match any arg --
// note: templates currently use positional args too.
func argIndex(i int) string {
	return "$" + strconv.Itoa(i+1)
}

func (pat *Pattern) initializeLocals(run rt.Runtime, rec *g.Record) (err error) {
	fin := len(pat.Labels) // locals start after labels
	for i, init := range pat.Locals {
		field := pat.Fields[i+fin]
		if init != nil {
			if v, e := init.GetAssignedValue(run); e != nil {
				err = errutil.New(pat.Name, "error determining local", i, field.Name, e)
				break
			} else if e := rec.SetIndexedField(i, v); e != nil {
				err = errutil.New(pat.Name, "error setting local", i, field.Name, e)
				break
			}
		}
	}
	return
}

// RunWithScope - note: assumes whatever scope is needed to run the pattern has already been setup.
func (pat *Pattern) executePattern(run rt.Runtime) (err error) {
	if inds, e := splitRules(run, pat.Rules); e != nil {
		err = e
	} else {
		for _, i := range inds {
			if e := safe.Run(run, pat.Rules[i].Execute); e != nil {
				err = e
				break
			}
			// NOTE: if we need to differentiate between "ran" and "not found",
			// "didnt run" should probably become an error code.
		}
	}
	return
}

func (pat *Pattern) getResult(rec *g.Record, aff affine.Affinity) (ret g.Value, err error) {
	// labels=parameters, inits=locals, the rest ( no more than 1 ) is the return.
	if res := pat.Return; len(res) > 0 {
		// get the value and check its result
		if res, e := rec.GetNamedField(res); e != nil {
			err = errutil.New("error trying to get return value", e)
		} else if e := safe.Check(res, aff); e != nil {
			err = errutil.New("error trying to get return value", e)
		} else if len(aff) == 0 {
			// the caller expects nothing but we have a return value.
			if res.Affinity() == affine.Text {
				core.HackTillTemplatesCanEvaluatePatternTypes = res.String()
			} else {
				err = errutil.New("the caller expects nothing but we returned", aff)
			}
		} else {
			ret = res
		}
	} else if len(aff) != 0 {
		err = errutil.New("caller expected", aff, "returned nothing")
	}
	return
}
