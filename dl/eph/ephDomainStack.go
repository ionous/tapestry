package eph

// DomainStack - keep track of the current block while marshaling.
// ( so that end block can be called. )
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
