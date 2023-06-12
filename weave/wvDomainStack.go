package weave

// DomainStack - tracks the current target of calls to schedule().
type DomainStack []*Domain

func (k *DomainStack) Push(b *Domain) {
	(*k) = append(*k, b)
}

func (k *DomainStack) Top() (ret *Domain, okay bool) {
	if end := len(*k) - 1; end >= 0 {
		ret, okay = (*k)[end], true
	}
	return
}

// return false if empty
func (k *DomainStack) Pop() (ret *Domain, okay bool) {
	if end := len(*k) - 1; end >= 0 {
		ret, (*k) = (*k)[end], (*k)[:end]
		okay = true
	}
	return
}
