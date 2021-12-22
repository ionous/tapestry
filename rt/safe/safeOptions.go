package safe

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/meta"
)

// safety for built in options means a panic if the option doesnt exist
func GetFlag(run rt.Runtime, opt meta.Options) (ret bool) {
	if v, e := run.GetField(meta.Option, opt.String()); e != nil {
		panic(e)
	} else if e := Check(v, affine.Bool); e != nil {
		panic(e)
	} else {
		ret = v.Bool()
	}
	return
}
