package story

import (
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"github.com/ionous/errutil"
)

func (op *VariableDecl) GetParam() (ret eph.EphParams, err error) {
	if t, aff, e := op.Type.GetParameterType(); e != nil {
		err = e
	} else {
		ret = eph.EphParams{
			Name:     op.Name.String(),
			Affinity: aff,
			Class:    t,
		}
	}
	return
}

// primitive type, object type, or ext
func (op *VariableType) GetParameterType() (retType string, retAff eph.Affinity, err error) {
	if opt, ok := op.Value.(variableTypeImporter); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", UnhandledSwap, op.Value))
	} else {
		retType, retAff, err = opt.GetParameterType()
	}
	return
}

func (op *ObjectType) GetParameterType() (retType string, retAff eph.Affinity, err error) {
	retType = op.Kind.String()
	retAff = eph.Affinity{eph.Affinity_Text}
	return
}

// returns one of the evalType(s) as a "Named" value --
// we return a name to normalize references to object kinds which are also used as variables
func (op *PrimitiveType) GetParameterType() (retType string, retAff eph.Affinity, err error) {
	switch str := op.Str; str {
	case PrimitiveType_Number, PrimitiveType_Text, PrimitiveType_Bool:
		retAff = eph.Affinity{str}
	default:
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", InvalidValue, str))
	}
	return
}

type variableTypeImporter interface {
	GetParameterType() (string, eph.Affinity, error)
}
