package print

import "git.sr.ht/~ionous/tapestry/rt/writer"

// Filter - sends incoming chunks to one of three functions: first, rest, or last.
// a user of filter can alter those chunks before sending them onward somewhere else.
// this doesnt define that destination.
type Filter struct {
	First, Rest func(writer.Chunk) (int, error)
	Last        func(int) error
	cnt         int
}

func (f *Filter) WriteChunk(c writer.Chunk) (ret int, err error) {
	if c.IsClosed() {
		if f.Last != nil {
			err = f.Last(f.cnt)
		}
	} else {
		if f.cnt == 0 && f.First != nil {
			ret, err = f.First(c)
		} else if f.Rest != nil {
			ret, err = f.Rest(c)
		}
		f.cnt += ret
	}
	return
}
