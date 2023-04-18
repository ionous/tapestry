package story

import (
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *DefineScene) PreImport(k *imp.Importer) (ret interface{}, err error) {
	if name, e := safe.GetOptionalText(k, op.Scene, ""); e != nil {
		err = e
	} else if dependsOn, e := safe.GetOptionalTexts(k, op.DependsOn, nil); e != nil {
		err = e
	} else if e := k.BeginDomain(name.String(), dependsOn.Strings()); e != nil {
		err = e
	} else {
		ret = op
	}
	return
}

func (op *DefineScene) PostImport(k *imp.Importer) (err error) {
	return k.EndDomain() // note: could cache the name in op.Markup at a special key if need be.
}
