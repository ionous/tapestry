package list

// for slice and splice: normalize the starting index
// assumes a one based index as input, returns a zero based index as output.
func clipStart(i, cnt int) (ret int) {
	if i == 0 {
		ret = 0 // unspecified: start at the front of the list
	} else if i > cnt {
		ret = -1 // negative return means an empty list
	} else if i > 0 {
		ret = i - 1 // convert to a one-based index
	} else if ofs := cnt + i; ofs > 0 {
		// offset from the end: slice(-2) extracts the last two elements in the sequence.
		ret = ofs
	} else {
		ret = 0
	}
	return
}

// for slice: normalize the ending index
// assumes a one based index as input, returns a zero based index as output.
func clipEnd(j, cnt int) (ret int) {
	if j > cnt {
		ret = cnt
	} else if j > 0 {
		ret = j - 1
	} else if ofs := cnt + j; ofs > 0 {
		ret = ofs // convert to a one-based index
	} else {
		ret = -1 // negative return means an empty list
	}
	return
}

// for splice: turn a starting index and a number of elements from that index into an ending index
func clipRange(start, rng, cnt int) (ret int) {
	if rng <= 0 {
		ret = start // cut start to start, ie. nothing
	} else if end := start + rng; end < cnt {
		ret = end // up to but not including the end
	} else {
		ret = cnt // every element
	}
	return
}
