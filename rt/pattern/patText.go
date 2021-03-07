package pattern

import (
	"strings"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

// filterText - if not of text affinity, return the value
// otherwise, if its an object / id make sure it fits
// FIX: id like to remove this -- the first step is getting rid of object affinities i think.
func filterText(run rt.Runtime, f g.Field, v g.Value) (ret g.Value, err error) {
	ret = v // provisionally
	switch f.Affinity {

	case affine.Text:
		objectPrefix := "object=" // set by ObjectAsName
		if strings.HasPrefix(f.Type, objectPrefix) {
			kind := f.Type[len(objectPrefix):] // chop off the prefix
			//
			switch aff := v.Affinity(); aff {
			// fix: templates ( other things? ) are still giving us object arguments
			// tbd: maybe cant be fixed until affine.Object is completely removed.
			case affine.Object:
				ret, err = convertObject(run, v, kind)

			case affine.Text:
				ret, err = convertName(run, v.String(), kind)
			}
		}
	}
	return
}

// converts an object value to an object id
// a nil kind is okay -- it allows any type
func convertObject(run rt.Runtime, obj g.Value, kind string) (ret g.Value, err error) {
	if !safe.Compatible(obj, kind, false) {
		err = errutil.New("object", obj, "not compatible with", kind)
	} else {
		ret = g.ObjectAsText(obj)
	}
	return
}

// converts a text value to a valid object id
func convertName(run rt.Runtime, n string, kind string) (ret g.Value, err error) {
	// look up the named object...
	if len(n) == 0 {
		ret = g.StringFrom("", "object="+kind)
	} else if obj, e := safe.ObjectFromString(run, n); e != nil {
		err = e
	} else {
		ret, err = convertObject(run, obj, kind)
	}
	return
}
