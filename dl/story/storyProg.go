package story

// note: there's only one kind of hook now: the activity
// func (op *ProgramHook) ImportProgram(k *Importer) (ret rt.Execute, err error) {
// 	if opt, ok := op.Value.(*core.Activity); !ok {
// 		err = ImportError(op, errutil.Fmt("%w for %T", UnhandledSwap, op.Value))
// 	} else {
// 		switch len(opt.Exe) {
// 		case 0:
// 			ret = &debug.DoNothing{}
// 		case 1:
// 			ret = opt.Exe[0]
// 		default:
// 			ret = opt
// 		}
// 	}
// 	return
// }
