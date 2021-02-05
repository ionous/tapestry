package core

import (
	"strconv"
	"strings"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/term"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"github.com/ionous/errutil"
)

type Argument struct {
	Name string // argument name
	From Assignment
}

type Arguments struct {
	Args []*Argument
}

func (*Argument) Compose() composer.Spec {
	return composer.Spec{
		Spec:  " {name:variable_name}: {from:assignment}",
		Group: "patterns",
		Stub:  true,
	}
}

func (*Arguments) Compose() composer.Spec {
	return composer.Spec{
		Spec:  " {arguments%args+argument|comma-and}",
		Group: "patterns",
		Stub:  true,
	}
}

//
func (op *Arguments) Distill(run rt.Runtime, out *g.Record) (err error) {
	k := out.Kind()
	for _, arg := range op.Args {
		if name, e := getParamName(k, arg.Name); e != nil {
			err = errutil.Append(err, e)
		} else if val, e := GetAssignedValue(run, arg.From); e != nil {
			err = errutil.Append(err, e)
		} else if fin := k.FieldIndex(name); fin < 0 {
			e := errutil.New("unknown field", name)
			err = errutil.Append(err, e)
		} else if v, e := convertTerm(run, k.Field(fin), val); e != nil {
			err = errutil.Append(err, e)
		} else if e := out.SetNamedField(name, v); e != nil {
			err = errutil.Append(err, e)
			// fix: we have to set by name to handle traits
			// this doesnt mesh very well with the affine.Object conversion(
		}
	}
	return
}

func convertTerm(run rt.Runtime, f g.Field, v g.Value) (ret g.Value, err error) {
	ret = v                   // provisionally
	objectPrefix := "object=" // set by ObjectAsName
	if strings.HasPrefix(f.Type, objectPrefix) {
		kind := f.Type[len(objectPrefix):]
		switch aff := v.Affinity(); aff {
		// fix: templates ( other things? ) are still giving us object arguments
		// tbd: maybe cant be fixed until affine.Object is completely removed.
		case affine.Object:
			ret, err = term.ConvertObject(run, v, kind)

		case affine.Text:
			ret, err = term.ConvertName(run, v.String(), kind)
		}
	}
	return
}

// change a argument name ( which could be an index ) into a valid param name
// fix: this should happen at assembly time...
func getParamName(k *g.Kind, arg string) (ret string, err error) {
	if usesIndex := len(arg) > 1 && arg[:1] == "$"; !usesIndex {
		ret = lang.Breakcase(arg)
	} else if storedIdx, e := strconv.Atoi(arg[1:]); e != nil {
		err = errutil.New("couldnt parse index", arg)
	} else if i := storedIdx - 1; i < 0 || i >= k.NumField() {
		err = errutil.New("field", arg, "not found")
	} else {
		ret = k.Field(i).Name
	}
	return
}
