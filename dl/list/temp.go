package list

import "git.sr.ht/~ionous/tapestry/dl/core"

func (op *Erasing) RewriteActivity() {
	core.RewriteActivity(&op.Do, &op.Does)
}
func (op *ErasingEdge) RewriteActivity() {
	core.RewriteActivity(&op.Do, &op.Does)
}
func (op *ListEach) RewriteActivity() {
	core.RewriteActivity(&op.Do, &op.Does)
}
