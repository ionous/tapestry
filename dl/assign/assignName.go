package assign

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"github.com/ionous/errutil"
)

// note: currently templates can't analyze the context they're being called in, so:
// if the command is being used to get an object value
// and no variable can be found in the current context of the requested name,
// see if the requested name is an object instead.
func ResolveName(run rt.Runtime, name string, path DottedPath) (ret RootValue, err error) {
	// uppercase names are assumed to be requests for object names.
	if lang.IsCapitalized(name) {
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
				ret = RootValue{
					RefValue: RefValue{
						Object: meta.Variables,
						Field:  name,
						Path:   path,
					},
					RootValue: value,
				}
			}

		case g.Unknown:
			// no such variable? try as an object:
			if v, e := tryAsObject(run, name, path); g.IsUnknown(e) {
				err = g.UnknownName(name)
			} else if err != nil {
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

func tryAsObject(run rt.Runtime, name string, path DottedPath) (ret RootValue, err error) {
	if id, e := run.GetField(meta.ObjectId, name); e != nil {
		err = e
	} else if len(path) == 0 {
		ret = RootValue{
			RefValue: RefValue{
				Object: meta.ObjectId,
				Field:  name,
				Path:   path,
			},
			RootValue: id,
		}
	} else {
		// manually unpack the first dot:
		// it has to be a field name ( not an index )
		dot, remainingPath := path[0], path[1:]
		if field, ok := dot.(DotField); !ok {
			errutil.Fmt("fields should be access by name (failed trying %q.%T)", name, dot)
		} else if value, e := run.GetField(id.String(), field.Field()); e != nil {
			err = e
		} else {
			ret = RootValue{
				RefValue: RefValue{
					Object: id.String(),
					Field:  field.Field(),
					Path:   remainingPath,
				},
				RootValue: value,
			}
		}
	}
	return
}
