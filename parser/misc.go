package parser

// Words changes the list of text into individual options
func Words(s []string) (ret Scanner) {
	if cnt := len(s); cnt == 1 {
		ret = &Word{s[0]}
	} else {
		words := make([]Scanner, cnt)
		for i, w := range s {
			words[i] = &Word{w}
		}
		ret = &AnyOf{words}
	}
	return
}
