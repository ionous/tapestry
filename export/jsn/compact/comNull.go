package compact

import "git.sr.ht/~ionous/iffy/export/jsn"

// a do-nothing state after encountering an unrecoverable error.
type comNull struct {
	jsn.NullMarshaler
}

func (d *comNull) named() string         { return "null" }
func (d *comNull) writeData(interface{}) {}
func (d *comNull) readData() interface{} { return nil }
