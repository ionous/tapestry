package safe

import (
	"git.sr.ht/~ionous/iffy/rt"
)

func ApplyRule(run rt.Runtime, rule rt.Rule, allow rt.Flags) (ret rt.Flags, err error) {
	if flags := rule.GetFlags(); allow&flags != 0 {
		if ok, e := GetOptionalBool(run, rule.Filter, true); e != nil {
			err = e
		} else if ok.Bool() { // the rule returns false if it didnt apply
			if e := Run(run, rule.Execute); e != nil {
				err = e
			} else {
				ret = flags
			}
		}
	}
	return
}
