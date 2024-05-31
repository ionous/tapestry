package object

import (
	"log"

	"git.sr.ht/~ionous/tapestry/dl/assign/dot"
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func MakeDot(path ...any) (ret []Dot) {
	if cnt := len(path); cnt > 0 {
		out := make([]Dot, len(path))
		for i, p := range path {
			switch el := p.(type) {
			case string:
				out[i] = &AtField{Field: literal.T(el)}
			case int:
				out[i] = &AtIndex{Index: literal.I(el)}
			case Dot:
				out[i] = el
			default:
				log.Panicf("expected an int or string element; got %T", el)
			}
		}
		ret = out
	}
	return
}

func Object(name string, path ...any) *ObjectDot {
	return &ObjectDot{
		Name: literal.T(name),
		Dot:  MakeDot(path...),
	}
}

// generate a statement which extracts a variable's value.
// path can include strings ( for reading from records ) or integers ( for reading from lists )
func Variable(name string, path ...any) *VariableDot {
	return &VariableDot{
		Name: literal.T(name),
		Dot:  MakeDot(path...),
	}
}

// return a dot.Index
func (op *AtIndex) Resolve(run rt.Runtime) (ret rt.Dotted, err error) {
	if idx, e := safe.GetNum(run, op.Index); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = dot.Index(idx.Int() - 1)
	}
	return
}

// return a dot.Field
func (op *AtField) Resolve(run rt.Runtime) (ret rt.Dotted, err error) {
	if field, e := safe.GetText(run, op.Field); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = dot.Field(field.String())
	}
	return
}
