package weave

import (
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

type PhaseActions map[assert.Phase]PhaseAction

type PhaseAction struct {
	Flags PhaseFlags
	Do    func(d *Domain) (err error)
}

type PhaseFlags struct {
	NoDuplicates bool
}

// shared generic marshal prog to text
func marshalout(cmd interface{}) (ret string, err error) {
	if cmd != nil {
		if m, ok := cmd.(jsn.Marshalee); !ok {
			err = errutil.Fmt("can only marshal autogenerated types (%T)", cmd)
		} else {
			ret, err = cout.Marshal(m, literal.CompactEncoder)
		}
	}
	return
}
