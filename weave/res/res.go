package res

import (
	"github.com/ionous/errutil"
)

// hrm
const Missing = errutil.Error("Missing")

// internal error code
const resolved = errutil.Error("Resolved")

type Result interface {
	// returns nil when the resource is ready;
	// Missing if not, or a critical error of some sort.
	Resolve() (any, error)
}

// Resolved generates a resource that immediately returns the passed value
func Resolved(res any) Result {
	return &Cache{
		err: resolved,
		res: res,
	}
}

// Error generates a resource that always errors
// panics if the passed error is nil
func Errored(err error) Result {
	if err == nil {
		panic("error shouldnt be nil")
	}
	return &Cache{
		err: err,
	}
}

// Resolve converts a function to a resource.
func Resolve(req func() (any, error)) Result {
	return &Cache{req: req}
}

// Result unwraps the result from a successfully completed resource.
// the caller needs to know the resource is available and valid;
// panics if the error isnt nil.
func GetResult(r Result) any {
	res, e := r.Resolve()
	if e != nil {
		panic(e)
	}
	return res
}

// Await generates an array from a number of other resources.
func Await(all ...Result) (ret Result) {
	out := make([]any, 0, len(all))
	return Resolve(func() (ret any, err error) {
		for len(out) < len(all) {
			next := all[len(out)]
			if v, e := next.Resolve(); e != nil {
				break
			} else {
				out = append(out, v)
			}
		}
		if err == nil {
			ret = out
		}
		return
	})
}
