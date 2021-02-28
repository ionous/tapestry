package pattern

import "git.sr.ht/~ionous/iffy/rt"

// previously we matched the rules, and then ran them.
// now: they are sorted first, and then matched so rules can affect each other.
// FIX: presort everything in the assembler?
func SortRules(rules []rt.Rule) (ret []int, retFlags rt.Flags) {
	cnt := len(rules)
	var pre, post []int
	in := make([]int, 0, cnt)
	for i := cnt - 1; i >= 0; i-- {
		flags := rules[i].GetFlags()
		var at *[]int
		switch flags {
		case rt.Prefix:
			at = &pre
		case rt.Postfix:
			at = &post
		case rt.Infix:
			at = &in
		}
		if at != nil {
			*at = append(*at, i)
			retFlags |= flags
		}
	}
	if retFlags == rt.Infix {
		ret = in // this is the most common
	} else {
		ret = append(pre, append(in, post...)...)
	}
	return
}
