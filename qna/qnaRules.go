package qna

import (
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

// return the runtime rules matching the passed pattern and target
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
	if c, e := run.values.cache(func() (interface{}, error) {
		return run.buildRules(pat, tgt)
	}, "rules", pat, tgt); e != nil {
		err = e
	} else {
		ret = c.(cachedRules)
	}
	return
}

// build the rules for the passed pat,tgt pair
func (run *Runner) buildRules(pat, tgt string) (ret cachedRules, err error) {
	if els, e := run.query.RulesFor(pat, tgt); e != nil {
		err = e
	} else {
		var rules []rt.Rule
		var sum rt.Flags
		for _, el := range els {
			if filter, e := run.decode.DecodeFilter(el.Filter); e != nil {
				err = errutil.Append(err, errutil.New("decoding filter", pat, tgt, el.Id, e))
			} else if prog, e := run.decode.DecodeProg(el.Prog); e != nil {
				err = errutil.Append(err, errutil.New("decoding prog", pat, tgt, el.Id, e))
			} else {
				flags := rt.MakeFlags(rt.Phase(el.Phase))
				rules = append(rules, rt.Rule{
					Name:     el.Id,
					Filter:   filter,
					Execute:  prog,
					RawFlags: float64(flags),
				})
				sum |= flags
			}
		}
		if err == nil {
			ret = cachedRules{rules: rules, flags: sum}
		}
	}
	return
}
