package encode

import (
	"errors"
	"fmt"
	r "reflect"

	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/lang/walk"
)

func WriteEnum(out *FlowBuilder, f walk.Field, instr string) (err error) {
	c := r.New(f.Type()).Interface() // ugh. fix.
	if n, ok := c.(composer.Composer); !ok {
		err = fmt.Errorf("is %s not a generated type?", f.Type())
	} else if outstr, ok := xformString(n.Compose(), instr); !ok {
		err = errors.New("invalid string")
	} else {
		err = out.WriteField(f, outstr)
	}
	return
}

func xformString(spec composer.Spec, str string) (ret string, okay bool) {
	if res, i := spec.IndexOfChoice(str); i >= 0 {
		ret, okay = res, true
	} else if spec.OpenStrings {
		ret, okay = str, true
	}
	return
}
