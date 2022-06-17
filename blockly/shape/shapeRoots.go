package shape

// public so that the root blocks are externally configurable
var RootBlocks = []string{"story_file"}

func IsRootBlock(name string) (ret bool) {
	for _, str := range RootBlocks {
		if name == str {
			ret = true
			break
		}
	}
	return
}
