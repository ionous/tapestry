package qna

import (
	"database/sql"
	"io"
	"log"
	"os"

	"git.sr.ht/~ionous/tapestry/jsn/cin"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/print"
	"git.sr.ht/~ionous/tapestry/rt/scope"
	"git.sr.ht/~ionous/tapestry/rt/writer"
	"git.sr.ht/~ionous/tapestry/web/markup"
	"github.com/ionous/errutil"
)

func NewRuntime(db *sql.DB, signatures cin.Signatures) *Runner {
	opt := NewOptions()
	qdb, e := qdb.NewQueries(db, true)
	if e != nil {
		panic(e)
	}
	w := print.NewLineSentences(markup.ToText(os.Stdout))
	return NewRuntimeOptions(w, qdb, opt, signatures)
}

func NewRuntimeOptions(w io.Writer, qdb *qdb.Query, options Options, signatures cin.Signatures) *Runner {
	run := &Runner{
		qdb:        qdb,
		values:     make(cache),
		nounValues: make(cache),
		counters:   make(counters),
		signatures: signatures,
		options:    options,
	}
	run.SetWriter(w)
	return run
}

type Runner struct {
	qdb        *qdb.Query
	values     cache
	nounValues cache
	counters
	signatures cin.Signatures
	options    Options
	//
	scope.Stack
	Randomizer
	writer.Sink
	currentPatterns
}

func (run *Runner) NounIsNamed(noun, name string) (bool, error) {
	return run.qdb.NounIsNamed(noun, name)
}

func (run *Runner) ActivateDomain(domain string) (ret string, err error) {
	if prev, e := run.qdb.ActivateDomain(domain); e != nil {
		run.Report(e)
		err = e
	} else {
		run.values = make(cache) // fix? focus cache clear to just the domains that became inactive?
		ret = prev
	}
	return
}

func (run *Runner) SingularOf(plural string) (ret string) {
	// fix: do we want to cache? we do for everything else.
	if n, e := run.qdb.PluralToSingular(plural); e != nil {
		run.Report(e)
	} else if len(n) > 0 {
		ret = n
	} else {
		ret = lang.Singularize(plural)
	}
	return
}

func (run *Runner) PluralOf(singular string) (ret string) {
	// fix: see singularOf.
	if n, e := run.qdb.PluralFromSingular(singular); e != nil {
		run.Report(e)
	} else if len(n) > 0 {
		ret = n
	} else {
		ret = lang.Pluralize(singular)
	}
	return
}

// the last value is the results; blank if need be
func (run *Runner) PatternLabels(pat string) (ret []string, err error) {
	if c, e := run.values.cache(func() (ret interface{}, err error) {
		ret, err = run.qdb.PatternLabels(pat)
		return
	}, "PatternLabels", pat); e != nil {
		run.Report(e)
	} else {
		ret = c.([]string)
	}
	return
}

func (run *Runner) Report(e error) error {
	// fix? now what?
	log.Println(e)
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
			err = run.qdb.Relate(k.Name(), na.Id, nb.Id)
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
	} else if vs, e := run.qdb.RelativesOf(k.Name(), n.Id); e != nil {
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
	} else if vs, e := run.qdb.ReciprocalsOf(k.Name(), n.Id); e != nil {
		err = e
	} else {
		fa := k.Field(0)
		ret = g.StringsFrom(vs, fa.Type)
	}
	return
}

func (run *Runner) SetField(target, rawField string, val g.Value) (err error) {
	// fix: pre-transform field
	if field := lang.Underscore(rawField); len(field) == 0 {
		err = errutil.Fmt("invalid targeted field '%s.%s'", target, rawField)
	} else if target[0] != meta.Prefix {
		// an object from the author's story
		if obj, e := run.getObjectInfo(target); e != nil {
			err = e
		} else {
			// copies internally once it knows the type of field
			err = run.setObjectField(obj, field, val)
		}
	} else {
		// one of the predefined faux objects:
		switch target {
		case meta.Variables:
			val := g.CopyValue(val)
			err = run.Stack.SetFieldByName(field, val)

		case meta.Option:
			val := g.CopyValue(val)
			err = run.options.SetOptionByName(field, val)

		case meta.Counter:
			// doesnt copy because it assumes the value is a number anyway.
			err = run.setCounter(field, val)

		default:
			err = errutil.Fmt("invalid targeted field '%s.%s'", target, field)
		}
	}
	return
}

func (run *Runner) GetField(target, rawField string) (ret g.Value, err error) {
	// fix: pre-transform field
	if field := lang.Underscore(rawField); len(field) == 0 {
		err = errutil.Fmt("requested empty field from %q", target)
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

		case meta.Counter:
			ret, err = run.getCounter(field)

		case meta.Domain:
			if b, e := run.qdb.IsDomainActive(field); e != nil {
				err = run.Report(e)
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
				err = run.Report(e)
			} else {
				ret = g.StringsOf(k.Path())
			}

		// given a noun, return the name declared by the author
		case meta.ObjectName:
			if n, e := run.getObjectName(field); e != nil {
				err = run.Report(e)
			} else {
				ret = g.StringOf(n)
			}

		// all objects of the named kind
		case meta.ObjectsOfKind:
			if ns, e := run.qdb.NounsByKind(field); e != nil {
				err = run.Report(e)
			} else {
				ret = g.StringsOf(ns)
			}

		// custom options
		case meta.Option:
			// note: uses raw field so that it matches the meta.Options go generated stringer strings.
			if t, e := run.options.Option(rawField); e != nil {
				err = run.Report(e)
			} else {
				ret = t
			}

		case meta.PatternLabels:
			if vs, e := run.PatternLabels(field); e != nil {
				err = run.Report(e)
			} else {
				ret = g.StringsOf(vs)
			}

		case meta.PatternRunning:
			b := run.currentPatterns.runningPattern(field)
			ret = g.IntOf(b)

		case meta.Variables:
			ret, err = run.Stack.FieldByName(field)
		}
	}
	return
}
