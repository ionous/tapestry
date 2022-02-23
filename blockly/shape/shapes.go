package shape

var rootBlocks = RootBlocks{"story_file"}

// root blocks have no output
type RootBlocks []string

func (x RootBlocks) IsRoot(name string) (ret bool) {
	for _, str := range x {
		if name == str {
			ret = true
			break
		}
	}
	return
}
