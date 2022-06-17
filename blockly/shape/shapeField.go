package shape

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"github.com/ionous/errutil"
)

type shapeField struct {
	term     spec.TermSpec
	typeSpec *spec.TypeSpec
	slot     bconst.SlotRule
}

func (w *ShapeWriter) newFieldDef(term spec.TermSpec) (ret *shapeField, err error) {
	typeName := term.TypeName() // lookup spec
	if typeSpec, ok := w.Types[typeName]; !ok {
		err = errutil.New("missing named type", typeName)
	} else {
		var slot bconst.SlotRule
		if typeSpec.Spec.Choice == spec.UsesSpec_Slot_Opt {
			// inputType might be a statement_input stack, or a single ( maybe repeatable ) input
			// regardless, it only has the input, no special fields.
			slot = bconst.FindSlotRule(typeSpec.Name)
		}
		ret = &shapeField{term, typeSpec, slot}
	}
	return
}

func (fd *shapeField) name() string {
	return strings.ToUpper(fd.term.Field())
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

func (fd *shapeField) shadow() (ret string) {
	switch fd.termType() {
	case spec.UsesSpec_Num_Opt, spec.UsesSpec_Flow_Opt:
		ret = fd.typeSpec.Name
	}
	return
}

// handle our fake field for the leading label.
func (fd *shapeField) termType() (ret string) {
	if fd.typeSpec != nil {
		ret = fd.typeSpec.Spec.Choice
	}
	return
}
