package eph

type DomainList []*Domain

// for now, duplicates are okay.
func (dl *DomainList) add(d *Domain) {
	if !dl.contains(d) {
		(*dl) = append(*dl, d)
	}
}

func (dl *DomainList) contains(d *Domain) (okay bool) {
	for _, el := range *dl {
		if el == d {
			okay = true
			break
		}
	}
	return
}
