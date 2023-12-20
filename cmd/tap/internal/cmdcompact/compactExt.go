package cmdcompact

import "git.sr.ht/~ionous/tapestry/support/files"

func formatOf(path string) (ret format) {
	if ext := files.Ext(path); ext.Tell() {
		ret = tellFormat
	} else if ext.Json() {
		ret = jsonFormat
	}
	return
}
