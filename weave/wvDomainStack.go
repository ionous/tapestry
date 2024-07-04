package weave

// tracks the current target of calls to schedule().
// ( often two deep: the tapestry domain,
// . and the domain being processed. )
type SceneStack []*Domain

func (k *SceneStack) Push(d *Domain) {
	// println("PUSHING", d.name)
	(*k) = append(*k, d)
}

func (k *SceneStack) Top() (ret *Domain, okay bool) {
	if end := len(*k) - 1; end >= 0 {
		ret, okay = (*k)[end], true
	}
	return
}

// return false if empty
func (k *SceneStack) Pop() (ret *Domain, okay bool) {
	if end := len(*k) - 1; end >= 0 {
		ret, (*k) = (*k)[end], (*k)[:end]
		// println("POPPING", ret.name)
		okay = true
	}
	return
}
