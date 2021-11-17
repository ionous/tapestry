package eph

import "sort"

type UniqueNames []string

// returns true if newly added
func (ns *UniqueNames) AddName(name string) (okay bool) {
	s := *ns
	if i := sort.SearchStrings(s, name); i == len(s) {
		(*ns) = append(s, name)
		okay = true
	} else if s[i] != name {
		s = append(s, "")
		copy(s[i+1:], s[i:])
		s[i] = name
		(*ns) = s
		okay = true
	}
	return
}
