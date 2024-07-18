package weave

import (
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/object"
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

type Domain struct {
	name        string
	cat         *Catalog
	finished    bool           //
	initializes []rt.Execute   // all of type object.SetValue
	pos         compact.Source // location of scene begin
	startup     []rt.Execute
	proc        Processing
}

func (d *Domain) Name() string {
	return d.name

}
func (d *Domain) AddInitialValue(noun, field string, val rt.Assignment) {
	d.initializes = append(d.initializes, &object.SetValue{
		Target: &object.ObjectDot{
			NounName: literal.T(noun),
			Dot:      object.MakeDot(field),
		},
		Value: val,
	})
}

func (d *Domain) runAll() error {
	return d.proc.runAll()
}

// write initial values....
func (d *Domain) finalizeDomain(pen *mdl.Pen) (err error) {
	d.finished = true
	if len(d.initializes) > 0 {
		if e := pen.WriteSceneStart(0, d.initializes); e != nil {
			err = e
		} else {
			d.initializes = nil
		}
	}
	if len(d.startup) > 0 && err == nil {
		if e := pen.WriteSceneStart(1, d.startup); e != nil {
			err = e
		} else {
			d.startup = nil
		}
	}
	return
}
