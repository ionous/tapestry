package shortcut

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestTokens(t *testing.T) {
	if e := match("@var.5", AtSign, "var", 5); e != nil {
		t.Fatal(e)
	} else if e := match("#obj.field", HashMark, "obj", "field"); e != nil {
		t.Fatal(e)
	} else if e := match("#`i have spaces`.ok", HashMark, "i have spaces", "ok"); e != nil {
		t.Fatal(e)
	} else if e := match("@i_have__spaces", AtSign, "i have spaces"); e != nil {
		t.Fatal(e)
	} else if e := match("@var.5.6.lsat", AtSign, "var", 5, 6, "lsat"); e != nil {
		t.Fatal(e)
	} else {
		var short NotShort
		if e := match("@@ignore", AtSign); !errors.As(e, &short) || short != 1 {
			t.Fatal(e)
		} else if e := match("##ignore", HashMark); !errors.As(e, &short) || short != 1 {
			t.Fatal(e)
		} else if e := match("oh so normal"); !errors.As(e, &short) || short != 0 {
			t.Fatal(e)
		}
	}
}
func match(str string, want ...any) (err error) {
	if e := Tokenize(str, &matcher{want: want}); e != nil {
		err = fmt.Errorf("%q failed %w", str, e)
	}
	return
}

type matcher struct {
	want []any
	ofs  int
}

func (a *matcher) Decoded(t Type, v any) (err error) {
	if q, ok := v.(rune); ok {
		log.Printf("decoded %q", q)
	} else {
		log.Printf("decoded %v", v)
	}
	want := a.want[a.ofs]
	if !reflect.DeepEqual(v, want) {
		err = fmt.Errorf("wanted %v got %v at %d", want, v, a.ofs)
	} else {
		a.ofs++
	}
	return
}
