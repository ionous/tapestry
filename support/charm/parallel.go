package charm

// Parallel region; run all of the passed states until they all return nil.
// if any return error, this returns error.
func Parallel(name string, rs ...State) State {
	return Self(name, func(self State, r rune) (ret State) {
		var cnt int
	Loop:
		for _, s := range rs {
			switch next := s.NewRune(r); next.(type) {
			case nil:
				// skip
			case Terminal:
				ret = next
				break Loop
			default:
				rs[cnt] = next
				cnt++
			}
		}
		if cnt > 0 && ret == nil {
			rs = rs[:cnt]
			ret = self
		}
		return
	})
}
