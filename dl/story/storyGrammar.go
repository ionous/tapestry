package story

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

// Execute - called by the macro runtime during weave.
func (op *DefineAlias) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefineAlias) Weave(cat *weave.Catalog) (err error) {
	return cat.Schedule(weave.RequireAll, func(w *weave.Weaver) (err error) {
		if name, e := safe.GetText(w, op.NounName); e != nil {
			err = e
		} else if names, e := safe.GetTextList(w, op.Names); e != nil {
			err = e
		} else if n, e := w.GetClosestNoun(lang.Normalize(name.String())); e != nil {
			err = e
		} else {
			pen := w.Pin()
			for _, a := range names.Strings() {
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
	return cat.Schedule(weave.RequireRules, func(w *weave.Weaver) error {
		return w.Pin().AddGrammar(name, &grammar.Directive{
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
	name := lang.Normalize(op.Name)
	return cat.Schedule(weave.RequireRules, func(w *weave.Weaver) error {
		return w.Pin().AddGrammar(name, &grammar.Directive{
			Name:   name,
			Series: op.Scans,
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
