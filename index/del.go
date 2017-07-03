package index

// DeletionCursor cuts a series of pairs from an index.
type DeletionCursor struct {
	*Index
	start, next int
}

// NewDeletionCursor
func NewDeletionCursor(index *Index) *DeletionCursor {
	return &DeletionCursor{Index: index}
}

// Flush preforms the actual deletion. It should only be called after some succesful call to DeletePair.
func (dc *DeletionCursor) Flush() {
	a, i, n := dc.Lines, dc.start, dc.next
	dc.Lines = a[:i+copy(a[i:], a[n:])]
	// println("cut", i, n)
}

// DeletePair marks a pair for deletion, returning true if it found the pair.
func (dc *DeletionCursor) DeletePair(majorKey, minorKey string) (okay bool) {
	// initially, just find the pair and record the cut point.
	// println("searching for", majorKey, minorKey)
	if dc.start == dc.next {
		if i, ok := dc.FindPair(dc.start, majorKey, minorKey); ok {
			dc.start, dc.next, okay = i, i+1, true
			// println("marked start", dc.start, dc.next)
		}
	} else {
		// later, check the direct next slot first in the hopes of consolidating cuts.
		a, major, minor := dc.Lines, dc.Major, (dc.Major+1)&1
		if n := dc.next; n < len(a) && a[n].Key[major] == majorKey && a[n].Key[minor] == minorKey {
			dc.next, okay = n+1, true
			// println("consolidated", dc.start, dc.next)
		} else {
			// we couldnt consolidate this time, cut out the recorded range.
			dc.Flush()
			// we still need to find the pair in question, and record the cut point.
			// we cut everything down to start, so the slot we just checked is now one beyond start.
			if i, ok := dc.FindPair(dc.start+1, majorKey, minorKey); ok {
				dc.start, dc.next, okay = i, i+1, true
				// println("marked next", dc.start, dc.next)
			} else {
				// if we didnt find it, then we have nothing to cut.
				// we're at our initial state again.
				dc.start, dc.next = i, i
				// println("start again at", dc.start, dc.next)
			}
		}
	}
	return
}
