package jsonexp_test

import (
	"encoding/json"
	r "reflect"
	"testing"

	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/ephemera/debug"
	"git.sr.ht/~ionous/iffy/ephemera/story"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

func TestDetails(t *testing.T) {
	n := Ctx{Name: ""}
	src := debug.FactorialStory
	// save the story
	if b, e := src.MarshalDetailed(&n); e != nil {
		t.Fatal(e)
	} else {
		// load the story into a map just for fun
		m := make(map[string]interface{})
		if e := json.Unmarshal(b, &m); e != nil {
			t.Fatal(e)
		} else {
			// read the story from raw bytes
			var dst story.Story
			if e := dst.UnmarshalDetailed(&n, b); e != nil {
				t.Fatal(e)
			} else if diff := pretty.Diff(dst, *src); len(diff) > 0 {
				t.Fatal(diff)
				// did everything check out?
			}
		}
	}
}

type Ctx struct{ Name string }

func (ctx *Ctx) Source() string {
	return ctx.Name
}

func (ctx *Ctx) Finalize(ptr interface{}) (interface{}, error) {
	return ptr, nil
}

func (ctx *Ctx) NewType(s, t string) (ret interface{}, err error) {
	if s := newType(t, story.Slats); s != nil {
		ret = s
	} else if s := newType(t, core.Slats); s != nil {
		ret = s
	} else {
		err = errutil.New("unknown type", t)
	}
	return
}

func newType(t string, cs []composer.Composer) (ret interface{}) {
	for _, c := range cs {
		if c.Compose().Name == t {
			ret = r.New(r.TypeOf(c).Elem()).Interface()
			break
		}
	}
	return
}
