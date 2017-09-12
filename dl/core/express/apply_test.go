package express

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec/ops"
	r "reflect"
	"testing"
)

type TestThe struct {
	Obj rt.ObjectEval
}

func (*TestThe) GetText(run rt.Runtime) (ret string, err error) {
	err = errutil.New("not implemented")
	return
}

func xTestApply(t *testing.T) {
	const (
		partStr = "{status.score}"
		cmdStr  = "{go TestThe example}"
		// ifElseStr    = "{if x}{status.score}{else}{story.turnCount}{endif}"
	)
	classes := make(unique.Types)
	cmds := ops.NewOps(classes)

	unique.PanicBlocks(cmds,
		(*std.Commands)(nil), // for Render
		(*core.Commands)(nil))

	unique.PanicTypes(cmds,
		(*TestThe)(nil))

	t.Run("parts", func(t *testing.T) {
		testEqual(t, partsFn(),
			templatize(t, partStr, cmds))
	})
	// t.Run("cmds", func(t *testing.T) {
	// testEqual(t, cmdsFn(),
	// 	templatize(t, cmdStr, cmds))
	// })
}

func templatize(t *testing.T, s string, cmds *ops.Ops) (ret rt.TextEval) {
	xf := Xform{cmds: cmds}
	rtype := r.TypeOf((*rt.TextEval)(nil)).Elem()
	if r, e := xf.TransformValue(s, rtype); e != nil {
		t.Fatal(e)
	} else {
		ret = r.(rt.TextEval)
	}
	return
}

func partsFn() rt.TextEval {
	return &core.Buffer{[]rt.Execute{
		&core.Say{
			&core.Get{
				Obj:  &core.GetAt{"status"},
				Prop: "score",
			},
		},
		&core.Say{
			&core.Text{"/"},
		},
		&core.Say{
			&core.Get{
				Obj:  &core.GetAt{"story"},
				Prop: "turnCount",
			},
		},
	},
	}
}

func cmdsFn() rt.TextEval {
	return &core.Buffer{[]rt.Execute{
		&core.Say{
			Text: &TestThe{
				&core.Object{Name: "example"},
			},
		},
	},
	}
}
