package tapestry

import (
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/reader"
	"git.sr.ht/~ionous/tapestry/dl/rel"
	"git.sr.ht/~ionous/tapestry/dl/render"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"

	r "reflect"
)

var AllSlats = [][]composer.Composer{
	core.Slats,
	debug.Slats,
	grammar.Slats,
	literal.Slats,
	list.Slats,
	reader.Slats,
	rel.Slats,
	render.Slats,
	story.Slats,
}
var AllSignatures = []map[uint64]interface{}{
	core.Signatures,
	debug.Signatures,
	grammar.Signatures,
	literal.Signatures,
	list.Signatures,
	reader.Signatures,
	rel.Signatures,
	render.Signatures,
	rt.Signatures,
	story.Signatures,
}

// note: if it were ever to turn out this is the only reflection needed
// could replace with function callbacks that new() rather than just a generic r.Type
type TypeRegistry map[string]r.Type

var reg TypeRegistry

func Registry() TypeRegistry {
	if reg == nil {
		for _, slats := range AllSlats {
			reg.RegisterTypes(slats)
		}
	}
	return reg
}

func (reg *TypeRegistry) RegisterTypes(cmds []composer.Composer) (err error) {
	if *(reg) == nil {
		*(reg) = make(TypeRegistry)
	}
	for _, cmd := range cmds {
		if spec := cmd.Compose(); len(spec.Name) == 0 {
			e := errutil.Fmt("Missing type name %T", cmd)
			errutil.Append(err, e)
		} else if was, exists := (*reg)[spec.Name]; exists {
			e := errutil.Fmt("Duplicate type name %q now: %T, was: %s", spec.Name, cmd, was.String())
			errutil.Append(err, e)
			break
		} else {
			(*reg)[spec.Name] = r.TypeOf(cmd).Elem()
		}
	}
	return
}
