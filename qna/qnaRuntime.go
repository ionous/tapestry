package qna

import (
	"database/sql"
	"log"

	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/qna/pdb"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/meta"
	"git.sr.ht/~ionous/iffy/rt/print"
	"git.sr.ht/~ionous/iffy/rt/scope"
	"git.sr.ht/~ionous/iffy/rt/writer"
	"github.com/ionous/errutil"
)

func NewRuntime(db *sql.DB, signatures signatures) *Runner {
	var run *Runner
	if qdb, e := pdb.NewQueries(db); e != nil {
		panic(e) //fix: report
	} else {
		run = &Runner{
			db:         db,
			qdb:        qdb,
			values:     make(cache),
			nounValues: make(cache),
			counters:   make(counters),
			signatures: signatures,
			qnaOptions: qnaOptions{
				meta.PrintResponseNames.String(): g.BoolOf(false),
			},
		}
		run.SetWriter(print.NewAutoWriter(writer.NewStdout()))
	}
	return run
}

type Runner struct {
	db         *sql.DB
	qdb        *pdb.Query
	values     cache
	nounValues cache
	counters
	signatures
	qnaOptions
	//
	scope.Stack
	Randomizer
	writer.Sink
	currentPatterns
}

type signatures []map[uint64]interface{}

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
		ret = lang.Singularize(singular)
	}
	return
}

// func (run *Runner) IsNounInScope(id string) (ret bool) {
// 	if run.values.cache(func() (ret interface{}, err error) {
// 		ret, err = run.qdb.NounActive(id)
// 		return
// 	}, "NounActive", id); e != nil {
// 		run.Report(e)
// 	} else {
// 		ret = cached.(bool)
// 	}
// 	return
// }

func (run *Runner) PatternLabels(pat string) (ret pdb.PatternLabels, err error) {
	if c, e := run.values.cache(func() (ret interface{}, err error) {
		ret, err = run.qdb.PatternLabels(pat)
		return
	}, "PatternLabels", pat); e != nil {
		run.Report(e)
	} else {
		ret = c.(pdb.PatternLabels)
	}
	return
}

func (run *Runner) Report(e error) error {
	// fix? now what?
	log.Println(e)
	return e
}

// doesnt reformat the names.
func (run *Runner) RelateTo(a, b, relation string) error {
	return run.qdb.Relate(relation, a, b)
}

// assumes a is a valid noun
func (run *Runner) RelativesOf(rel, a string) (ret []string, err error) {
	// doesnt cache because relateTo would have to clear the cache.
	if vs, e := run.qdb.RelativesOf(rel, a); e != nil {
		err = e
	} else {
		ret = vs
	}
	return
}

// assumes b is a valid noun
func (run *Runner) ReciprocalsOf(rel, b string) (ret []string, err error) {
	if vs, e := run.qdb.ReciprocalsOf(rel, b); e != nil {
		err = e
	} else {
		ret = vs
	}
	return
}

func (run *Runner) SetField(target, rawField string, val g.Value) (err error) {
	// fix: pre-transform field
	if field := lang.Underscore(rawField); len(field) == 0 {
		err = errutil.Fmt("invalid targeted field '%s.%s'", target, rawField)
	} else {
		switch target {
		case meta.Variables:
			err = run.Stack.SetFieldByName(field, val)

		case meta.Option:
			err = run.setOption(field, val)

		case meta.Counter:
			err = run.setCounter(field, val)

		default:
			// maybe they meant to get the object?
			err = errutil.Fmt("invalid targeted field '%s.%s'", target, field)
		}
	}
	return
}

func (run *Runner) GetField(target, rawField string) (ret g.Value, err error) {
	if field := lang.Underscore(rawField); len(field) == 0 {
		err = errutil.Fmt("invalid targeted field '%s.%s'", target, rawField)
	} else {
		switch target {
		default:
			err = errutil.Fmt("GetField: unknown target %q (with field %q)", target, rawField)

		case meta.Counter:
			ret, err = run.getCounter(field)

		case meta.Domain:
			if b, e := run.qdb.IsDomainActive(field); e != nil {
				err = run.Report(e)
			} else {
				ret = g.BoolOf(b)
			}

		case meta.ObjectId:
			if ok, e := run.getObjectKind(field); e != nil {
				err = e
			} else {
				ret = g.StringOf(ok.Name)
			}

		// type of a a game object
		case meta.ObjectKind:
			if ok, e := run.getObjectKind(field); e != nil {
				err = e
			} else {
				ret = g.StringOf(ok.Kind)
			}

		case meta.ObjectKinds:
			if ok, e := run.getObjectKind(field); e != nil {
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

		// fix: see notes in qnaObject --
		case meta.ObjectValue:
			ret, err = run.getObjectByName(field)

		// all objects of the named kind
		case meta.ObjectsOfKind:
			if ns, e := run.qdb.NounsByKind(field); e != nil {
				err = run.Report(e)
			} else {
				ret = g.StringsOf(ns)
			}

		// custom options
		case meta.Option:
			if t, e := run.option(field); e != nil {
				err = run.Report(e)
			} else {
				ret = t
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
