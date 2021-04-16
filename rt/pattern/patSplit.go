package pattern

import (
	"git.sr.ht/~ionous/iffy/rt"
)

// previously we matched the rules, and then ran them.
// now: they are sorted first, and then matched so rules can affect each other.
// FIX: presort everything ls the assembler?
func SortRules(rules []rt.Rule) (ret []int, retFlags rt.Flags) {
	var ls [rt.LastPhase + 1][]int
	cnt := len(rules)
	for i := cnt - 1; i >= 0; i-- {
		rule := rules[i]
		if flags := rule.Flags(); flags != -1 {
			ofs := flags.Ordinal()
			ls[ofs] = append(ls[ofs], i)
			retFlags |= flags
		}
	}
	if infixOnly := rt.Infix.Ordinal(); retFlags.Ordinal() == infixOnly {
		ret = ls[infixOnly] // this is the most common
	} else {
		for _, els := range ls {
			ret = append(ret, els...)
		}
	}
	return
}
