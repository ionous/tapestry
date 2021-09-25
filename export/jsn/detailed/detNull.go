package detailed

import "git.sr.ht/~ionous/iffy/export/jsn"

// extend NullMarshaler with empty detailedWriter methods
type detNull struct {
	jsn.NullMarshaler
}

func (detNull) named() string               { return "null" }
func (detNull) writeData(value interface{}) {}
func (detNull) readData() interface{}       { return nil }
