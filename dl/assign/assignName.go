package assign

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign/dot"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

// note: currently templates can't analyze the context they're being called in, so:
// if the command is being used to get an object value
// and no variable can be found in the current context of the requested name,
// see if the requested name is an object instead.
func GetNamedValue(run rt.Runtime, name string, dots []Dot) (ret rt.Value, err error) {
	if path, e := resolveDots(run, dots); e != nil {
		err = fmt.Errorf("%w resolving path for %s", e, name)
	} else {
		ret, err = getNamedValue(run, name, path)
	}
	return
}

func getNamedValue(run rt.Runtime, name string, path dot.Path) (ret rt.Value, err error) {
	// uppercase names are assumed to be requests for object names.
	if inflect.IsCapitalized(name) {
		ret, err = getObjValue(run, name, path)

	} else {
		// otherwise, try as a variable first:
		switch value, e := run.GetField(meta.Variables, name); e.(type) {
		case nil:
			// variables and objects are both specified with text evals
			// because the template cant tell, we try to use that text as an object:
			// {Say: RenderTemplate: {.obj.indefinite_article} {.name}}}
			if aff := value.Affinity(); aff == affine.Text && len(path) > 0 {
				ret, err = getObjValue(run, value.String(), path)
			} else {
				at := dot.MakeReference(run, meta.Variables)
				if at, e := at.Dot(dot.Field(name)); e != nil {
					err = e
				} else if at, e := at.DotPath(path); e != nil {
					err = e
				} else {
					ret, err = at.GetValue()
				}
			}

		case rt.Unknown:
			// no such variable? try as an object:
			if v, e := getObjValue(run, name, path); rt.IsUnknown(e) {
				err = rt.UnknownName(name)
			} else if e != nil {
				err = e
			} else {
				ret = v
			}
		default:
			err = e
		}
	}
	return
}

func getObjValue(run rt.Runtime, name string, path dot.Path) (ret rt.Value, err error) {
	if id, e := run.GetField(meta.ObjectId, name); e != nil {
		err = e
	} else if len(path) == 0 {
		ret = id
	} else {
		at := dot.MakeReference(run, id.String())
		if at, e := at.DotPath(path); e != nil {
			err = e
		} else {
			ret, err = at.GetValue()
		}
	}
	return
}
