package story

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/rt"
)

func (op *TestRule) RewriteActivity() {
	if len(op.Does) == 0 {
		var act = &core.Activity{Exe: []rt.Execute{&debug.DoNothing{}}}
		if opt, ok := op.Hook.Value.(*core.Activity); ok {
			if len(opt.Exe) > 0 {
				act.Exe = opt.Exe
			}
			op.Hook = ProgramHook{}
		}
		core.RewriteActivity(&act, &op.Does)
	}
}

func (op *PatternRule) RewriteActivity() {
	if len(op.Does) == 0 {
		var act = &core.Activity{Exe: []rt.Execute{&debug.DoNothing{}}}
		if opt, ok := op.Hook.Value.(*core.Activity); ok {
			if len(opt.Exe) > 0 {
				act.Exe = opt.Exe
			}
			op.Hook = ProgramHook{}
		}
		core.RewriteActivity(&act, &op.Does)
	}
}
