package story

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/eph"
)

// export an internal function just for testing.
func ImportCall(op *core.CallPattern) *eph.EphRefs {
	return importCall(op)
}
