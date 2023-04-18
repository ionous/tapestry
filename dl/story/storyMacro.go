package story

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/imp/assert"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// Execute - called by the macro runtime during weave.
func (op *DefineMacro) Execute(macro rt.Runtime) error {
	return imp.StoryStatement(macro, op)
}

// PostImport - register the macro with the importer;
// subsequent CallMacro(s) will be able to run it.
func (op *DefineMacro) PostImport(k *imp.Importer) (err error) {
	if name, e := safe.GetText(k, op.MacroName); e != nil {
		err = e
	} else {
		macro := name.String()
		if e := k.AssertAncestor(macro, kindsOf.Macro.String()); e != nil {
			err = e
		} else {
			if res := op.Result; res != nil {
				err = res.DeclareField(func(name, class string, aff affine.Affinity, init assign.Assignment) error {
					return k.AssertResult(macro, name, class, aff, init)
				})
			}
			if err == nil {
				if e := declareFields(op.Params, func(name, class string, aff affine.Affinity, init assign.Assignment) error {
					return k.AssertLocal(macro, name, class, aff, init)
				}); e != nil {
					err = e
				} else if e := declareFields(op.Locals, func(name, class string, aff affine.Affinity, init assign.Assignment) error {
					return k.AssertLocal(macro, name, class, aff, init)
				}); e != nil {
					err = e
				} else {
					err = k.AssertRule(macro, "", nil, 0, op.MacroStatements)
				}
			}
		}
	}
	return
}

// PostImport for macros calls Execute... eventually... to generate dynamic assertions.
func (op *CallMacro) PostImport(k *imp.Importer) error {
	return k.Schedule(assert.MacroPhase, func(assert.World, assert.Assertions) error {
		return op.Execute(k)
	})
}

func (op *CallMacro) Execute(run rt.Runtime) error {
	_, err := op.determine(run, affine.None)
	return err
}

func (op *CallMacro) GetBool(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.Bool)
}

func (op *CallMacro) GetNumber(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.Number)
}

func (op *CallMacro) GetText(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.Text)
}

func (op *CallMacro) GetRecord(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.Record)
}

func (op *CallMacro) GetNumList(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.NumList)
}

func (op *CallMacro) GetTextList(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.TextList)
}

func (op *CallMacro) GetRecordList(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.RecordList)
}

func (op *CallMacro) determine(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	name := lang.Underscore(op.MacroName)
	if rec, e := assign.MakeRecord(run, name, op.Arguments...); e != nil {
		err = assign.CmdError(op, e)
	} else if v, e := run.Call(rec, aff); e != nil {
		err = assign.CmdError(op, e)
	} else {
		ret = v
	}
	return
}
