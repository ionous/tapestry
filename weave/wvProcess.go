package weave

import (
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

// walk the domains and run the commands remaining in their queues
func (cat *Catalog) AssembleCatalog() (err error) {
	var ds []*Domain
	for {
		if len(cat.processing) > 0 {
			err = errutil.New("mismatched begin/end domain")
			break
		} else if len(cat.pendingDomains) == 0 {
			break
		} else if was, e := cat.assembleNext(); e != nil {
			err = e
			break
		} else {
			ds = append(ds, was)
		}
	}
	if err == nil {
		// REMOVE:
		err = cat.WriteValues(cat.writer)
	}
	return
}

func (cat *Catalog) assembleNext() (ret *Domain, err error) {
	found := -1
	// without resolving, have we resolved our parents?
	for i := 0; i < len(cat.pendingDomains); i++ {
		next := cat.pendingDomains[i]
		if next.isReadyForProcessing() {
			found = i
			break
		}
	}
	if found < 0 {
		first := cat.pendingDomains[0]
		err = errutil.New("circular or unknown domain %q", first.name)
	} else {
		// chop this one out, then process
		at := cat.pendingDomains[found]
		cat.pendingDomains = append(cat.pendingDomains[:found], cat.pendingDomains[found+1:]...)
		if e := cat.processDomain(at); e != nil {
			err = e
		} else {
			ret = at
		}
	}
	return
}

func (cat *Catalog) processDomain(d *Domain) (err error) {
	if _, e := cat.run.ActivateDomain(d.name); e != nil {
		err = e
	} else if e := cat.findRivals(); e != nil {
		err = e
	} else {
		cat.processing.Push(d)
		for p := assert.Phase(0); p < assert.NumPhases; p++ {
			ctx := Weaver{d: d, phase: p, Runtime: cat.run}
			if e := d.runPhase(&ctx); e != nil {
				err = e
				break
			}
		}
		cat.processing.Pop()
	}
	return
}
