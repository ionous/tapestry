package assign

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign/dot"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

// note: currently templates can't analyze the context they're being called in, so:
// if the command is being used to get an object value
// and no variable can be found in the current context of the requested name,
// see if the requested name is an object instead.
func ResolveName(run rt.Runtime, name string, path dot.Path) (ret dot.Endpoint, err error) {
	// uppercase names are assumed to be requests for object names.
	if inflect.IsCapitalized(name) {
		ret, err = tryAsObject(run, name, path)

	} else {
		// otherwise, try as a variable first:
		switch value, e := run.GetField(meta.Variables, name); e.(type) {
		case nil:
			// variables and objects are both specified with text evals
			// because the template cant tell, we try to use that text as an object:
			// {Say: RenderTemplate: {.obj.indefinite_article} {.name}}}
			if aff := value.Affinity(); aff == affine.Text && len(path) > 0 {
				ret, err = tryAsObject(run, value.String(), path)
			} else {
				ret, err = resolveVariable(run, name, path)
			}

		case g.Unknown:
			// no such variable? try as an object:
			if v, e := tryAsObject(run, name, path); g.IsUnknown(e) {
				err = g.UnknownName(name)
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

func tryAsObject(run rt.Runtime, name string, path dot.Path) (ret dot.Endpoint, err error) {
	if len(path) == 0 {
		ret, err = dot.FindEndpoint(run, meta.ObjectId, dot.Path{dot.Field(name)})
	} else {
		if id, e := run.GetField(meta.ObjectId, name); e != nil {
			err = e
		} else {
			ret, err = dot.FindEndpoint(run, id.String(), path)
		}
	}
	return
}
