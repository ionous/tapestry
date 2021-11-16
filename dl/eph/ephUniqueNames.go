package eph

import "sort"

type UniqueNames []string

func (ns *UniqueNames) AddName(name string) {
	s := *ns
	i := sort.SearchStrings(s, name)
	if i == len(s) {
		(*ns) = append(s, name)
	} else if s[i] != name {
		s = append(s, "")
		copy(s[i+1:], s[i:])
		s[i] = name
		(*ns) = s
	}
}
