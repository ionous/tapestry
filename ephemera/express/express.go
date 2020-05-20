package express

import (
	r "reflect"
	"strconv"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/template"
	"github.com/ionous/iffy/template/postfix"
	"github.com/ionous/iffy/template/types"
)

// Express converts a postfix expression into iffy commands.
func Convert(xs template.Expression) (ret interface{}, err error) {
	c := Converter{}
	return c.Convert(xs)
}

type Converter struct {
	stack rstack // the stack is empty initially, and we fill it with converted commands
	// ( to be used later by other commands )
	AutoCounter int
}

func (c *Converter) Convert(xs template.Expression) (ret interface{}, err error) {
	if e := c.convert(xs); e != nil {
		err = e
	} else {
		ret, err = c.stack.flush()
	}
	return
}

func (c *Converter) convert(xs template.Expression) (err error) {
	for _, fn := range xs {
		if e := c.addFunction(fn); e != nil {
			err = e
			break
		}
	}
	return
}

func (c *Converter) buildOne(cmd interface{}) {
	c.stack.push(r.ValueOf(cmd))
}

func (c *Converter) buildTwo(cmd interface{}) (err error) {
	return c.buildCommand(cmd, 2)
}

func (c *Converter) buildCommand(cmd interface{}, arity int) (err error) {
	if args, e := c.stack.pop(arity); e != nil {
		err = e
	} else {
		ptr := r.ValueOf(cmd)
		if e := assignProps(ptr.Elem(), args); e != nil {
			err = e
		} else {
			c.stack.push(ptr)
		}
	}
	return
}

// fix? this is where a Scalar value could come in handy.
func (c *Converter) buildCompare(cmp core.Comparator) (err error) {
	if args, e := c.stack.pop(2); e != nil {
		err = e
	} else {
		var ptr r.Value
		switch a, b := args[0], args[1]; {
		case implements(a, b, typeNumEval):
			ptr = r.New(compareNum)
		case implements(a, b, typeTextEval):
			ptr = r.New(compareText)
		default:
			err = errutil.New("unknown commands")
		}
		if err == nil {
			cmp := r.ValueOf(cmp)
			args = []r.Value{args[0], cmp, args[1]}
			if e := assignProps(ptr.Elem(), args); e != nil {
				err = e
			} else {
				c.stack.push(ptr)
			}
		}
	}
	return
}

func (c *Converter) buildSequence(cmd rt.TextEval, seq *core.Sequence, count int) (err error) {
	if args, e := c.stack.pop(count); e != nil {
		err = e
	} else {
		var parts []rt.TextEval
		for i, a := range args {
			if text, ok := a.Interface().(rt.TextEval); !ok {
				err = errutil.Fmt("couldn't convert sequence part %d to text", i)
				break
			} else {
				parts = append(parts, text)
			}
		}
		if err == nil {
			c.AutoCounter++
			counter := "autoexp" + strconv.Itoa(c.AutoCounter)
			// seq is part of cmd
			seq.Parts = parts
			seq.Seq = counter
			// after filling out the cmd, we push it for later processing
			c.buildOne(cmd)
		}
	}
	return
}

// build an command named in the export Slat
// names in templates are currently "mixedCase" rather than "underscore_case".
func (c *Converter) buildExport(name string, arity int) (err error) {
	if a, ok := exportsCache.get(name); !ok {
		err = errutil.New("unknown command", name, arity)
	} else if args, e := c.stack.pop(arity); e != nil {
		err = e
	} else {
		rtype := r.TypeOf(a).Elem()
		ptr := r.New(rtype)
		if e := assignProps(ptr.Elem(), args); e != nil {
			err = e
		} else {
			c.stack.push(ptr)
		}
	}
	return
}

func (c *Converter) buildUnless(cmd interface{}, arity int) (err error) {
	if args, e := c.stack.pop(arity); e != nil {
		err = e
	} else if len(args) > 0 {
		if a, ok := args[0].Interface().(rt.BoolEval); !ok {
			err = errutil.New("argument is not a bool")
		} else {
			args[0] = r.ValueOf(&core.IsNot{a}) // rewrite the arg.
			c.stack.push(args...)               //
			err = c.buildCommand(cmd, arity)
		}
	}
	return
}

// convert the passed postfix template function into iffy commands.
func (c *Converter) addFunction(fn postfix.Function) (err error) {
	switch fn := fn.(type) {
	case types.Quote:
		txt := fn.Value()
		c.buildOne(&core.Text{txt})

	case types.Number:
		num := fn.Value()
		c.buildOne(&core.Number{num})

	case types.Bool:
		b := fn.Value()
		c.buildOne(&core.Bool{b})

	case types.Command: // see decode
		err = c.buildExport(fn.CommandName, fn.CommandArity)

	case types.Reference:
		// fields are an array of strings a.b.c
		if fields := fn.Value(); len(fields) == 0 {
			err = errutil.New("empty reference")
		} else {
			// build a chain of GetFields
			// to start: we either want the object named "text"
			// or, we want the object name that's stored in the local variable called "text"
			var op rt.TextEval
			if name := fields[0]; lang.IsCapitalized(name) {
				// fix: this should add ephemera that there's an object of name
				op = &core.Text{name}
			} else {
				// fix: can this add ephemera that there's a local of name?
				op = &core.GetVar{name}
			}
			// a.b: from the named object a, we want its field b
			// a.b.c: after getting the object name in field b, get that object's field c
			for _, field := range fields[1:] {
				op = &core.GetField{op, &core.Text{field}}
			}
			// the whole chain becomes a single "function"
			c.buildOne(op)
		}

	case types.Builtin:
		switch k := fn.Type; k {
		case types.IfStatement:
			// it would be nice if this could be choose text or choose number based on context
			// choose scalar might simplify things....
			err = c.buildCommand(&core.ChooseText{}, fn.ParameterCount)

		case types.UnlessStatement:
			err = c.buildUnless(&core.ChooseText{}, fn.ParameterCount)

		// case types.Span:
		// 	err = c.buildCommand(&core.Join{}, fn.ParameterCount)

		case types.Stopping:
			var seq core.StoppingText
			err = c.buildSequence(&seq, &seq.Sequence, fn.ParameterCount)
		case types.Shuffle:
			var seq core.ShuffleText
			err = c.buildSequence(&seq, &seq.Sequence, fn.ParameterCount)
		case types.Cycle:
			var seq core.CycleText
			err = c.buildSequence(&seq, &seq.Sequence, fn.ParameterCount)
		}

	case types.Operator:
		switch fn {
		case types.MUL:
			err = c.buildTwo(&core.ProductOf{})
		case types.QUO:
			err = c.buildTwo(&core.QuotientOf{})
		case types.REM:
			err = c.buildTwo(&core.RemainderOf{})
		case types.ADD:
			err = c.buildTwo(&core.SumOf{})
		case types.SUB:
			err = c.buildTwo(&core.DiffOf{})

		case types.EQL:
			err = c.buildCompare(&core.EqualTo{})
		case types.NEQ:
			err = c.buildCompare(&core.NotEqualTo{})
		case types.LSS:
			err = c.buildCompare(&core.LessThan{})
		case types.LEQ:
			err = c.buildCompare(&core.LessOrEqual{})
		case types.GTR:
			err = c.buildCompare(&core.GreaterThan{})
		case types.GEQ:
			err = c.buildCompare(&core.GreaterOrEqual{})

		case types.LAND:
			err = c.buildTwo(&core.AllTrue{})
		case types.LOR:
			err = c.buildTwo(&core.AnyTrue{})
		default:
			err = errutil.Fmt("unknown operator %s", fn)
		}

	default:
		err = errutil.Fmt("unknown function %T", fn)
	}
	return
}