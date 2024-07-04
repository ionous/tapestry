package rt

// error code for system commands
// used by play/step to switch take some action as a side-effect a command
type Signal int

//go:generate stringer -type=Signal -trimprefix=Signal
const (
	SignalUnknown Signal = iota
	SignalQuit
	SignalSave
	SignalLoad
	SignalBreak
	SignalContinue
)

func (s Signal) Error() string {
	return s.String()
}
