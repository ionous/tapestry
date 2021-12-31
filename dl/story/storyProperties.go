package story

import (
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"github.com/ionous/errutil"
)

func (op *PropertyDecl) ImportProperty(k *Importer, kind string) (err error) {
	return op.PropertyType.ImportPropertyType(k, kind, op.Property.Str)
}

func (op *PropertyType) ImportPropertyType(k *Importer, kind, prop string) (err error) {
	type propertyTypeImporter interface {
		ImportPropertyType(k *Importer, kind, prop string) error
	}
	if opt, ok := op.Value.(propertyTypeImporter); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", UnhandledSwap, op.Value))
	} else {
		err = opt.ImportPropertyType(k, kind, prop)
	}
	return
}

// inform gives these the name "<noun> condition"
// we could only do that with an after the fact reduction, and with some additional mdl data.
// ( ex. in case the same aspect is assigned twice, or twice at difference depths )
// for now the name of the field is the name of the aspect
func (op *PropertyAspect) ImportPropertyType(k *Importer, kind, prop string) (err error) {
	// record the existence of an aspect with the same name as the property
	k.Write(&eph.EphKinds{Kinds: kind, Contain: []eph.EphParams{eph.AspectParam(prop)}})
	return
}

// "{a number%number}, {some text%text}, or {a true/false value%bool}");
// bool properties become implicit aspects
func (op *PrimitiveType) ImportPropertyType(k *Importer, kind, aspect string) (err error) {
	if op.Str != PrimitiveType_Bool {
		primType := eph.Affinity{op.Str} // these are the same
		k.Write(&eph.EphKinds{Kinds: kind, Contain: []eph.EphParams{{Name: aspect, Affinity: primType}}})
	} else {
		// ex. innumerable, not innumerable, is innumerable
		k.AddImplicitAspect(aspect, kind,
			"not_"+aspect, // false first ( so that the default is the zero value )
			"is_"+aspect,
		)
		k.Write(&eph.EphKinds{Kinds: kind, Contain: []eph.EphParams{eph.AspectParam(aspect)}})
	}
	return
}

// number_list, text_list, record_type, record_list
func (op *ExtType) GetParameterType() (retType string, retAff eph.Affinity, err error) {
	if imp, ok := op.Value.(primTypeAffer); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", UnhandledSwap, op.Value))
	} else {
		retType, retAff = imp.GetTypeAffinity()
	}
	return
}

// number_list, text_list, record_type, record_list
func (op *ExtType) ImportPropertyType(k *Importer, kind, prop string) (err error) {
	if imp, ok := op.Value.(primTypeAffer); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", UnhandledSwap, op.Value))
	} else {
		primType, primAff := imp.GetTypeAffinity()
		k.Write(&eph.EphKinds{Kinds: kind, Contain: []eph.EphParams{
			{Name: prop, Affinity: primAff, Class: primType},
		}})
	}
	return
}

type primTypeAffer interface{ GetTypeAffinity() (string, eph.Affinity) }

func (op *NumberList) GetTypeAffinity() (string, eph.Affinity) {
	return "", eph.Affinity{eph.Affinity_NumList}
}

func (op *TextList) GetTypeAffinity() (string, eph.Affinity) {
	return "", eph.Affinity{eph.Affinity_TextList}
}

func (op *RecordType) GetTypeAffinity() (string, eph.Affinity) {
	return op.Kind.Str, eph.Affinity{eph.Affinity_Record}
}

func (op *RecordList) GetTypeAffinity() (string, eph.Affinity) {
	return op.Kind.Str, eph.Affinity{eph.Affinity_RecordList}
}
