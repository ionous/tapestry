package story

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/action"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/ionous/errutil"
)

// marker interface
type ActionParam interface {
	GetParam(rt.Runtime) (string, error)
}

// Execute - called by the macro runtime during weave.
func (op *StoryAlias) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *StoryAlias) Weave(cat *weave.Catalog) (err error) {
	return cat.Schedule(weave.RequireAll, func(w *weave.Weaver) (err error) {
		name := lang.Normalize(op.AsNoun)
		if n, e := w.GetClosestNoun(name); e != nil {
			err = e
		} else {
			pen := w.Pin()
			for _, a := range op.Names {
				if a := lang.Normalize(a); len(a) > 0 {
					if e := pen.AddName(n, a, -1); e != nil {
						err = e
						break
					}
				}
			}
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *StoryDirective) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// an ugly way to ensure that grammar ( and therefore the runtime )
// isnt dependent on story / weave
func (op *StoryDirective) Weave(cat *weave.Catalog) (err error) {
	// jump/skip/hop	{"Directive:scans:":[["jump","skip","hop"],[{"As:":"jumping"}]]}
	name := strings.Join(op.Lede, "/")
	return cat.Schedule(weave.RequireRules, func(w *weave.Weaver) error {
		return w.Pin().AddGrammar(name, &grammar.Directive{
			Lede:  op.Lede,
			Scans: op.Scans,
		})
	})
}

// scheduled by importStory
// verifies that a pattern exists for this action
func importAction(cat *weave.Catalog, op *grammar.Action) error {
	return cat.Schedule(weave.RequirePatterns, func(w *weave.Weaver) error {
		act := lang.Normalize(op.Action) // todo: a simpler way of handling references
		return w.Pin().ExtendPattern(mdl.NewPatternBuilder(act).Pattern)
	})
}

func (op *StoryAction) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (*ActionParamNothing) GetParam(rt.Runtime) (_ string, _ error) {
	return
}

func (op *ActionParamNoun) GetParam(run rt.Runtime) (ret string, err error) {
	if name, e := op.KindName.GetText(run); e != nil {
		err = e
	} else {
		ret = lang.Normalize(name.String())
	}
	return
}

func (op *StoryAction) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequireDependencies, func(w *weave.Weaver) (err error) {
		if act, e := safe.GetText(w, op.Action); e != nil {
			err = e
		} else {
			act := lang.Normalize(act.String())
			for evt := action.FirstEvent; evt < action.NumEvents; evt++ {
				pb := mdl.NewPatternSubtype(evt.Name(act), evt.Kind())
				if name, e := getParamName(w.Runtime, op.FirstNoun); e != nil {
					err = e
					break
				} else {
					addParam(pb, action.Noun, name)
				}
				if name, e := getParamName(w.Runtime, op.SecondNoun); e != nil {
					err = e
					break
				} else {
					addParam(pb, action.OtherNoun, name)
				}
				addParam(pb, action.Actor, action.Actor.String())
				// other information for events:
				addParam(pb, action.Target, nil)
				addParam(pb, action.CurrentTarget, nil)
				addParam(pb, action.Interupt, affine.Bool)
				addParam(pb, action.Cancel, affine.Bool)
				// write the pattern
				if e := cat.Schedule(weave.RequirePatterns, func(w *weave.Weaver) error {
					return w.Pin().AddPattern(pb.Pattern)
				}); e != nil {
					err = e
					break
				}
			}
		}
		return
	})
}

func addParam(pb *mdl.PatternBuilder, f action.Field, affineOrCls any) {
	var clsname string
	aff := affine.Text
	if a, ok := affineOrCls.(affine.Affinity); ok {
		aff = a
	} else if affineOrCls != nil {
		clsname = affineOrCls.(string)
	}
	pb.AddParam(mdl.FieldInfo{
		Name:     f.String(),
		Affinity: aff,
		Class:    clsname,
	})
}

func getParamName(run rt.Runtime, op ActionParam) (ret string, err error) {
	if op == nil {
		err = errutil.New("invalid")
	} else if txt, e := op.GetParam(run); e != nil {
		err = e
	} else {
		ret = txt
	}
	return
}
