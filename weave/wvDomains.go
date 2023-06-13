package weave

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

type Domain struct {
	name       string
	catalog    *Catalog
	currPhase  assert.Phase                     // updated during weave, ends at NumPhases
	scheduling [assert.RequireAll + 1][]memento // separates commands into phases
}

type memento struct {
	cb func(*Weaver) error
	at string
}

func (op *memento) call(ctx *Weaver) error {
	ctx.at = op.at
	return op.cb(ctx)
}

func (cat *Catalog) Schedule(when assert.Phase, what func(*Weaver) error) (err error) {
	if d, ok := cat.processing.Top(); !ok {
		err = errutil.New("unknown top level domain")
	} else {
		err = d.schedule(cat.cursor, when, what)
	}
	return
}

// have all parent domains been processed?
func (d *Domain) isReadyForProcessing() bool {
	return nil == d.visit(func(uses *Domain) (err error) {
		if d != uses && uses.currPhase <= assert.RequireAll {
			err = errutil.New("break")
		}
		return
	})
}

func (d *Domain) schedule(at string, when assert.Phase, what func(*Weaver) error) (err error) {
	if d.currPhase > when {
		err = errutil.Fmt("scheduling error for %q: currently in %s asking, for %s.",
			d.name, d.currPhase, when)
	} else /*if when == d.currPhase {
		err = what(&ctx)
	} else */{
		d.scheduling[when] = append(d.scheduling[when], memento{what, at})
	}
	return
}

// return the domain hierarchy: the ancestors ending just before the domain itself.
// direct parents may not be contiguous ( depending on whether their ancestors overlap. )
func (d *Domain) Resolve() (ret []string, err error) {
	c := d.catalog // we shouldnt have to worry about dupes, because in theory we didnt add them.
	if rows, e := c.db.Query(`select uses from domain_tree 
		where base = ?1 order by dist desc`, d.name); e != nil {
		err = e
	} else {
		ret, err = tables.ScanStrings(rows)
	}
	return
}

// fix? used by "isReadyForProcessing" -- is there a better way.
func (d *Domain) visit(visit func(d *Domain) error) (err error) {
	cat := d.catalog
	if tree, e := d.Resolve(); e != nil {
		err = e
	} else {
		for _, el := range tree {
			if p, ok := cat.GetDomain(el); !ok {
				err = errutil.Fmt("unexpected domain %q", el)
				break
			} else if e := visit(p); e != nil {
				err = e
				break
			}
		}
	}
	return
}

func (d *Domain) runPhase(ctx *Weaver) (err error) {
	w := ctx.phase
	d.currPhase = w // hrmm
	redo := struct {
		cnt int
		err error
	}{}
	// don't range over the slice since the contents can change during traversal.
	// tbd: have "Schedule" immediately execute the statement if in the correct phase?
	els := &d.scheduling[w]
Loop:
	for len(*els) > 0 {
		// slice the next element out of the list
		next := (*els)[0]
		(*els) = (*els)[1:]

		switch e := next.call(ctx); {
		case e == nil:
			redo.cnt, redo.err = 0, nil
		case errors.Is(e, mdl.Missing):
			redo.err = errutil.Append(redo.err, e)
			if redo.cnt < len((*els)) {
				// add redo elements back into the list
				(*els) = append((*els), next)
				redo.cnt++
			} else {
				if d.catalog.warn != nil {
					e := errutil.New(w, "didn't finish")
					d.catalog.warn(e)
				}
				err = errutil.Append(err, redo.err)
				break Loop
			}
		case errors.Is(e, mdl.Duplicate):
			if d.catalog.warn != nil {
				d.catalog.warn(e)
			}
		default:
			err = errutil.Append(err, e)
		}
	}
	d.currPhase++
	return
}
