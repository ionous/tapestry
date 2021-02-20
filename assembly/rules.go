package assembly

import (
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

type BuildRule struct {
	Query        string
	NewContainer func(name string) interface{}
	NewEl        func(c interface{}) interface{}
}

// the first parameter should be a *string, the second some *bytes
func (b *BuildRule) buildFromRule(asm *Assembler, args ...interface{}) (err error) {
	list := make(map[string]interface{})
	var last string
	var curr interface{}
	if e := tables.QueryAll(asm.cache.DB(), b.Query,
		func() (err error) {
			name, prog := *args[0].(*string), *args[1].(*[]byte)
			if name != last || curr == nil {
				curr = b.NewContainer(name)
				list[name] = curr
				last = name
			}
			return tables.DecodeGob(prog, b.NewEl(curr))
		}, args...); e != nil {
		err = errutil.New("buildFromRule", e)
	} else {
		// write the passed list of gobs into the assembler db
		err = asm.WriteGobs(list)
	}
	return
}
