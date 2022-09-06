package debug

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/rt"
)

func (op *Comment) ImportPhrase(k *story.Importer) (err error) {
	// do nothing for now.
	return
}

// turn a comment execute statement into a debug log.
func (op *Comment) ImportStub(k *story.Importer) (ret interface{}, err error) {
	if !k.InProgram() {
		ret = op
	} else {
		ret = &DebugLog{
			Value: &core.FromText{
				Val: core.T(op.Lines.String()),
			},
			LogLevel: LoggingLevel{
				Str: LoggingLevel_Note,
			},
		}
	}
	return
}

var _ story.StubImporter = (*Comment)(nil)

func (*Comment) Execute(rt.Runtime) error {
	panic("unexpected use of story method")
}
