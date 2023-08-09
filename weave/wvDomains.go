package weave

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/ionous/errutil"
)

type Domain struct {
	name       string
	cat        *Catalog
	currPhase  Phase                     // updated during weave, ends at NumPhases
	scheduling [RequireAll + 1][]memento // separates commands into phases
	suspended  []memento                 // for missing definitions
}

type memento struct {
	cb    ScheduledCallback
	at    string
	phase Phase
	err   error
}

func (d *Domain) Name() string {
	return d.name
}

func (op *memento) call(w *Weaver) error {
	w.At = op.at
	return op.cb(w)
}

// have all parent domains been processed?
func (d *Domain) isReadyForProcessing() (okay bool, err error) {
	cat := d.cat
	// get the domain hierarchy: the ancestors ending just before the domain itself.
	// direct parents may not be contiguous ( depending on whether their ancestors overlap. )
	if rows, e := cat.db.Query(`select uses from domain_tree 
		where base = ?1 order by dist desc`, d.name); e != nil {
		err = e
	} else if tree, e := tables.ScanStrings(rows); e != nil {
		err = e
	} else {
		okay = true // provisionally
		for _, name := range tree {
			if uses, ok := cat.domains[name]; !ok {
				okay = false
				break
			} else if (d != uses) && (uses.currPhase != -1) {
				okay = false
				break
			}
		}
	}
	return
}

func (d *Domain) schedule(at string, when Phase, what ScheduledCallback) (err error) {
	if d.currPhase < 0 {
		err = errutil.Fmt("domain %q already finished", d.name)
	} else if d.currPhase <= when {
		d.scheduling[when] = append(d.scheduling[when], memento{what, at, when, nil})
	} else {
		w := Weaver{Catalog: d.cat, Domain: d.name, Phase: d.currPhase, Runtime: d.cat.run}
		if e := what(&w); errors.Is(e, mdl.Missing) {
			d.suspended = append(d.suspended, memento{what, at, when, e})
		} else {
			err = e
		}
	}
	return
}

func (d *Domain) runPhase(w *Weaver) (err error) {
	phase := w.Phase
	d.currPhase = phase // hrmm
	// don't range over the slice since the contents can change during traversal.
	// tbd; may no longer be true.
	els := &d.scheduling[phase]

	for len(*els) > 0 {
		// slice the next element out of the list
		next := (*els)[0]
		(*els) = (*els)[1:]

		switch e := next.call(w); {
		case errors.Is(e, mdl.Missing):
			next.err = e
			d.suspended = append(d.suspended, next)

		case e != nil:
			err = errutil.Append(err, e)
		}
	}
	d.currPhase++
	return
}

// tbd: the suspended mnd flush model is a little low bar.
// it'd be better to categorize individual statements -- even from macros
// and trigger them as their needs are satisfied; how exactly? not sure.
// the big mish mash are patterns,params,locals, and returns --
// which need to be written in a specific order.
func (d *Domain) flush(ignore bool) (err error) {
	w := Weaver{Catalog: d.cat, Domain: d.name, Runtime: d.cat.run}
	redo := struct {
		cnt int
		err error
	}{}

Loop:
	for len(d.suspended) > 0 {
		// slice the next element out of the list
		next := d.suspended[0]
		d.suspended = d.suspended[1:]
		w.Phase = next.phase

		switch e := next.call(&w); {
		case e == nil:
			// every success, abandon all old errors and try everything over again.
			redo.cnt, redo.err = 0, nil

		case errors.Is(e, mdl.Missing):
			// append to rack all that are missing
			redo.err = errutil.Append(redo.err, e)
			// add redo elements back into the list
			next.err = e
			d.suspended = append(d.suspended, next)
			// might still have statements to try?
			// keep going
			if redo.cnt = redo.cnt + 1; redo.cnt > len(d.suspended) {
				// if we have visited every suspended element
				// an haven't progressed; we're done.
				// return all the errors.
				if !ignore {
					err = redo.err
				}
				break Loop
			}

		default:
			// accumulate all errors
			err = errutil.Append(err, e)
		}
	}
	return
}
