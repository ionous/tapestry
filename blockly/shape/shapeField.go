package shape

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"github.com/ionous/errutil"
)

type shapeField struct {
	term typeinfo.Term
	slot bconst.SlotRule
}

func (w *ShapeWriter) newFieldDef(term typeinfo.Term) (ret *shapeField, err error) {
	var slot bconst.SlotRule
	if term.Type == nil {
		// could happen with private fields
		err = errutil.New("missing term type", term.Name)
	} else {
		if t, ok := term.Type.(*typeinfo.Slot); ok {
			// inputType might be a statement_input stack, or a single ( maybe repeatable ) input
			// regardless, it only has the input, no special fields.
			slot = bconst.MakeSlotRule(t)
		}
		ret = &shapeField{term, slot}
	}
	return
}

func (fd *shapeField) name() string {
	return strings.ToUpper(fd.term.Name)
}

// will we need a label for required anonymous terms?
// maybe at least for the mui?
func (fd *shapeField) blocklyLabel() (ret string) {
	if !fd.term.IsAnonymous() {
		ret = strings.Join(strings.Split(fd.term.Label, "_"), " ")
	} /*else {
		ret = fd.term.Name
	}*/
	return
}

// // handle our fake field for the leading label.
// func (fd *shapeField) termType() (ret string) {
// 	if fd.typeSpec != nil {
// 		ret = fd.typeSpec.Spec.Choice
// 	}
// 	return
// }
