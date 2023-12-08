package cmdcompact

import "path/filepath"

const (
	SpecExt = ".ifspecs"
	//
	DetailedExt = ".ifx"
	CompactExt  = ".if"
	BlockExt    = ".block"
	//
	TellStory = ".tell"
	TellSpec  = ".tells"
)

var allExts = []string{SpecExt, DetailedExt, CompactExt, BlockExt, TellSpec, TellStory}

// fix: probably should have an ext type
// maybe use comment generator to produce it.
// then you could get rid of all allExts

func oppositeExt(ext string) (ret string) {
	if ext == CompactExt {
		ret = DetailedExt
	} else if ext == DetailedExt {
		ret = CompactExt
	} else {
		ret = ext
	}
	return
}

func isSpecExt(ext string) bool {
	return ext == SpecExt || ext == TellSpec
}

func formatOf(path string) (ret format) {
	ext := filepath.Ext(path)
	if ext == TellSpec || ext == TellStory {
		ret = tellFormat
	} else {
		ret = jsonFormat
	}
	return
}
