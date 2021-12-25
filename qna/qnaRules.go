package qna

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/rt"
	"github.com/ionous/errutil"
)

func (run *Runner) GetRules(pattern, target string, pflags *rt.Flags) (ret []rt.Rule, err error) {
	pat, tgt := lang.Underscore(pattern), lang.Underscore(target) // FIX: caller normalization would be best.
	if rs, e := run.getRules(pat, tgt); e != nil {
		err = e
	} else {
		if pflags != nil {
			*pflags = rs.flags
		}
		ret = rs.rules
	}
	return
}

// stored in Runner.cache
type cachedRules struct {
	rules []rt.Rule
	flags rt.Flags // sum of flags of each rule
}

// get the rules from the cache, or build them and add them to the cache
func (run *Runner) getRules(pat, tgt string) (ret cachedRules, err error) {
	if c, e := run.values.cache(func() (ret interface{}, err error) {
		if rules, flags, e := run.buildRules(pat, tgt); e != nil {
			err = run.Report(e)
		} else {
			ret = cachedRules{rules: rules, flags: flags}
		}
		return
	}, "rules", pat, tgt); e != nil {
		err = e
	} else {
		ret = c.(cachedRules)
	}
	return
}

// build the rules for the passed pat,tgt pair
func (run *Runner) buildRules(pat, tgt string) (retRules []rt.Rule, retFlags rt.Flags, err error) {
	if els, e := run.qdb.RulesFor(pat, tgt); e != nil {
		err = e
	} else {
		var rules []rt.Rule
		var flags rt.Flags
		for _, el := range els {
			var filter rt.BoolEval
			// fix: we dont want to be bound to core here,
			// even though we need its custom decoding handlers
			// probably best is to pass a decoder object in -- not even just run.signatures
			// ( then you could run it against detail encoding too if you wanted )
			// might also stack the custom decoders just like with the signatures --
			// then you just "loop" over them maybe
			if e := core.Decode(rt.BoolEval_Slot{&filter}, el.Filter, run.signatures); e != nil {
				e = errutil.New("error decoding filter for", pat, tgt, el.Id, e)
				err = errutil.Append(err, e)
			} else {
				var prog rt.Execute
				if e := core.Decode(rt.Execute_Slot{&prog}, el.Prog, run.signatures); e != nil {
					e = errutil.New("error decoding prog for", pat, tgt, el.Id, e)
					err = errutil.Append(err, e)
				} else {
					flags := rt.MakeFlags(rt.Phase(el.Phase))
					rules = append(rules, rt.Rule{
						Name:     el.Id,
						Filter:   filter,
						Execute:  prog,
						RawFlags: float64(flags),
					})
					flags |= flags
				}
			}
		}
		if err == nil {
			retRules, retFlags = rules, flags
		}
	}
	return
}
