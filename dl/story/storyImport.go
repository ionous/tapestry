package story

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/weave"
	"github.com/ionous/errutil"

	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// Execute - called by the macro runtime during weave.
func (op *DefineTraits) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// (the) colors are red, blue, or green.
func (op *DefineTraits) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequirePlurals, func(w *weave.Weaver) (err error) {
		if aspect, e := safe.GetText(w, op.Aspect); e != nil {
			err = e
		} else if traits, e := safe.GetTextList(w, op.Traits); e != nil {
			err = e
		} else {
			aspect, traits := lang.Normalize(aspect.String()), traits.Strings()
			for i, t := range traits {
				traits[i] = lang.Normalize(t)
			}
			err = w.Pin().AddAspect(aspect, traits)
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *GrammarDecl) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *GrammarDecl) Weave(cat *weave.Catalog) (err error) {
	switch el := op.Grammar.(type) {
	// fix: why have a generic "grammar" decl, just to switch on two sub decls
	// they should be top level.
	case *grammar.Alias:
		err = cat.Schedule(weave.RequireAll, func(w *weave.Weaver) (err error) {
			name := lang.Normalize(el.AsNoun)
			if n, e := w.GetClosestNoun(name); e != nil {
				err = e
			} else {
				pen := w.Pin()
				for _, a := range el.Names {
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

	case *grammar.Directive:
		// jump/skip/hop	{"Directive:scans:":[["jump","skip","hop"],[{"As:":"jumping"}]]}
		name := strings.Join(el.Lede, "/")
		err = cat.Schedule(weave.RequireRules, func(w *weave.Weaver) error {
			return w.Pin().AddGrammar(name, el)
		})
	default:
		err = errutil.Fmt("unknown grammar %T", el)
	}
	return
}

// Execute - called by the macro runtime during weave.
func (op *DefineNounTraits) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefineNounTraits) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequirePlurals, func(w *weave.Weaver) (err error) {
		if nouns, e := safe.GetTextList(w, op.Nouns); e != nil {
			err = e
		} else if kind, e := safe.GetOptionalText(w, op.Kind, ""); e != nil {
			err = e
		} else if traits, e := safe.GetTextList(w, op.Traits); e != nil {
			err = e
		} else {
			pen := w.Pin()
			names := nouns.Strings()
			if kind, e := grok.StripArticle(kind.String()); e != nil {
				err = e
			} else if len(kind) > 0 {
				for i, name := range names {
					if name, e := grok.StripArticle(name); e != nil {
						err = errutil.Append(err, e)
					} else {
						n := lang.Normalize(name)
						if e := pen.AddNoun(n, name, kind); e != nil {
							err = errutil.Append(err, e)
						} else {
							names[i] = n // replace for the traits loop
						}
					}
				}
			}
			if traits := traits.Strings(); len(traits) > 0 && err == nil {
				for _, t := range traits {
					t := lang.Normalize(t)
					for _, n := range names {
						if e := pen.AddFieldValue(n, t, truly()); e != nil {
							err = errutil.Append(err, e)
							break // out of the traits to the next noun
						}
					}
				}
			}
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *DefinePhrase) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefinePhrase) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequireAncestry, func(w *weave.Weaver) (err error) {
		if phrase, e := safe.GetText(w, op.Phrase); e != nil {
			err = e
		} else if macro, e := safe.GetText(w, op.Macro); e != nil {
			err = e
		} else if rev, e := safe.GetOptionalBool(w, op.Reversed, false); e != nil {
			err = e
		} else if macro := lang.Normalize(macro.String()); len(macro) == 0 {
			err = errutil.New("missing macro name")
		} else {
			err = w.Pin().AddPhrase(macro, phrase.String(), rev.Bool())
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *DefineNouns) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefineNouns) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequirePlurals, func(w *weave.Weaver) (err error) {
		if nouns, e := safe.GetTextList(w, op.Nouns); e != nil {
			err = e
		} else if kind, e := safe.GetText(w, op.Kind); e != nil {
			err = e
		} else {
			names := nouns.Strings()
			if kind := kind.String(); len(kind) > 0 {
				if kind, e := grok.StripArticle(kind); e != nil {
					err = e
				} else {
					pen := w.Pin()
					for _, noun := range names {
						if noun, e := grok.StripArticle(noun); e != nil {
							err = errutil.Append(err, e)
						} else if e := pen.AddNoun(lang.Normalize(noun), noun, kind); e != nil {
							err = errutil.Append(err, e)
						}
					}
				}
			}
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *NounAssignment) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// ex. The description of the nets is xxx
func (op *NounAssignment) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequirePlurals, func(w *weave.Weaver) (err error) {
		if nouns, e := safe.GetTextList(w, op.Nouns); e != nil {
			err = e
		} else if field, e := safe.GetText(w, op.FieldName); e != nil {
			err = e
		} else if lines, e := ConvertText(op.Lines.String()); e != nil {
			err = e
		} else {
			pen := w.Pin()
			subjects := nouns.Strings()
			field, lines := field.String(), text(lines, "")
			for _, noun := range subjects {
				if noun, e := grok.StripArticle(noun); e != nil {
					err = errutil.Append(err, e)
				} else {
					n := lang.Normalize(noun)
					if e := pen.AddFieldValue(n, field, lines); e != nil {
						err = errutil.Append(err, e)
					}
				}
			}
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *DefineRelatives) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefineRelatives) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequirePlurals, func(w *weave.Weaver) (err error) {
		if rel, e := safe.GetText(w, op.Relation); e != nil {
			err = e
		} else if nouns, e := safe.GetTextList(w, op.Nouns); e != nil {
			err = e
		} else if otherNouns, e := safe.GetTextList(w, op.OtherNouns); e != nil {
			err = e
		} else {
			err = defineRelatives(w, rel.String(), nouns.Strings(), otherNouns.Strings())
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *DefineOtherRelatives) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefineOtherRelatives) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequirePlurals, func(w *weave.Weaver) (err error) {
		if rel, e := safe.GetText(w, op.Relation); e != nil {
			err = e
		} else if nouns, e := safe.GetTextList(w, op.Nouns); e != nil {
			err = e
		} else if otherNouns, e := safe.GetTextList(w, op.OtherNouns); e != nil {
			err = e
		} else {
			err = defineRelatives(w, rel.String(), otherNouns.Strings(), nouns.Strings())
		}
		return
	})
}

func defineRelatives(w *weave.Weaver, rel string, nouns, otherNouns []string) (err error) {
	pen, rel := w.Pin(), lang.Normalize(rel)
	for _, one := range nouns {
		if a, e := w.GetClosestNoun(lang.Normalize(one)); e != nil {
			err = errutil.Append(err, e)
		} else {
			for _, other := range otherNouns {
				if b, e := w.GetClosestNoun(lang.Normalize(other)); e != nil {
					err = errutil.Append(err, e)
				} else {
					if e := pen.AddPair(rel, a, b); e != nil {
						err = errutil.Append(err, e)
					}
				}
			}
		}
	}
	return
}
