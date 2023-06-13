package weave

import (
	"database/sql"
	"errors"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

// StoryStatement - import a single story statement.
// used during weave, and expects that the runtime is the importer's own runtime.
// ( as opposed to the story's playtime. )
func StoryStatement(run rt.Runtime, op Schedule) (err error) {
	if k, ok := run.(*Weaver); !ok {
		err = errutil.Fmt("runtime %T doesn't support story statements", run)
	} else {
		err = op.Schedule(k.d.catalog)
	}
	return
}

func (cat *Catalog) AssertAlias(opShortName string, opAliases ...string) error {
	return cat.Schedule(assert.RequireNouns, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		if shortName, ok := UniformString(opShortName); !ok {
			err = errutil.New("invalid name", opShortName)
		} else if n, e := d.GetClosestNoun(shortName); e != nil {
			err = e
		} else {
			for _, a := range opAliases {
				if a, ok := UniformString(a); !ok {
					err = errutil.Append(err, InvalidString(a))
				} else {
					err = cat.writer.Name(d.name, n.name, a, -1, at)
				}
			}
		}
		return
	})
}

// generates traits and adds them to a custom aspect kind.
func (cat *Catalog) AssertAspectTraits(opAspects string, opTraits []string) error {
	// uses the ancestry phase because it generates kinds ( one per aspect. )
	return cat.Schedule(assert.RequireDeterminers, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		// we dont singularize aspects even thought its a kind;
		// most are really singularizable anyway, and some common things like "darkness" dont singularize correctly.
		if aspect, ok := UniformString(opAspects); !ok {
			err = InvalidString(opAspects)
		} else if traits, e := UniformStrings(opTraits); e != nil {
			err = e
		} else {
			aspect := d.singularize(aspect)
			if e := d.addKind(aspect, kindsOf.Aspect.String(), at); e != nil {
				err = e
			} else if len(traits) > 0 {
				err = d.schedule(at, assert.RequireResults, func(ctx *Weaver) error {
					return cat.writer.Aspect(d.name, aspect, at, traits)
				})
			}
		}
		return
	})
}

//
func (d *Domain) singularize(name string) (ret string) {
	if d.currPhase <= assert.RequireDependencies {
		panic("singularizing before plurals are known")
	}
	if len(name) < 2 {
		ret = name //
	} else if e := d.catalog.db.QueryRow(`
	select one
	from mdl_plural
	join domain_tree
		on (uses = domain)
	where base = ?1 and many = ?2
	limit 1`, d.name, name).Scan(&ret); e != nil {
		if e != sql.ErrNoRows && d.catalog.warn != nil {
			err := errutil.Fmt("%v while singularizing %q", e, name)
			d.catalog.warn(err)
		}
		ret = lang.Singularize(name)
	}
	return
}

func (d *Domain) addKind(name, parent, at string) (err error) {
	if e := d.catalog.writer.Kind(d.name, name, parent, at); e != nil {
		if !errors.Is(e, mdl.Duplicate) {
			err = e
		} else if d.catalog.warn != nil {
			d.catalog.warn(e)
		}
	}
	return
}

func (cat *Catalog) AssertCheck(opName string, prog []rt.Execute, expect literal.LiteralValue) error {
	// uses domain phase, because it needs to ensure a domain exists
	return cat.Schedule(assert.RequireAll, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		if name, ok := UniformString(opName); !ok {
			err = InvalidString(opName)
		} else {
			err = cat.writer.Check(d.name, name, expect, prog, at)
		}
		return
	})
}

func (cat *Catalog) AssertDefinition(path ...string) error {
	return cat.Schedule(assert.RequireAll, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		if end := len(path) - 1; end <= 0 {
			err = errutil.New("path too short", path)
		} else {
			key, value := strings.Join(path[:end], "/"), path[end]
			err = cat.writer.Fact(d.name, key, value, at)
		}
		return
	})
}
