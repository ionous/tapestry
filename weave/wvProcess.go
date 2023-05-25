package weave

import (
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

var TempSplit = assert.AncestryPhase + 1

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
		err = cat.oldProcess(ds)
	}
	return
}

// FIX: DONT want to walk across all domains first.
// currently: walks across all domains for each phase to support things like fields:
// which exist per kind but which can be added to by multiple domains.
func (cat *Catalog) oldProcess(ds []*Domain) (err error) {
Loop:
	for p := TempSplit; p < assert.NumPhases; p++ {
		for _, d := range ds {
			cat.processing.Push(d)
			ctx := Weaver{d: d, phase: p, Runtime: cat.run}
			if e := d.runPhase(&ctx); e != nil {
				err = e
				break Loop
			} else if e := cat.postPhase(p, d); e != nil {
				err = e
				break Loop
			}
			cat.processing.Pop()
		}
	}
	if err == nil && cat.writer != nil {
		for p := assert.Phase(0); p < assert.NumPhases; p++ {
			if e := cat.writePhase(p); e != nil {
				err = e
				break
			}
		}
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
		err = errutil.New("circular reference or unknown domain in", first.name)
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
	} else if e := cat.findConflicts(); e != nil {
		err = e
	} else {
		cat.processing.Push(d)
	Loop:
		for p := assert.Phase(0); p < TempSplit; p++ {
			ctx := Weaver{d: d, phase: p, Runtime: cat.run}
			if e := d.runPhase(&ctx); e != nil {
				err = e
				break Loop
			} else {
				if p == assert.AncestryPhase {
					if ks, e := d.resolveKinds(); e != nil {
						err = e
						break Loop
					} else {
						for _, dep := range ks {
							k, ancestors := dep.Leaf().(*ScopedKind), dep.Strings(true)
							if e := cat.writer.Kind(d.name, k.name, ancestors, k.at); e != nil {
								err = e
								break Loop
							}
						}
					}
				}
			}
		}
		cat.processing.Pop()
	}
	return
}
