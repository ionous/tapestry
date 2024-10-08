package express

import (
	r "reflect"

	"github.com/ionous/errutil"
)

func assignProps(out r.Value, args []r.Value) (err error) {
	outType := out.Type()
	WalkProperties(outType, func(f *r.StructField, path []int) (done bool) {
		if len(args) <= 0 {
			done = true
		} else {
			field := out.FieldByIndex(path)
			if f.Type.Kind() != r.Slice {
				if arg, rest := popArg(f.Type, args); !arg.IsValid() {
					err = errutil.Fmt("cant assign %s to field %s.%s (%s)",
						args[0].Type(), outType, f.Name, f.Type)
				} else {
					field.Set(arg)
					args = rest
				}
			} else {
				// when assigning to a slice, eat as many elements as possible.
				// it makes having slices as the last element of a command a good idea.
				slice, elType := field, f.Type.Elem()
				for len(args) > 0 {
					if arg, rest := popArg(elType, args); !arg.IsValid() {
						break
					} else {
						slice = r.Append(slice, arg)
						args = rest
					}
				}
				field.Set(slice)
			}
		}
		return done || err != nil // returns "done" when there is an error.
	})
	if err == nil && len(args) > 0 {
		err = errutil.New("unable to consume all args in", outType.Name())
	}
	return
}

func unpackArg(arg r.Value) r.Value {
	if on, ok := arg.Interface().(dotName); ok {
		arg = r.ValueOf(on.getNamedValue())
	}
	return arg
}

func popArg(elType r.Type, args []r.Value) (ret r.Value, rest []r.Value) {
	arg := unpackArg(args[0])
	if argType := arg.Type(); argType.AssignableTo(elType) {
		ret, rest = arg, args[1:] // pop
	}
	return
}
