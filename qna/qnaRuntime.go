package qna

import (
	"database/sql"
	"log"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/qna/decoder"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/scope"
	"git.sr.ht/~ionous/tapestry/rt/writer"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"github.com/ionous/errutil"
)

// Callbacks for important system level changes
type Notifier struct {
	StartedScene    func(domains []string)
	EndedScene      func(domains []string)
	ChangedState    func(noun, aspect, oldState, newState string)
	ChangedRelative func(a, b, rel string)
	// ChangedValue()
}

func NewRuntime(db *sql.DB, d decoder.Decoder) (*Runner, error) {
	return NewRuntimeOptions(db, d, NewOptions())
}

func NewRuntimeOptions(db *sql.DB, d decoder.Decoder, opt Options) (ret *Runner, err error) {
	cacheErrors := opt.cacheErrors()
	if q, e := qdb.NewQueries(db); e != nil {
		err = e
	} else {
		if d == nil {
			d = decoder.DecodeNone("unsupported decoder")
		}
		ret = &Runner{
			db:         db,
			query:      q,
			decode:     d,
			values:     makeCache(cacheErrors),
			nounValues: makeCache(cacheErrors),
			counters:   make(counters),
			options:    opt,
			scope:      scope.Chain{Scope: scope.Empty{}},
		}
		ret.SetWriter(log.Writer())
	}
	return
}

type Runner struct {
	db         *sql.DB
	query      query.Query
	decode     decoder.Decoder
	notify     Notifier
	values     cache
	nounValues cache // read isn't kept across domains
	counters
	options Options
	//
	scope scope.Chain
	Randomizer
	writer.Sink
	currentPatterns
}

func (run *Runner) SetNotifier(n Notifier) {
	run.notify = n
}

func (run *Runner) reportError(e error) error {
	log.Println(e) // fix? now what?
	return e
}

func (run *Runner) ActivateDomain(domain string) (err error) {
	if ends, begins, e := run.query.ActivateDomains(domain); e != nil {
		err = run.reportError(e)
	} else {
		if len(ends) > 0 {
			// fix? the domain is already out of scope
			// might want to rewind one by one just after running each end event
			// ( begin has a similar issue if some subdomain sets a different value say than its parent )
			if e := run.domainChanged(ends, "ends"); e != nil {
				err = errutil.Append(err, e)
			} else if notify := run.notify.EndedScene; notify != nil {
				notify(ends)
			}
		}
		run.values.reset() // fix? focus cache clear to just the domains that became inactive?
		if len(domain) == 0 {
			run.nounValues.reset()
		}
		if len(begins) > 0 {
			if e := run.domainChanged(begins, "begins"); e != nil {
				err = errutil.Append(err, e)
			} else if notify := run.notify.StartedScene; notify != nil {
				notify(begins)
			}
		}
	}
	return
}

// shift to notifier?
func (run *Runner) domainChanged(ds []string, evt string) (err error) {
	for _, d := range ds {
		name := d + " " + evt
		if pat, e := run.getKind(name); e == nil {
			if _, e := run.call(pat, rt.NewRecord(pat), affine.None); e != nil {
				err = errutil.Append(err, run.reportError(e))
				break
			}
		}
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
		ret = inflect.Singularize(plural)
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
		ret = inflect.Pluralize(singular)
	}
	return
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
		} else if e := run.query.Relate(k.Name(), na.Id, nb.Id); e != nil {
			err = e
		} else if notify := run.notify.ChangedRelative; notify != nil {
			notify(a, b, rel)
		}
	}
	return
}

func (run *Runner) RelativesOf(a, rel string) (ret rt.Value, err error) {
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
		ret = rt.StringsFrom(vs, fb.Type)
	}
	return
}

func (run *Runner) ReciprocalsOf(b, rel string) (ret rt.Value, err error) {
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
		ret = rt.StringsFrom(vs, fa.Type)
	}
	return
}

func (run *Runner) SetField(target, rawField string, val rt.Value) (err error) {
	// fix: pre-transform field name
	if field := inflect.Normalize(rawField); len(field) == 0 {
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
			cpy := rt.CopyValue(val)
			err = run.scope.SetFieldByName(field, cpy)

		case meta.Option:
			cpy := rt.CopyValue(val)
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

func (run *Runner) GetField(target, rawField string) (ret rt.Value, err error) {
	// fix: pre-transform field
	if field := inflect.Normalize(rawField); len(field) == 0 {
		err = errutil.Fmt("GetField given an empty field for target %q", target)
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
			if b, e := run.query.IsDomainActive(field); e != nil {
				err = run.reportError(e)
			} else {
				ret = rt.BoolOf(b)
			}

		case meta.FieldsOfKind:
			if k, e := run.getKind(field); e != nil {
				err = run.reportError(e)
			} else {
				var fs []string
				for i, cnt := 0, k.NumField(); i < cnt; i++ {
					f := k.Field(i)
					fs = append(fs, f.Name)
				}
				ret = rt.StringsOf(fs)
			}

		case meta.KindAncestry:
			if ks, e := run.getAncestry(field); e != nil {
				err = run.reportError(e)
			} else {
				ret = rt.StringsOf(ks)
			}

		case meta.ObjectAliases:
			if ns, e := run.getObjectNames(field); e != nil {
				err = run.reportError(e)
			} else {
				ret = rt.StringsOf(ns)
			}

		case meta.ObjectId:
			if ok, e := run.getObjectInfo(field); e != nil {
				err = e
			} else {
				ret = rt.StringFrom(ok.Id, ok.Kind)
			}

		// type of a game object
		case meta.ObjectKind:
			if ok, e := run.getObjectInfo(field); e != nil {
				err = e
			} else {
				ret = rt.StringOf(ok.Kind)
			}

		case meta.ObjectKinds:
			if ok, e := run.getObjectInfo(field); e != nil {
				err = e
			} else if ks, e := run.getAncestry(ok.Kind); e != nil {
				err = run.reportError(e)
			} else {
				ret = rt.StringsOf(ks)
			}

		// given a noun, return the name declared by the author
		case meta.ObjectName:
			if n, e := run.getObjectName(field); e != nil {
				err = run.reportError(e)
			} else {
				ret = rt.StringOf(n) // tbd: should these have a type?
			}

		// all objects of the named kind
		case meta.ObjectsOfKind:
			if k, e := run.getKind(field); e != nil {
				err = run.reportError(e)
			} else if ns, e := run.query.NounsByKind(k.Name()); e != nil {
				err = run.reportError(e)
			} else {
				ret = rt.StringsOf(ns)
			}

		// custom options
		case meta.Option:
			// note: uses raw field so that it matches the meta.Options go generated stringer strings.
			// fix? specify those strings as their inflect.Normalized versions using -linecomment?
			if t, e := run.options.OptionByName(rawField); e != nil {
				err = run.reportError(e)
			} else {
				ret = t
			}

		case meta.PatternLabels:
			if vs, e := run.query.PatternLabels(field); e != nil {
				err = run.reportError(e)
			} else {
				ret = rt.StringsOf(vs)
			}

		case meta.PatternRunning:
			b := run.currentPatterns.runningPattern(field)
			ret = rt.IntOf(b)

		case meta.Response:
			if flag, e := run.options.Option(meta.PrintResponseNames); e != nil {
				err = e
			} else if flag.Affinity() == affine.Bool && flag.Bool() {
				// fix? hen response are printed inline, they get munged together
				// manually add some spacing.
				ret = rt.StringOf(field + ". ")
			} else if k, e := run.getKind(kindsOf.Response.String()); e != nil {
				err = errutil.New("couldnt find response table", e)
			} else if i := k.FieldIndex(field); i < 0 {
				// ^ fix: this is an in-order search; have a custom cached lookup for responses?
				err = rt.UnknownResponse(rawField)
			} else if vs, e := run.getKindValues(k); e != nil {
				err = e
			} else {
				ret = rt.Empty // provisionally
				for _, kv := range vs {
					if kv.i > i {
						break
					} else if kv.i == i {
						if val, e := kv.val.GetAssignedValue(run); e != nil {
							err = e
							break
						} else {
							ret = val
						}
					}
				}
			}
		case meta.Variables:
			ret, err = run.scope.FieldByName(field)
		}
	}
	return
}
