package weave

import "sort"

type UniqueNames []string

// returns the index if newly added ( or -1 if not )
func (ns *UniqueNames) AddName(name string) (ret int) {
	s := *ns
	if i := sort.SearchStrings(s, name); i == len(s) {
		(*ns) = append(s, name)
		ret = len(s)
	} else if s[i] != name {
		s = append(s, "")
		copy(s[i+1:], s[i:])
		s[i] = name
		(*ns) = s
		ret = i
	} else {
		ret = -1 // not found
	}
	return
}
