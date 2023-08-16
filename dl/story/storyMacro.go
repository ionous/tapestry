package story

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

// Execute - called by the macro runtime during weave.
func (op *DefineMacro) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// Schedule - register the macro with the importer;
// subsequent CallMacro(s) will be able to run it.
func (op *DefineMacro) Weave(cat *weave.Catalog) (err error) {
	return cat.Schedule(weave.RequirePlurals, func(w *weave.Weaver) (err error) {
		if name, e := safe.GetText(cat.Runtime(), op.MacroName); e != nil {
			err = e
		} else {
			pb := mdl.NewPatternSubtype(name.String(), kindsOf.Macro)
			if e := addFields(pb, mdl.PatternLocals, op.Locals); e != nil {
				err = e
			} else if e := addFields(pb, mdl.PatternParameters, op.Params); e != nil {
				err = e
			} else if e := addOptionalField(pb, mdl.PatternResults, op.Result); e != nil {
				err = e
			} else {
				/**/ pb.AddRule("", &core.Always{}, 2, op.MacroStatements)
				if e := w.Pin().AddPattern(pb.Pattern); e != nil {
					err = e
				}
			}
		}
		return
	})
}

// Schedule for macros calls Execute... eventually... to generate dynamic assertions.
func (op *CallMacro) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequireNouns, func(w *weave.Weaver) error {
		return op.Execute(w)
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
	name := lang.Normalize(op.MacroName)
	if k, v, e := assign.ExpandArgs(run, op.Arguments); e != nil {
		err = assign.CmdError(op, e)
	} else if v, e := run.Call(name, aff, k, v); e != nil {
		err = assign.CmdError(op, e)
	} else {
		ret = v
	}
	return
}
