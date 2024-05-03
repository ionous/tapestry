package generic

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"github.com/ionous/errutil"
)

// Kinds database
// this isnt used by package generic, but its a common enough interface for tests and the runtime
type Kinds interface {
	GetKindByName(n string) (*Kind, error)
}
