package express

import (
	r "reflect"
	"strconv"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/render"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/template"
	"git.sr.ht/~ionous/tapestry/template/postfix"
	"git.sr.ht/~ionous/tapestry/template/types"
	"github.com/ionous/errutil"
)

// Express converts a postfix expression into commands.
func Convert(xs template.Expression) (ret interface{}, err error) {
	c := Converter{}
	return c.Convert(xs)
}

type Converter struct {
	stack cmdStack // the stack is empty initially, and we fill it with converted commands
	// ( to be used later by other commands )
	autoCounter int
}

func (c *Converter) Convert(xs template.Expression) (ret interface{}, err error) {
	if e := c.convert(xs); e != nil {
		err = e
	} else if op, e := c.stack.flush(); e != nil {
		err = e
	} else if on, ok := op.(dotName); ok {
		// if the entire template can be reduced to an dotName ( ex. `{.lantern}` )
		// then we treat it as a request for the friendly name of an object
		ret = on.getPrintedName()
	} else {
		ret = op
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
func (c *Converter) buildCompare(cmp core.Comparison) (err error) {
	if args, e := c.stack.pop(2); e != nil {
		err = e
	} else {
		var ptr r.Value
		a, b := unpackArg(args[0]), unpackArg(args[1])
		an, bn := a.String(), b.String() // here for debugging
		switch {
		case implements(a, b, typeNumEval):
			ptr = r.New(compareNum)
		case implements(a, b, typeTextEval):
			ptr = r.New(compareText)
		default:
			err = errutil.Fmt("unknown commands %v %v", an, bn)
		}
		if err == nil {
			cmp := r.ValueOf(cmp)
			args = []r.Value{a, cmp, b}
			if e := assignProps(ptr.Elem(), args); e != nil {
				err = e
			} else {
				c.stack.push(ptr)
			}
		}
	}
	return
}

func (c *Converter) buildSequence(cmd rt.TextEval, pAt *string, pParts *[]rt.TextEval, count int) (err error) {
	if args, e := c.stack.pop(count); e != nil {
		err = e
	} else {
		var parts []rt.TextEval
		for i, a := range args {
			a := unpackArg(a)
			if text, ok := a.Interface().(rt.TextEval); !ok {
				err = errutil.Fmt("couldn't convert sequence part %d to text", i)
				break
			} else {
				parts = append(parts, text)
			}
		}
		if err == nil {
			c.autoCounter++
			counter := "autoexp" + strconv.Itoa(c.autoCounter)
			// seq is part of cmd
			(*pParts) = parts
			(*pAt) = counter
			// after filling out the cmd, we push it for later processing
			c.buildOne(cmd)
		}
	}
	return
}

// build an command named in the export Slat
// names in templates are currently "mixedCase" rather than "underscore_case".
func (c *Converter) buildExport(name string, arity int) (err error) {
	// if its in the coreCache its a known command
	if a, ok := coreCache.get(name); !ok {
		// if its not, the user is probably calling a patter
		err = c.buildPattern(name, arity)
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

func (c *Converter) buildPattern(name string, arity int) (err error) {
	// pull the number of arguments needed ( already converted into commands )
	// args is a slice of package reflect r.Value(s)
	if args, e := c.stack.pop(arity); e != nil {
		err = e
	} else {
		values := make([]render.RenderEval, len(args))
		for i, arg := range args {
			values[i] = unpackPatternArg(arg)
		}
		if err == nil {
			c.buildOne(&render.RenderPattern{
				PatternName: name,
				Render:      values,
			})
		}
	}
	return
}

func unpackPatternArg(arg r.Value) render.RenderEval {
	var out rt.Assignment
	switch arg := arg.Interface().(type) {
	default:
		panic(errutil.Fmt("unknown argument type %T", arg))
	case dotName:
		return arg.getNamedValue()
	case *render.RenderPattern:
		return arg
	case rt.BoolEval:
		out = &assign.FromBool{Value: arg}
	case rt.NumEval:
		out = &assign.FromNum{Value: arg}
	case rt.TextEval:
		out = &assign.FromText{Value: arg}
	case rt.RecordEval:
		out = &assign.FromRecord{Value: arg}
	case rt.NumListEval:
		out = &assign.FromNumList{Value: arg}
	case rt.TextListEval:
		out = &assign.FromTextList{Value: arg}
	case rt.RecordListEval:
		out = &assign.FromRecordList{Value: arg}
	}
	// fall through handling for assignments
	return &render.RenderValue{Value: out}
}

// an eval h
func (c *Converter) buildUnless(cmd interface{}, arity int) (err error) {
	if args, e := c.stack.pop(arity); e != nil {
		err = e
	} else if len(args) > 0 {
		arg := unpackArg(args[0])
		if a, ok := arg.Interface().(rt.BoolEval); !ok {
			err = errutil.New("argument is not a bool", arg.Type().String())
		} else {
			args[0] = r.ValueOf(&core.Not{Test: a}) // rewrite the arg.
			c.stack.push(args...)                   //
			err = c.buildCommand(cmd, arity)
		}
	}
	return
}

func (c *Converter) buildSpan(arity int) (err error) {
	if args, e := c.stack.pop(arity); e != nil {
		err = e
	} else {
		var txts []rt.TextEval
		for _, el := range args {
			switch el := el.Interface().(type) {
			// in a list of text evaluations,
			// for example maybe "{.bennie} and the {.jets}"
			// single occurrences of dotted names are treated as requests for a friendly name
			case dotName:
				txts = append(txts, el.getPrintedName())
			case rt.TextEval:
				txts = append(txts, el)
			default:
				e := errutil.Fmt("argument %T is not a text eval", el)
				err = errutil.Append(err, e)
			}
		}
		if err == nil {
			c.buildOne(&core.Join{Parts: txts})
		}
	}
	return
}

// convert the passed postfix template function into commands.
func (c *Converter) addFunction(fn postfix.Function) (err error) {
	switch fn := fn.(type) {
	case types.Quote:
		txt := fn.Value()
		c.buildOne(T(txt))

	case types.Number:
		num := fn.Value()
		c.buildOne(F(num))

	case types.Bool:
		b := fn.Value()
		c.buildOne(B(b))

	case types.Command: // see decode
		err = c.buildExport(fn.CommandName, fn.CommandArity)

	case types.Reference:
		// fields are an array of strings .a.b.c
		if fields := fn.Value(); len(fields) == 0 {
			err = errutil.New("empty reference")
		} else {
			// fix: this should add ephemera that there's an object of name
			// fix: can this add ephemera that there's a local of name?
			if firstField := fields[0]; len(fields) == 1 {
				// we dont know yet how { .name.... } is being used:
				// - a command arg, so the desired type is known.
				// - a pattern arg, so the desired type isn't known.
				// - a request to print an object name
				//
				// the name itself could refer to:
				// - the name of an object,
				// - the name of a pattern parameter,
				// - a loop counter,
				// - etc.
				c.buildOne(dotName(firstField))
			} else {
				// a chain of dots indicates one or more fields of a record
				// ex.  .object.fieldContainingAnRecord.otherField
				dot := make([]assign.Dot, len(fields)-1)
				for i, field := range fields[1:] {
					dot[i] = &assign.AtField{Field: T(field)}
				}
				c.buildOne(&render.RenderRef{Name: T(firstField), Dot: dot})
			}
		}

	case types.Builtin:
		switch k := fn.Type; k {
		case types.IfStatement:
			// it would be nice if this could be choose text or choose number based on context
			// choose scalar might simplify things....
			err = c.buildCommand(&core.ChooseText{}, fn.ParameterCount)
		case types.UnlessStatement:
			err = c.buildUnless(&core.ChooseText{}, fn.ParameterCount)

		case types.Stopping:
			var seq core.CallTerminal
			err = c.buildSequence(&seq, &seq.Name, &seq.Parts, fn.ParameterCount)
		case types.Shuffle:
			var seq core.CallShuffle
			err = c.buildSequence(&seq, &seq.Name, &seq.Parts, fn.ParameterCount)
		case types.Cycle:
			var seq core.CallCycle
			err = c.buildSequence(&seq, &seq.Name, &seq.Parts, fn.ParameterCount)
		case types.Span:
			err = c.buildSpan(fn.ParameterCount)

		default:
			// fix? span is supposed to join text sections.... but there were no tests or examples in the og code.
			err = errutil.New("unhandled builtin", k.String())
		}

	case types.Operator:
		switch fn {
		case types.MUL:
			err = c.buildTwo(&core.MultiplyValue{})
		case types.QUO:
			err = c.buildTwo(&core.DivideValue{})
		case types.REM:
			err = c.buildTwo(&core.ModValue{})
		case types.ADD:
			err = c.buildTwo(&core.AddValue{})
		case types.SUB:
			err = c.buildTwo(&core.SubtractValue{})

		case types.EQL:
			err = c.buildCompare(core.C_Comparison_EqualTo)
		case types.NEQ:
			err = c.buildCompare(core.C_Comparison_OtherThan)
		case types.LSS:
			err = c.buildCompare(core.C_Comparison_LessThan)
		case types.LEQ:
			err = c.buildCompare(core.C_Comparison_AtMost)
		case types.GTR:
			err = c.buildCompare(core.C_Comparison_GreaterThan)
		case types.GEQ:
			err = c.buildCompare(core.C_Comparison_AtLeast)

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
