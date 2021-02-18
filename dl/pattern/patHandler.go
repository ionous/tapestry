package pattern

import "git.sr.ht/~ionous/iffy/rt"

type Handler struct {
	Filter rt.BoolEval
	rt.Execute
}
