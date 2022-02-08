package core

import "git.sr.ht/~ionous/tapestry/rt"

func RewriteActivity(act **Activity, out *[]rt.Execute) {
	if *act != nil {
		*out = (*act).Exe
		*act = nil
	} else if len(*out) == 0 {
		*out = make([]rt.Execute, 0)
	}
}

func (op *BracketText) RewriteActivity() {
	RewriteActivity(&op.Do, &op.Does)
}
func (op *BufferText) RewriteActivity() {
	RewriteActivity(&op.Do, &op.Does)
}
func (op *ChooseAction) RewriteActivity() {
	RewriteActivity(&op.Do, &op.Does)
}
func (op *ChooseMore) RewriteActivity() {
	RewriteActivity(&op.Do, &op.Does)
}
func (op *ChooseMoreValue) RewriteActivity() {
	RewriteActivity(&op.Do, &op.Does)
}
func (op *ChooseNothingElse) RewriteActivity() {
	RewriteActivity(&op.Do, &op.Does)
}
func (op *ChooseValue) RewriteActivity() {
	RewriteActivity(&op.Do, &op.Does)
}
func (op *CommaText) RewriteActivity() {
	RewriteActivity(&op.Do, &op.Does)
}
func (op *Row) RewriteActivity() {
	RewriteActivity(&op.Do, &op.Does)
}
func (op *Rows) RewriteActivity() {
	RewriteActivity(&op.Do, &op.Does)
}
func (op *SlashText) RewriteActivity() {
	RewriteActivity(&op.Do, &op.Does)
}
func (op *SpanText) RewriteActivity() {
	RewriteActivity(&op.Do, &op.Does)
}
func (op *While) RewriteActivity() {
	RewriteActivity(&op.Do, &op.Does)
}
