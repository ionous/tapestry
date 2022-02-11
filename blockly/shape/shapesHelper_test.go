package shape

import (
	"io/fs"
	"sort"

	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/web/js"
)

func ReadSpec(files fs.FS, fileName string) (ret *spec.TypeSpec, err error) {
	return readSpec(files, fileName)
}

func Lookup(k string) (ret *spec.TypeSpec, okay bool) {
	ret, okay = lookup[k]
	return
}

func LookupKeys() []string {
	var keys []string
	for k, _ := range lookup {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func ResetLookup() {
	lookup = make(TypeSpecs) // reset
}

func Write(block *js.Builder, blockType *spec.TypeSpec) bool {
	return writeShape(block, blockType)
}
