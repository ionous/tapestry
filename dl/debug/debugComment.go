package debug

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/rt"
)

func (op *Comment) PostImport(k *imp.Importer) (err error) {
	// do nothing for now.
	return
}

// turn a comment execute statement into a debug log.
func (op *Comment) PreImport(k *imp.Importer) (ret interface{}, err error) {
	if !k.Env().InProgram() {
		ret = op
	} else {
		ret = &DebugLog{
			Value: core.AssignFromText(core.T(op.Lines.String())),
			LogLevel: LoggingLevel{
				Str: LoggingLevel_Note,
			},
		}
	}
	return
}

var _ imp.PreImport = (*Comment)(nil)

func (*Comment) Execute(rt.Runtime) error {
	panic("unexpected use of story method")
}
