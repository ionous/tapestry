package story

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/ephemera"
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

type variableDecl struct {
	name, typeName ephemera.Named
	affinity       string
}

func (op *VariableDecl) ImportVariable(k *Importer, cat string) (ret variableDecl, err error) {
	if n, e := NewVariableName(k, op.Name, cat); e != nil {
		err = e
	} else if t, aff, e := op.Type.ImportVariableType(k); e != nil {
		err = e
	} else {
		ret = variableDecl{n, t, aff}
	}
	return
}

// primitive type, object type, or ext
func (op *VariableType) ImportVariableType(k *Importer) (retType ephemera.Named, retAff string, err error) {
	type variableTypeImporter interface {
		ImportVariableType(*Importer) (ephemera.Named, string, error)
	}
	if opt, ok := op.Opt.(variableTypeImporter); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", UnhandledSwap, op.Opt))
	} else {
		retType, retAff, err = opt.ImportVariableType(k)
	}
	return
}

func (op *ObjectType) ImportVariableType(k *Importer) (retType ephemera.Named, retAff string, err error) {
	retType, err = NewSingularKind(k, op.Kind)
	retAff = affine.Object.String()
	return
}

func (op *PrimitiveType) ImportPrimType(k *Importer) (ret string, err error) {
	if str, ok := composer.FindChoice(op, op.Str); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w %q", InvalidValue, op.Str))
	} else {
		ret = str
	}
	return
}

// returns one of the evalType(s) as a "Named" value --
// we return a name to normalize references to object kinds which are also used as variables
func (op *PrimitiveType) ImportVariableType(k *Importer) (retType ephemera.Named, retAff string, err error) {
	// fix -- shouldnt this be a different type ??
	// ie. we should be able to use FindChoie here.
	var namedType string
	switch str := op.Str; str {
	case PrimitiveType_Number:
		namedType = "number_eval"
	case PrimitiveType_Text:
		namedType = "text_eval"
	case PrimitiveType_Bool:
		namedType = "bool_eval"
	default:
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", InvalidValue, str))
	}
	if err == nil {
		retType = k.NewName(namedType, tables.NAMED_TYPE, op.At.String())
	}
	return
}
