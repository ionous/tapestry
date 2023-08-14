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
	GetParamNames(rt.Runtime) (one, other string, err error)
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

func (*ActionParamNothing) GetParamNames(rt.Runtime) (_, _ string, _ error) {
	return
}

func (op *ActionParamNoun) GetParamNames(run rt.Runtime) (one, other string, err error) {
	if kind, e := safe.GetText(run, op.KindName); e != nil {
		err = e
	} else if otherKind, e := safe.GetText(run, op.OtherKindName); e != nil {
		err = e
	} else {
		one, other = lang.Normalize(kind.String()),
			lang.Normalize(otherKind.String())
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
				if params := op.Params; params == nil {
					err = errutil.New("invalid param names")
					break
				} else if one, other, e := params.GetParamNames(w); e != nil {
					err = e
					break
				} else {
					// for now we add both, even if they are blank.
					// fix: what about a separate event object;
					// then we only have to declare it once.
					// could put it on the stack during event processing.
					addParam(pb, action.Noun, one)
					addParam(pb, action.OtherNoun, other)
				}
				addParam(pb, action.Actor, action.Actor.String())
				// other information for events:
				addParam(pb, action.Target, nil)
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
