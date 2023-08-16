package qna

import (
	"io"
	"log"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/qna/decoder"
	"git.sr.ht/~ionous/tapestry/qna/query"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/scope"
	"git.sr.ht/~ionous/tapestry/rt/writer"
	"github.com/ionous/errutil"
)

func NewRuntimeOptions(w io.Writer, q query.Query, d decoder.Decoder, options Options) *Runner {
	run := &Runner{
		query:      q,
		decode:     d,
		values:     make(cache),
		nounValues: make(cache),
		counters:   make(counters),
		options:    options,
		scope:      scope.Chain{Scope: scope.Empty{}},
	}
	run.SetWriter(w)
	return run
}

type Runner struct {
	query  query.Query
	decode decoder.Decoder

	values     cache
	nounValues cache
	counters
	options Options
	//
	scope scope.Chain
	Randomizer
	writer.Sink
	currentPatterns
}

func (run *Runner) NounIsNamed(noun, name string) (bool, error) {
	return run.query.NounIsNamed(noun, name)
}

func (run *Runner) ActivateDomain(domain string) (ret string, err error) {
	if prev, e := run.query.ActivateDomain(domain); e != nil {
		err = run.reportError(e)
	} else {
		run.values = make(cache) // fix? focus cache clear to just the domains that became inactive?
		ret = prev
	}
	return
}

func (run *Runner) SingularOf(plural string) (ret string) {
	// fix: do we want to cache? we do for everything else. ( see: singular of )
	if len(plural) < 2 {
		ret = plural
	} else if n, e := run.query.PluralToSingular(plural); e != nil {
		run.reportError(e)
	} else if len(n) > 0 {
		ret = n
	} else {
		ret = lang.Singularize(plural)
	}
	return
}

func (run *Runner) PluralOf(singular string) (ret string) {
	// cache? see also singularOf.
	// dont bother with one letter kinds ( ex. tests )
	if len(singular) < 2 {
		ret = singular
	} else if n, e := run.query.PluralFromSingular(singular); e != nil {
		run.reportError(e)
	} else if len(n) > 0 {
		ret = n
	} else {
		ret = lang.Pluralize(singular)
	}
	return
}

func (run *Runner) OppositeOf(word string) (ret string) {
	if n, e := run.query.OppositeOf(word); e != nil {
		run.reportError(e)
	} else if len(n) > 0 {
		ret = n
	} else {
		ret = word
	}
	return
}

// the last value is the results; blank if need be
func (run *Runner) getPatternLabels(pat string) (ret []string, err error) {
	if c, e := run.values.cache(func() (ret any, err error) {
		ret, err = run.query.PatternLabels(pat)
		return
	}, "PatternLabels", pat); e != nil {
		err = run.reportError(e)
	} else {
		ret = c.([]string)
	}
	return
}

func (run *Runner) reportError(e error) error {
	log.Println(e) // fix? now what?
	return e
}

// doesnt reformat the names.
func (run *Runner) RelateTo(a, b, rel string) (err error) {
	if k, e := run.getKindOf(rel, kindsOf.Relation.String()); e != nil {
		err = e
	} else if na, e := run.getObjectInfo(a); e != nil {
		err = e
	} else if nb, e := run.getObjectInfo(b); e != nil {
		err = e
	} else {
		fa, fb := k.Field(0), k.Field(1)
		if _, e := run.getKindOf(na.Kind, fa.Type); e != nil {
			err = e
		} else if _, e := run.getKindOf(nb.Kind, fb.Type); e != nil {
			err = e
		} else {
			err = run.query.Relate(k.Name(), na.Id, nb.Id)
		}
	}
	return
}

func (run *Runner) RelativesOf(a, rel string) (ret g.Value, err error) {
	// note: we dont have to validate the type of the noun....
	// if its not valid, it wont appear in the relation.
	if n, e := run.getObjectInfo(a); e != nil {
		err = e
	} else if k, e := run.getKindOf(rel, kindsOf.Relation.String()); e != nil {
		err = e
	} else if vs, e := run.query.RelativesOf(k.Name(), n.Id); e != nil {
		err = e // doesnt cache because relateTo would have to clear the cache.
	} else {
		fb := k.Field(1)
		ret = g.StringsFrom(vs, fb.Type)
	}
	return
}

func (run *Runner) ReciprocalsOf(b, rel string) (ret g.Value, err error) {
	// note: we dont have to validate the type of the noun....
	// if its not valid, it wont appear in the relation.
	if n, e := run.getObjectInfo(b); e != nil {
		err = e
	} else if k, e := run.getKindOf(rel, kindsOf.Relation.String()); e != nil {
		err = e
	} else if vs, e := run.query.ReciprocalsOf(k.Name(), n.Id); e != nil {
		err = e
	} else {
		fa := k.Field(0)
		ret = g.StringsFrom(vs, fa.Type)
	}
	return
}

func (run *Runner) SetField(target, rawField string, val g.Value) (err error) {
	// fix: pre-transform field name
	if field := lang.Normalize(rawField); len(field) == 0 {
		err = errutil.Fmt("invalid targeted field '%s.%s'", target, rawField)
	} else if target[0] != meta.Prefix {
		// an object from the author's story
		if obj, e := run.getObjectInfo(target); e != nil {
			err = e
		} else {
			// copies val internally once it knows the type of field
			err = run.setObjectField(obj, field, val)
		}
	} else {
		// one of the predefined faux objects:
		switch target {
		case meta.Variables:
			cpy := g.CopyValue(val)
			err = run.scope.SetFieldByName(field, cpy)

		case meta.Option:
			cpy := g.CopyValue(val)
			err = run.options.SetOptionByName(field, cpy)

		case meta.Counter:
			// doesnt copy because it errors if the value isn't a number
			// ( and numbers dont need to be copied ).
			err = run.setCounter(field, val)

		case meta.ValueChanged:
			if val.Affinity() != affine.Text {
				err = errutil.New("the value of value changed should be the name of the field that changed")
			} else {
				// unpack the real target and field
				switch target, field := field, val.String(); target {
				case meta.Variables:
					err = run.scope.SetFieldDirty(field)
				default:
					// todo: example, flag object or db for save.
					// for now, simply verify that the field exists.
					_, err = run.GetField(target, field)
				}
			}

		default:
			err = errutil.Fmt("invalid targeted field '%s.%s'", target, field)
		}
	}
	return
}

func (run *Runner) GetField(target, rawField string) (ret g.Value, err error) {
	// fix: pre-transform field
	if field := lang.Normalize(rawField); len(field) == 0 {
		err = errutil.Fmt("requested an empty %q", target)
	} else if target[0] != meta.Prefix {
		// an object from the author's story
		if obj, e := run.getObjectInfo(target); e != nil {
			err = e
		} else {
			ret, err = run.getObjectField(obj, field)
		}
	} else {
		// one of the predefined faux objects:
		switch target {
		default:
			err = errutil.Fmt("GetField: unknown target %q (with field %q)", target, rawField)
			// not one of the predefined options?
		case meta.Response:
			// response arent implemented yet.
			// note: uses raw field so that it matches the meta.Options go generated stringer strings.
			if flag, e := run.options.Option(meta.PrintResponseNames.String()); e != nil {
				err = e
			} else if flag.Affinity() == affine.Bool && flag.Bool() {
				ret = g.StringOf(field)
			} else {
				err = g.UnknownResponse(rawField)
			}

		case meta.Counter:
			ret, err = run.getCounter(field)

		case meta.Domain:
			if b, e := run.query.IsDomainActive(field); e != nil {
				err = run.reportError(e)
			} else {
				ret = g.BoolOf(b)
			}

		case meta.ObjectId:
			if ok, e := run.getObjectInfo(field); e != nil {
				err = e
			} else {
				ret = g.StringFrom(ok.Id, ok.Kind)
			}

		// type of a game object
		case meta.ObjectKind:
			if ok, e := run.getObjectInfo(field); e != nil {
				err = e
			} else {
				ret = g.StringOf(ok.Kind)
			}

		case meta.ObjectKinds:
			if ok, e := run.getObjectInfo(field); e != nil {
				err = e
			} else if k, e := run.GetKindByName(ok.Kind); e != nil {
				err = run.reportError(e)
			} else {
				ret = g.StringsOf(k.Path())
			}

		// given a noun, return the name declared by the author
		case meta.ObjectName:
			if n, e := run.getObjectName(field); e != nil {
				err = run.reportError(e)
			} else {
				ret = g.StringOf(n) // tbd: should these have a type?
			}

		// all objects of the named kind
		case meta.ObjectsOfKind:
			if ns, e := run.query.NounsByKind(field); e != nil {
				err = run.reportError(e)
			} else {
				ret = g.StringsOf(ns)
			}

		// custom options
		case meta.Option:
			// note: uses raw field so that it matches the meta.Options go generated stringer strings.
			if t, e := run.options.Option(rawField); e != nil {
				err = run.reportError(e)
			} else {
				ret = t
			}

		case meta.PatternLabels:
			if vs, e := run.query.PatternLabels(field); e != nil {
				err = run.reportError(e)
			} else {
				ret = g.StringsOf(vs)
			}

		case meta.PatternRunning:
			b := run.currentPatterns.runningPattern(field)
			ret = g.IntOf(b)

		case meta.Variables:
			ret, err = run.scope.FieldByName(field)
		}
	}
	return
}
