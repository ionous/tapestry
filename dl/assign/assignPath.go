package assign

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/dl/assign/dot"
	"git.sr.ht/~ionous/tapestry/rt"
)

// Dot - access to a value inside another value.
// ex. in objects, lists, or records.
type Dot interface {
	Resolve(rt.Runtime) (dot.Dotted, error)
}

// change num and text evals into concrete index and member names
// ( determining the values of the path elements in advance aids debugging. )
func ResolvePath(run rt.Runtime, dots []Dot) (ret dot.Path, err error) {
	if cnt := len(dots); cnt > 0 {
		path := make(dot.Path, 0, cnt)
		for _, el := range dots {
			if p, e := el.Resolve(run); e != nil {
				if str := path.String(); len(str) == 0 {
					err = e
				} else {
					err = fmt.Errorf("%w with partial path %s", e, str)
				}
				break
			} else {
				path = append(path, p)
			}
		}
		if err == nil {
			ret = path
		}
	}
	return
}
