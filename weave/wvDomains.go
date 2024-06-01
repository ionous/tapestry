package weave

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/object"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
	"github.com/ionous/errutil"
)

type Domain struct {
	name          string
	cat           *Catalog
	currPhase     weaver.Phase // updated during weave, ends at NumPhases
	steps         []StepFunction
	scheduling    [weaver.NumPhases][]memento // separates commands into phases
	initialValues initialValues               // all of type object.SetValue
}

type initialValues []rt.Execute

func (in initialValues) add(noun, field string, val rt.Assignment) initialValues {
	return append(in, &object.SetValue{
		Target: &object.ObjectDot{
			Name: &literal.TextValue{Value: noun},
			Dot:  object.MakeDot(field),
		},
		Value: val,
	})
}

type memento struct {
	cb ScheduledCallback
	at string
}

// write initial values....
func (d *Domain) finalizeDomain() (err error) {
	if len(d.initialValues) > 0 {
		domainName := d.name
		eventName := domainName + " " + "begins"
		pin := d.cat.Modeler.Pin(domainName, "")
		pb := mdl.NewPatternBuilder(eventName)
		pb.AppendRule(0, rt.Rule{
			Name: "initial value rule",
			Exe:  d.initialValues,
		})
		if e := pin.AddPattern(pb.Pattern); e != nil {
			err = e
		} else {
			d.initialValues = nil
		}
	}
	return
}

func (d *Domain) Name() string {
	return d.name
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

func (d *Domain) step(_ string, cb StepFunction) (err error) {
	if z := d.currPhase; z < 0 {
		err = errutil.Fmt("domain %q already finished", d.name)
	} else {
		d.steps = append(d.steps, cb)
	}
	return
}

func (d *Domain) schedule(at string, when weaver.Phase, what ScheduledCallback) (err error) {
	// when we are not running we are in phase zero; the first active phase is index 1
	if z := d.currPhase; z < 0 {
		err = errutil.Fmt("domain %q already finished", d.name)
	} else if z > when {
		err = errutil.Fmt("domain %q processing %s phase %s already passed",
			d.name, z, when)
	} else {
		m := memento{what, at}
		if z < when {
			d.scheduling[when] = append(d.scheduling[when], m)
		} else {
			// if its the same phase, try to run immediately;
			// this is important for jess which matches and then immediately schedules:
			// that way one sentence can match dependent on the results of the previous sentence.
			if e := d.runOne(m); e != nil && !errors.Is(e, mdl.Missing) {
				err = e
			} else if e != nil {
				d.scheduling[z] = append(d.scheduling[z], m)
			}
		}
	}
	return
}

func (d *Domain) runPhase(z weaver.Phase) (err error) {
	d.currPhase = z // hrmm
	if e := d.runSteps(z); e != nil {
		err = e
	} else {
		var phase []memento
		phase, d.scheduling[z] = d.scheduling[z], nil
		for len(phase) > 0 {
			var lastMissing error
			lastCnt := len(phase)
			if keep, e := d.runSchedule(phase, &lastMissing); e != nil {
				err = e
				break
			} else if add := len(d.scheduling[z]); lastCnt == keep && add == 0 {
				err = fmt.Errorf("%s; couldn't finish phase %s", lastMissing, z)
				break
			} else {
				phase = phase[:keep]
				if add > 0 {
					// add any newly scheduled elements to this phase
					phase = append(phase, d.scheduling[z]...)
					d.scheduling[z] = nil
				}
			}
		}
		return
	}
	return
}

// compacts the passed array, removing any mementos that finished
func (d *Domain) runSchedule(phase []memento, lastMissing *error) (ret int, err error) {
	var keep int
	for _, next := range phase {
		if e := d.runOne(next); e != nil && !errors.Is(e, mdl.Missing) {
			err = e
			break
		} else if e != nil {
			*lastMissing = e
			phase[keep] = next
			keep++
		}
	}
	if err == nil {
		ret = keep
	}
	return
}

func (d *Domain) runOne(m memento) error {
	pen := d.cat.Modeler.Pin(d.name, m.at)
	w := localWeaver{d, pen}
	run := d.cat.GetRuntime()
	return m.cb(w, run)
}

// primarily exists to simplify jess.... hrmmm...
// runs the same set of callback every phrase
// rather than having to generate a new function for every scheduled process.
func (d *Domain) runSteps(z weaver.Phase) (err error) {
	var keep int
	steps := d.steps
	d.steps = nil // watch out for any steps scheduled while we are running
	for _, cb := range steps {
		if ok, e := cb(z); e != nil {
			err = e
		} else {
			if !ok {
				if z+1 == weaver.NumPhases {
					err = errors.New("processing incomplete")
				} else {
					steps[keep] = cb
					keep++
				}
			}
		}
	}
	if err == nil {
		// these wont garbage collect without copying
		// that's probably fine.
		d.steps = append(steps[:keep], d.steps...)
	}
	return
}
