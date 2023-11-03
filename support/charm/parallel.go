package charm

// Parallel region; run all of the passed states until they all return nil.
// if none of them return nil, and one of them returns error, this will return the error.
func Parallel(name string, rs ...State) State {
	return Self(name, func(self State, r rune) (ret State) {
		var cnt int
		var lastError State
		for _, s := range rs {
			switch next := s.NewRune(r); next.(type) {
			case nil:
			case Terminal:
				// we could keep them all in the parallel list continually erroring
				// but i think if we just track the last one we can see erroring,
				// it should be fine.
				lastError = next
			default:
				rs[cnt] = next
				cnt++
			}
		}
		if cnt > 0 {
			rs = rs[:cnt]
			ret = self
		} else {
			// if any of them returned error and none returned nil
			// report on that.
			ret = lastError
		}
		return
	})
}
