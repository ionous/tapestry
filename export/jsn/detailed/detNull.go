package detailed

import "git.sr.ht/~ionous/iffy/export/jsn"

// extend NullMarshaler with empty detailedWriter methods
type detNull struct {
	jsn.NullMarshaler
}

func (detNull) commit(value interface{}) {}
