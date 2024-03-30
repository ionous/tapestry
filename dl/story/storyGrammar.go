package story

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// Execute - called by the macro runtime during weave.
func (op *DefineAlias) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefineAlias) Weave(cat *weave.Catalog) (err error) {
	return cat.Schedule(weaver.ValuePhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if name, e := safe.GetText(run, op.NounName); e != nil {
			err = e
		} else if names, e := safe.GetTextList(run, op.Names); e != nil {
			err = e
		} else if n, e := run.GetField(meta.ObjectId, name.String()); e != nil {
			err = e
		} else {
			n := n.String()
			for _, a := range names.Strings() {
				if a := inflect.Normalize(a); len(a) > 0 {
					if e := w.AddNounName(n, a, -1); e != nil {
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
func (op *DefineLeadingGrammar) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// an ugly way to ensure that grammar ( and therefore the runtime )
// isnt dependent on story / weave
func (op *DefineLeadingGrammar) Weave(cat *weave.Catalog) (err error) {
	// jump/skip/hop	{"Directive:scans:":[["jump","skip","hop"],[{"As:":"jumping"}]]}
	name := strings.Join(op.Lede, "/")
	words := &grammar.Words{Words: op.Lede}
	scans := append([]grammar.ScannerMaker{words}, op.Scans...)
	return cat.Schedule(weaver.ValuePhase, func(w weaver.Weaves, run rt.Runtime) error {
		return w.AddGrammar(name, &grammar.Directive{
			Name:   name,
			Series: scans,
		})
	})
}

// Execute - called by the macro runtime during weave.
func (op *DefineNamedGrammar) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// an ugly way to ensure that grammar ( and therefore the runtime )
// isnt dependent on story / weave
func (op *DefineNamedGrammar) Weave(cat *weave.Catalog) (err error) {
	name := inflect.Normalize(op.Name)
	return cat.Schedule(weaver.ValuePhase, func(w weaver.Weaves, run rt.Runtime) error {
		return w.AddGrammar(name, &grammar.Directive{
			Name:   name,
			Series: op.Scans,
		})
	})
}
