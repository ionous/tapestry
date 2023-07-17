package res

import "errors"

// Cache helper to cache results from Result.
type Cache struct {
	req func() (any, error)
	err error
	res any
}

func (r *Cache) Resolve() (ret any, err error) {
	// when nil, still need to check:
	if r.err == nil || errors.Is(r.err, Missing) {
		// record last error, might be missing or critical.
		if v, e := r.req(); e != nil {
			r.err = e
			err = e
		} else {
			// resolved okay:
			r.res, r.err = v, resolved
			ret = v
		}
	} else {
		if r.err == resolved {
			ret = r.res
		} else {
			err = r.err
		}
	}
	return
}
