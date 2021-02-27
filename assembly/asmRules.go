package assembly

import (
	"git.sr.ht/~ionous/iffy/dl/pattern"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

type TargetDomain struct {
	Target, Domain string
}

type RuleSet struct {
	name string // pattern
	// FIX: this is wrong --
	// probably we have to include domain hierarchy ( more derived first? ) when calculating rule order
	// that way we can turn on and off domains and keep good sorting
	rules map[TargetDomain][]rt.Rule
}

func WriteRules(asm *Assembler, pat, tgt, domain string, rules []rt.Rule) (err error) {
	inds, _ := pattern.SortRules(rules)
	for _, j := range inds {
		rule := rules[j]
		handler := rt.Handler{Filter: rule.Filter, Exe: rule.Execute}
		if prog, e := tables.EncodeGob(&handler); e != nil {
			err = errutil.Append(err, e)
		} else if e := asm.WriteRule(pat, tgt, domain, rule.GetFlags(), prog, rule.Name); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// return pattern meta data for anything that has rules.
func buildPatternRules(asm *Assembler, cache patternCache) (err error) {
	var curr *RuleSet
	var rules []RuleSet
	var name, target, domain string
	var prog []byte
	if e := tables.QueryAll(asm.cache.DB(),
		`select pattern, target, domain, prog 
		from asm_rule where type='rule'
		order by pattern, type, domain, idProg`,
		func() (err error) {
			newPattern := curr == nil || curr.name != name
			if newPattern && cache[name] == nil {
				err = errutil.New("unknown pattern", name)
			} else {
				if newPattern {
					rules = append(rules, RuleSet{name: name})
					curr = &rules[len(rules)-1]
					curr.rules = make(map[TargetDomain][]rt.Rule)
				}
				var rule rt.Rule
				if e := tables.DecodeGob(prog, &rule); e != nil {
					err = e
				} else {
					key := TargetDomain{target, domain}
					rs := curr.rules[key]
					curr.rules[key] = append(rs, rule)
				}
			}
			return
		}, &name, &target, &domain, &prog); e != nil {
		err = e
	} else {
		// array of all possible rules grouped by pattern
		for _, rs := range rules {
			// within each pattern are its rules grouped by possible event target and domain
			for td, ls := range rs.rules {
				// fix: domains should have hierarchy and we should list rules by decreasing depth
				// ( so that deeper, more derived domains take precedence )
				if td.Domain != "entire_game" {
					e := WriteRules(asm, rs.name, td.Target, td.Domain, ls)
					err = errutil.Append(err, e)
				}
			}
			for td, ls := range rs.rules {
				if td.Domain == "entire_game" {
					e := WriteRules(asm, rs.name, td.Target, td.Domain, ls)
					err = errutil.Append(err, e)
				}
			}
		}
	}
	return
}
