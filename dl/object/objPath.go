package object

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/rt"
)

// Dot - access to a value inside another value.
// ex. in objects, lists, or records.
type Dot interface {
	Resolve(rt.Runtime) (rt.Dotted, error)
}

// change num and text evals into concrete index and member names
// ( determining the values of the path elements in advance aids debugging. )
func resolveDots(run rt.Runtime, dots []Dot) (ret []rt.Dotted, err error) {
	if cnt := len(dots); cnt > 0 {
		path := make([]rt.Dotted, 0, cnt)
		for _, el := range dots {
			if p, e := el.Resolve(run); e != nil {
				// if str := path.String(); len(str) == 0 {
				// 	err = e
				// } else {
				err = fmt.Errorf("%w with partial path %v", e, path)
				// }
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
